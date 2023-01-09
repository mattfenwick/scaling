package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/mattfenwick/collections/pkg/json"
	"github.com/mattfenwick/scaling/pkg/cli"
	"github.com/mattfenwick/scaling/pkg/database"
	"github.com/mattfenwick/scaling/pkg/utils"
	"github.com/mattfenwick/scaling/pkg/webserver"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)

	isSimple := len(os.Args) < 2 || os.Args[1] != "false"
	if isSimple {
		myClient := webserver.NewClient("http://scaling-example.local:80")
		createResp, err := myClient.CreateUser(context.TODO(), &webserver.CreateUserRequest{Name: "abc", Email: "abc@def.org"})
		utils.DoOrDie(err)
		logrus.Infof("create response: %+v", createResp)

		getResp, err := myClient.GetUser(context.TODO(), &webserver.GetUserRequest{UserId: createResp.UserId})
		utils.DoOrDie(err)
		logrus.Infof("get response: %+v", getResp)

		db, err := database.Connect("postgres", "postgres", "localhost", "scaling")
		utils.DoOrDie(err)

		name, email := "roc", "XAN"
		dbUsers, err := database.SearchUsers(context.TODO(), db, name, email)
		utils.DoOrDie(err)
		fmt.Printf("db users: %+v\n", json.MustMarshalToString(dbUsers))

		apiUsers, err := myClient.SearchUsers(context.TODO(), &webserver.SearchUsersRequest{NamePattern: name, EmailPattern: email})
		utils.DoOrDie(err)
		fmt.Printf("api users: %s\n", json.MustMarshalToString(apiUsers))

		for _, user := range dbUsers {
			timelineMessages, err := database.GetUserTimeline(context.TODO(), db, user.UserId)
			utils.DoOrDie(err)
			fmt.Printf("timeline for user %s (%s, %s):\n%s\n\n", user.UserId.String(), user.Name, user.Email, json.MustMarshalToString(timelineMessages))

			userMessages, err := database.GetUserMessages(context.TODO(), db, user.UserId)
			utils.DoOrDie(err)
			fmt.Printf("messages sent by user %s (%s, %s):\n%s\n\n", user.UserId.String(), user.Name, user.Email, json.MustMarshalToString(userMessages))
		}

		userTimelineAndMessageTest(db, myClient)

		if false {
			messagesTest(myClient)
		}

		tableSizes(db)
		searchMessages(db)
	} else {
		cli.Run()
	}
}

