package loadgen

import (
	"context"
	"github.com/mattfenwick/scaling/pkg/utils"
	"github.com/mattfenwick/scaling/pkg/webserver"
	"github.com/pkg/errors"
)

type Config struct {
	Mode              string
	KeyCounts         []int
	Workers           int
	PauseMilliseconds int
}

func Cli(client *webserver.Client, config *Config) {
	uploader := NewUploader(client)

	ctx := context.TODO()

	switch config.Mode {
	case "canned":
		utils.DoOrDie(uploader.RunCannedUploads())
	case "by-key-count":
		uploader.RunRandomUploadsByKeyCount(config.KeyCounts)
	case "continuous":
		uploader.RunContinuous(ctx, config.KeyCounts, config.Workers, config.PauseMilliseconds)
		<-ctx.Done()
	default:
		utils.DoOrDie(errors.Errorf("invalid mode: %s", config.Mode))
	}
}
