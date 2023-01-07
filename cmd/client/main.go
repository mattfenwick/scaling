package main

import (
	"context"
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
	} else {
		cli.Run()
	}
}
