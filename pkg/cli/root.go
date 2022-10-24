package cli

import (
	"context"
	"github.com/mattfenwick/collections/pkg/json"
	"github.com/mattfenwick/scaling/pkg/client"
	"github.com/mattfenwick/scaling/pkg/parse"
	"github.com/mattfenwick/scaling/pkg/telemetry"
	"github.com/mattfenwick/scaling/pkg/utils"
	"github.com/mattfenwick/scaling/pkg/webserver"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
)

func Run() {
	RunVersionCommand()

	mode := os.Args[1]

	config, err := json.ParseFile[Config](os.Args[2])
	utils.DoOrDie(err)

	rootContext := context.Background()

	tp, err, cleanup := telemetry.Setup(rootContext, config.LogLevel, mode, config.PrometheusPort, config.JaegerURL)
	defer cleanup()
	utils.DoOrDie(err)

	switch mode {
	case "webserver":
		webserver.Run(config.Webserver.ContainerPort, tp)
	case "loadgen":
		utils.DoOrDie(client.RunSmallBatchOfRequests(config.Webserver.Host, config.Webserver.ServicePort))
	case "parser":
		result := parse.JsonObject("{}")
		logrus.Infof("%+v", json.MustMarshalToString(result))
		panic(errors.Errorf("TODO"))
	default:
		panic(errors.Errorf("invalid mode: %s", mode))
	}
}

var (
	version   = "development"
	gitSHA    = "development"
	buildTime = "development"
)

func RunVersionCommand() {
	jsonString, err := json.MarshalToString(map[string]string{
		"Version":   version,
		"GitSHA":    gitSHA,
		"BuildTime": buildTime,
	})
	utils.DoOrDie(err)
	logrus.Infof("scaling version: \n%s\n", jsonString)
}
