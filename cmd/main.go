package main

import (
	"github.com/mattfenwick/collections/pkg/json"
	"github.com/mattfenwick/scaling/pkg/cli"
	"github.com/mattfenwick/scaling/pkg/parse"
	"github.com/sirupsen/logrus"
)

func main() {
	cli.RunRootSchemaCommand()

	result := parse.JsonAST("{}")
	logrus.Infof("%+v", json.MustMarshalToString(result))
}
