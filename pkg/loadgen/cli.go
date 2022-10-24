package loadgen

import (
	"github.com/mattfenwick/scaling/pkg/utils"
	"github.com/mattfenwick/scaling/pkg/webserver"
	"github.com/pkg/errors"
)

type Config struct {
	Mode string

	KeyCounts []int
}

func Cli(client *webserver.Client, config *Config) {
	uploader := NewUploader(client)

	switch config.Mode {
	case "canned":
		utils.DoOrDie(uploader.RunCannedUploads())
	case "by-key-count":
		uploader.RunRandomUploadsByKeyCount(config.KeyCounts)
	default:
		utils.DoOrDie(errors.Errorf("invalid mode: %s", config.Mode))
	}
}
