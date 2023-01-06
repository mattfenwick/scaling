package main

import (
	"context"
	"os"

	"github.com/mattfenwick/scaling/pkg/cli"
	"github.com/mattfenwick/scaling/pkg/utils"
	"github.com/mattfenwick/scaling/pkg/webserver"
	"github.com/sirupsen/logrus"
)

func main() {
	isSimple := len(os.Args) < 2 || os.Args[1] != "false"
	if isSimple {
		myClient := webserver.NewClient("http://localhost:8765")
		createResp, err := myClient.CreateUser(context.TODO(), &webserver.CreateUserRequest{Name: "abc", Email: "abc@def.org"})
		utils.DoOrDie(err)
		logrus.Infof("create response: %+v", createResp)

		getResp, err := myClient.GetUser(context.TODO(), &webserver.GetUserRequest{UserId: createResp.UserId})
		utils.DoOrDie(err)
		logrus.Infof("get response: %+v", getResp)
	} else {
		cli.Run()
	}
}
