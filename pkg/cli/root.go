package cli

import (
	"context"
	"github.com/mattfenwick/collections/pkg/json"
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

	config, err := json.ParseFile[Config](os.Args[1])
	utils.DoOrDie(err)

	rootContext := context.Background()

	err, cleanup := telemetry.Setup(rootContext, config.LogLevel, config.Mode, config.PrometheusPort, config.JaegerURL)
	defer cleanup()
	utils.DoOrDie(err)

	switch config.Mode {
	case "webserver":
		webserver.Run(config.Port)
	case "parser":
		// TODO
		result := parse.JsonObject("{}")
		logrus.Infof("%+v", json.MustMarshalToString(result))
	default:
		panic(errors.Errorf("invalid mode: %s", config.Mode))
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
