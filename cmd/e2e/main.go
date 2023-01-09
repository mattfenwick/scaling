package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mattfenwick/collections/pkg/json"
	"github.com/mattfenwick/scaling/pkg/utils"
	"github.com/mattfenwick/scaling/pkg/webserver"
	"github.com/sirupsen/logrus"
)

/*
/users
 - POST => fancy search
 - GET => naive list of users

/message
 - POST => create
 - GET => by message's uuid
/messages
 - POST => fancy search
 - GET => naive list of messages

/follow
 - POST
/followers
 - GET => by user's uuid

/upvote
 - POST

*/

func main() {
	logrus.SetLevel(logrus.InfoLevel)

	serverUrl := "http://scaling-example.local:80"
	client := webserver.NewClient(serverUrl)

	// make user
	createUser1 := utils.DoOrDie(client.CreateUser(context.TODO(), &webserver.CreateUserRequest{Name: "e2e-user-1", Email: "e2e-user-1@scaling.local"}))
	fmt.Printf("create user1: %s\n", json.MustMarshalToString(createUser1))

	createUser2 := utils.DoOrDie(client.CreateUser(context.TODO(), &webserver.CreateUserRequest{Name: "e2e-user-2", Email: "e2e-user-2@scaling.local"}))
	fmt.Printf("create user2: %s\n", json.MustMarshalToString(createUser2))

	// make message
	mess11 := utils.DoOrDie(client.CreateMessage(context.TODO(), &webserver.CreateMessageRequest{SenderUserId: createUser1.UserId, Content: "this is message 1, from user 1 -- additional filler"}))
	fmt.Printf("create message 1-1: %s\n", json.MustMarshalToString(mess11))
	mess12 := utils.DoOrDie(client.CreateMessage(context.TODO(), &webserver.CreateMessageRequest{SenderUserId: createUser1.UserId, Content: "this is message 2, from user 1 -- additional filler"}))
	fmt.Printf("create message 1-2: %s\n", json.MustMarshalToString(mess12))
	mess21 := utils.DoOrDie(client.CreateMessage(context.TODO(), &webserver.CreateMessageRequest{SenderUserId: createUser2.UserId, Content: "this is message 1, from user 2 -- additional filler"}))
	fmt.Printf("create message 2-1: %s\n", json.MustMarshalToString(mess21))
	// get message
	fmt.Printf("get message 1-1 (%s): %s\n", mess11.MessageId.String(), json.MustMarshalToString(utils.DoOrDie(client.GetMessage(context.TODO(), &webserver.GetMessageRequest{MessageId: mess11.MessageId}))))
	fmt.Printf("get message 1-2 (%s): %s\n", mess12.MessageId.String(), json.MustMarshalToString(utils.DoOrDie(client.GetMessage(context.TODO(), &webserver.GetMessageRequest{MessageId: mess12.MessageId}))))
	fmt.Printf("get message 2-1 (%s): %s\n", mess21.MessageId.String(), json.MustMarshalToString(utils.DoOrDie(client.GetMessage(context.TODO(), &webserver.GetMessageRequest{MessageId: mess21.MessageId}))))

	// user2 follows user1
	utils.DoOrDie(client.FollowUser(context.TODO(), &webserver.FollowRequest{FolloweeUserId: createUser1.UserId, FollowerUserId: createUser2.UserId}))

	// upvotes
	utils.DoOrDie(client.UpvoteMessage(context.TODO(), &webserver.CreateUpvoteRequest{UserId: createUser1.UserId, MessageId: mess21.MessageId}))
	utils.DoOrDie(client.UpvoteMessage(context.TODO(), &webserver.CreateUpvoteRequest{UserId: createUser2.UserId, MessageId: mess11.MessageId}))
	utils.DoOrDie(client.UpvoteMessage(context.TODO(), &webserver.CreateUpvoteRequest{UserId: createUser2.UserId, MessageId: mess21.MessageId}))

	// get messages
	fmt.Printf("get messages: %s\n", json.MustMarshalToString(utils.DoOrDie(client.GetMessages(context.TODO(), &webserver.GetMessagesRequest{}))))
	// search messages
	for _, searchPhrase := range []string{"message 1", "message 2", "user 1", "user 2", "additional"} {
		fmt.Printf("search messages for '%s': %s\n", searchPhrase, json.MustMarshalToString(utils.DoOrDie(client.SearchMessages(context.TODO(), &webserver.SearchMessagesRequest{LiteralString: searchPhrase}))))
	}

	// search users
	for _, pair := range [][2]string{{"e2e", ""}, {"user-1", ""}, {"user-2", ""}} {
		client.SearchUsers(context.TODO(), &webserver.SearchUsersRequest{NamePattern: pair[0], EmailPattern: pair[1]})
	}

	// get: user, user timeline, user messages, followers
	for ix, userId := range []uuid.UUID{createUser1.UserId, createUser2.UserId} {
		fmt.Printf("get user %d (%s): %s\n", ix, userId.String(), json.MustMarshalToString(utils.DoOrDie(client.GetUser(context.TODO(), &webserver.GetUserRequest{UserId: userId}))))

		fmt.Printf("get user timeline %d (%s): %s\n", ix, userId.String(), json.MustMarshalToString(utils.DoOrDie(client.GetUserTimeline(context.TODO(), &webserver.GetUserTimelineRequest{UserId: userId}))))

		fmt.Printf("get user messages %d (%s): %s\n", ix, userId.String(), json.MustMarshalToString(utils.DoOrDie(client.GetUserMessages(context.TODO(), &webserver.GetUserMessagesRequest{UserId: userId}))))

		fmt.Printf("get user followers %d (%s): %s\n", ix, userId.String(), json.MustMarshalToString(utils.DoOrDie(client.GetFollowers(context.TODO(), &webserver.GetFollowersOfUserRequest{UserId: createUser1.UserId}))))
	}

	// get users
	fmt.Printf("get users: %s\n", json.MustMarshalToString(utils.DoOrDie(client.GetUsers(context.TODO(), &webserver.GetUsersRequest{}))))

	// search users
	name, email := "roc", "XAN"
	apiUsers, err := client.SearchUsers(context.TODO(), &webserver.SearchUsersRequest{NamePattern: name, EmailPattern: email})
	utils.Die(err)
	fmt.Printf("api users: %s\n", json.MustMarshalToString(apiUsers))
}
