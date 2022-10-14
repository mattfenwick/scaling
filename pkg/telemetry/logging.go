package telemetry

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func SetUpLogger(level string) error {
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return errors.Wrapf(err, "unable to parse the specified log level: '%s'", level)
	}
	logrus.SetLevel(logLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.Infof("log level set to '%s'", logrus.GetLevel())
	return nil
}
