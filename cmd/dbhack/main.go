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
	utils.DoOrDie(err)

	ctx := context.Background()

	// users
	user1 := database.NewUser("abc def", "abcdef@whatever.com")
	user2 := database.NewUser("qrs xyz", "qrsxyz@whatever.com")

	err = database.InsertUser(ctx, db, user1)
	utils.DoOrDie(err)
	err = database.InsertUser(ctx, db, user2)
	utils.DoOrDie(err)

	users, err := database.ReadAllUsers(ctx, db)
	utils.DoOrDie(err)
	fmt.Printf("users: %s\n", json.MustMarshalToString(users))

	// followers
	follower1 := database.NewFollower(user1.UserId, user2.UserId)
	err = database.InsertFollower(ctx, db, follower1)
	utils.DoOrDie(err)

	followers, err := database.ReadAllFollowers(ctx, db)
	utils.DoOrDie(err)
	fmt.Printf("followers: %s\n", json.MustMarshalToString(followers))

	// messages
	message1 := database.NewMessage(user1.UserId, "hi, i'm user 1")
	message2 := database.NewMessage(user2.UserId, "whereas I'm user 2")

	err = database.InsertMessage(ctx, db, message1)
	utils.DoOrDie(err)
	err = database.InsertMessage(ctx, db, message2)
	utils.DoOrDie(err)

	messages, err := database.ReadAllMessages(ctx, db)
	utils.DoOrDie(err)
	fmt.Printf("messages: %s\n", json.MustMarshalToString(messages))

	// upvotes
	upvote1 := database.NewUpvote(user1.UserId, message1.MessageId)
	upvote2 := database.NewUpvote(user1.UserId, message2.MessageId)

	err = database.InsertUpvote(ctx, db, upvote1)
	utils.DoOrDie(err)
	err = database.InsertUpvote(ctx, db, upvote2)
	utils.DoOrDie(err)

	upvotes, err := database.ReadAllUpvotes(ctx, db)
	utils.DoOrDie(err)
	fmt.Printf("upvotes: %s\n", json.MustMarshalToString(upvotes))

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
	utils.DoOrDie(err)

	//dbname := "scaling"
	docs, err := database.ReadDocuments(context.TODO(), db)
	utils.DoOrDie(err)
	fmt.Printf("docs: %s\n", json.MustMarshalToString(docs))
	for _, doc := range docs {
		bytes, err := base64.StdEncoding.DecodeString(string(doc.Parsed.([]uint8)))
		fmt.Printf("doc? %T\n", doc.Parsed)
		utils.DoOrDie(err)
		fmt.Printf("???? %s\n", bytes)
	}
}