func userTimelineAndMessageTest(db *sql.DB, client *webserver.Client) {
	if db != nil {
		// create user objects
		user1 := database.NewUser("utamt1", "utamt1@scaling.local")
		user2 := database.NewUser("utamt-two", "utamt-two@scaling.local")
		// insert users
		utils.DoOrDie(database.InsertUser(context.TODO(), db, user1))
		utils.DoOrDie(database.InsertUser(context.TODO(), db, user2))
		// have user2 follow user1
		utils.DoOrDie(database.InsertFollower(context.TODO(), db, database.NewFollower(user1.UserId, user2.UserId)))
		// create message objects
		message1user1 := database.NewMessage(user1.UserId, "this is message 1, from user 1")
		message2user1 := database.NewMessage(user1.UserId, "this is message 2, from user 1")
		message1user2 := database.NewMessage(user2.UserId, "this is message 1, from user 2")
		// insert messages
		utils.DoOrDie(database.InsertMessage(context.TODO(), db, message1user1))
		utils.DoOrDie(database.InsertMessage(context.TODO(), db, message2user1))
		utils.DoOrDie(database.InsertMessage(context.TODO(), db, message1user2))

		// look at timelines, messages
		timeline1, err := database.GetUserTimeline(context.TODO(), db, user1.UserId)
		utils.DoOrDie(err)
		messages1, err := database.GetUserMessages(context.TODO(), db, user1.UserId)
		utils.DoOrDie(err)
		fmt.Printf("user1 (%s) timeline and messages:\n%s\n\n", user1.UserId.String(), json.MustMarshalToString(map[string]any{"timeline": timeline1, "messages": messages1}))

		timeline2, err := database.GetUserTimeline(context.TODO(), db, user2.UserId)
		utils.DoOrDie(err)
		messages2, err := database.GetUserMessages(context.TODO(), db, user2.UserId)
		utils.DoOrDie(err)
		fmt.Printf("user2 (%s) timeline and messages:\n%s\n\n", user2.UserId.String(), json.MustMarshalToString(map[string]any{"timeline": timeline2, "messages": messages2}))
	}

	if client != nil {
		// create user objects
		user1, err := client.CreateUser(context.TODO(), &webserver.CreateUserRequest{Name: "utamt1-client", Email: "utamt1-client@scaling.local"})
		utils.DoOrDie(err)
		user2, err := client.CreateUser(context.TODO(), &webserver.CreateUserRequest{Name: "utamt-two-client", Email: "utamt-two-client@scaling.local"})
		utils.DoOrDie(err)

		// have user2 follow user1
		_, err = client.FollowUser(context.TODO(), &webserver.FollowRequest{FolloweeUserId: user1.UserId, FollowerUserId: user2.UserId})
		utils.DoOrDie(err)

		// create message objects
		_, err = client.CreateMessage(context.TODO(), &webserver.CreateMessageRequest{SenderUserId: user1.UserId, Content: "this is message 1, from user 1"})
		utils.DoOrDie(err)
		_, err = client.CreateMessage(context.TODO(), &webserver.CreateMessageRequest{SenderUserId: user1.UserId, Content: "this is message 2, from user 1"})
		utils.DoOrDie(err)
		_, err = client.CreateMessage(context.TODO(), &webserver.CreateMessageRequest{SenderUserId: user2.UserId, Content: "this is message 1, from user 2"})
		utils.DoOrDie(err)

		// look at timelines, messages
		timeline1, err := client.GetUserTimeline(context.TODO(), &webserver.GetUserTimelineRequest{UserId: user1.UserId})
		utils.DoOrDie(err)
		messages1, err := client.GetUserMessages(context.TODO(), &webserver.GetUserMessagesRequest{UserId: user1.UserId})
		utils.DoOrDie(err)
		fmt.Printf("(client) user1 (%s) timeline and messages:\n%s\n\n", user1.UserId.String(), json.MustMarshalToString(map[string]any{"timeline": timeline1, "messages": messages1}))

		timeline2, err := client.GetUserTimeline(context.TODO(), &webserver.GetUserTimelineRequest{UserId: user2.UserId})
		utils.DoOrDie(err)
		messages2, err := client.GetUserMessages(context.TODO(), &webserver.GetUserMessagesRequest{UserId: user2.UserId})
		utils.DoOrDie(err)
		fmt.Printf("(client) user2 (%s) timeline and messages:\n%s\n\n", user2.UserId.String(), json.MustMarshalToString(map[string]any{"timeline": timeline2, "messages": messages2}))
	}
}

func messagesTest(client *webserver.Client) {
	messages, err := client.GetMessages(context.TODO(), &webserver.GetMessagesRequest{})
	utils.DoOrDie(err)

	for _, message := range messages.Messages {
		refetchedMessage, err := client.GetMessage(context.TODO(), &webserver.GetMessageRequest{MessageId: message.MessageId})
		utils.DoOrDie(err)

		fmt.Printf("refetched: %+v\n", refetchedMessage)
	}
}

func tableSizes(db *sql.DB) {
	rowCounts, err := database.GetTableSizes(context.TODO(), db)
	utils.DoOrDie(err)
	fmt.Printf("row counts: %+v\n", rowCounts)
}

func searchMessages(db *sql.DB) {
	messages, err := database.SearchMessages(context.TODO(), db, "banan")
	utils.DoOrDie(err)
	fmt.Printf("searched messages, found %d:\n%s\n", len(messages), json.MustMarshalToString(messages))
}
