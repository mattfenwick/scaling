package loadgen

import (
	"context"

	"github.com/mattfenwick/scaling/pkg/utils"
	"github.com/mattfenwick/scaling/pkg/webserver"
	"github.com/pkg/errors"
)

type Config struct {
	Mode              string
	Workers           int
	PauseMilliseconds int
}

func Cli(client *webserver.Client, config *Config) {
	ctx := context.TODO()

	uploader := NewGenerator(ctx, client)

	switch config.Mode {
	case "create-users":
		uploader.CreateUsers(ctx, 10)
	default:
		utils.DoOrDie(errors.Errorf("invalid mode: %s", config.Mode))
	}
}
