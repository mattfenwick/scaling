package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/mattfenwick/collections/pkg/json"
	"github.com/mattfenwick/scaling/pkg/database"
	"github.com/mattfenwick/scaling/pkg/utils"
)

func main() {
	user := "postgres"
	pw := "postgres"
	host := "localhost"
	initDbName := "hack"
	db, err := database.Connect(user, pw, host, initDbName)
	utils.Die(err)

	ctx := context.Background()

	// users
	user1 := database.NewUser("abc def", "abcdef@whatever.com")
	user2 := database.NewUser("qrs xyz", "qrsxyz@whatever.com")

	err = database.InsertUser(ctx, db, user1)
	utils.Die(err)
	err = database.InsertUser(ctx, db, user2)
	utils.Die(err)

	users, err := database.GetUsers(ctx, db)
	utils.Die(err)
	fmt.Printf("users: %s\n", json.MustMarshalToString(users))

	// followers
	follower1 := database.NewFollower(user1.UserId, user2.UserId)
	err = database.InsertFollower(ctx, db, follower1)
	utils.Die(err)

	followers, err := database.GetFollowers(ctx, db)
	utils.Die(err)
	fmt.Printf("followers: %s\n", json.MustMarshalToString(followers))

	// messages
	message1 := database.NewMessage(user1.UserId, "hi, i'm user 1")
	message2 := database.NewMessage(user2.UserId, "whereas I'm user 2")

	err = database.InsertMessage(ctx, db, message1)
	utils.Die(err)
	err = database.InsertMessage(ctx, db, message2)
	utils.Die(err)

	messages, err := database.GetMessages(ctx, db)
	utils.Die(err)
	fmt.Printf("messages: %s\n", json.MustMarshalToString(messages))

	// upvotes
	upvote1 := database.NewUpvote(user1.UserId, message1.MessageId)
	upvote2 := database.NewUpvote(user1.UserId, message2.MessageId)

	err = database.InsertUpvote(ctx, db, upvote1)
	utils.Die(err)
	err = database.InsertUpvote(ctx, db, upvote2)
	utils.Die(err)

	upvotes, err := database.ReadAllUpvotes(ctx, db)
	utils.Die(err)
	fmt.Printf("upvotes: %s\n", json.MustMarshalToString(upvotes))

	// find followers by user
	allUsers, err := database.GetUsers(ctx, db)
	utils.Die(err)
	for _, user := range allUsers {
		userFollowers, err := database.GetFollowersOfUser(ctx, db, user.UserId)
		utils.Die(err)
		for _, follower := range userFollowers {
			fmt.Printf("follower of %s (%s): %s (%s)\n", user.Name, user.UserId, follower.Name, follower.UserId)
		}

		timelineMessages, err := database.GetUserTimeline(ctx, db, user.UserId)
		utils.Die(err)
		for _, message := range timelineMessages {
			fmt.Printf("timeline message for %s (%s): %d upvotes, %s (%s)\n", user.Name, user.UserId, message.UpvoteCount, message.Content, message.MessageId)
		}
		if len(timelineMessages) == 0 {
			fmt.Printf("no timeline messages for %s\n", user.UserId)
		}
	}

	// done
	if false {
		doDocumentStuff(db)
	}
}

func doDocumentStuff(db *sql.DB) {
	insert := `insert into documents (parsed, parse_error) values($1, $2)`
	_, err := db.ExecContext(
		context.TODO(),
		insert,
		json.MustMarshalToString([]any{1, 2, 3, "hi", map[string]string{"qrs": "tuv"}}),
		"")
	utils.Die(err)

	//dbname := "scaling"
	docs, err := database.ReadDocuments(context.TODO(), db)
	utils.Die(err)
	fmt.Printf("docs: %s\n", json.MustMarshalToString(docs))
	for _, doc := range docs {
		bytes, err := base64.StdEncoding.DecodeString(string(doc.Parsed.([]uint8)))
		fmt.Printf("doc? %T\n", doc.Parsed)
		utils.Die(err)
		fmt.Printf("???? %s\n", bytes)
	}
}
