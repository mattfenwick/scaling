package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/mattfenwick/collections/pkg/json"
	"github.com/mattfenwick/scaling/pkg/database"
	"github.com/mattfenwick/scaling/pkg/loadgen"
	"github.com/mattfenwick/scaling/pkg/telemetry"
	"github.com/mattfenwick/scaling/pkg/utils"
	"github.com/mattfenwick/scaling/pkg/webserver"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func Run() {
	RunVersionCommand()

	mode := os.Args[1]

	config, err := json.ParseFile[Config](os.Args[2])
	utils.DoOrDie(err)

	RunWithConfig(mode, config)
}

func RunWithConfig(mode string, config *Config) {
	rootContext := context.Background()

	tp, err, cleanup := telemetry.Setup(rootContext, config.LogLevel, mode, config.PrometheusPort, config.JaegerURL)
	defer cleanup()
	utils.DoOrDie(err)

	switch mode {
	case "schema":
		pg := config.Postgres

		adminDb, err := database.Connect(pg.User, pg.Password, pg.Host, pg.AdminDatabase)
		utils.DoOrDie(err)
		utils.DoOrDie(database.CreateDatabaseIfNotExists(rootContext, adminDb, pg.Database))

		db, err := database.Connect(pg.User, pg.Password, pg.Host, pg.Database)
		utils.DoOrDie(err)
		utils.DoOrDie(database.InitializeSchema(rootContext, db))
	case "webserver":
		pg := config.Postgres

		adminDb, err := database.Connect(pg.User, pg.Password, pg.Host, pg.AdminDatabase)
		utils.DoOrDie(err)
		utils.DoOrDie(database.CreateDatabaseIfNotExists(rootContext, adminDb, pg.Database))

		db, err := database.Connect(pg.User, pg.Password, pg.Host, pg.Database)
		utils.DoOrDie(err)
		webserver.Run(config.Webserver.ContainerPort, tp, db)
	case "loadgen":
		url := fmt.Sprintf("http://%s:%d", config.Webserver.Host, config.Webserver.ServicePort)
		client := webserver.NewClient(url)
		loadgen.Cli(client, &config.LoadGen)
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
