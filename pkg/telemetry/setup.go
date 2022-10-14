package telemetry

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

func Setup(ctx context.Context, logLevel string, serviceName string, jaegerURL string) (error, func()) {
	cleanup := func() {
		logrus.Infof("noop tracing cleanup")
	}

	// logs
	logrus.Infof("setting up logging for level %s", logLevel)
	err := SetUpLogger(logLevel)
	if err != nil {
		return err, cleanup
	}

	// metrics
	logrus.Infof("setting up metrics for namespace %s", serviceName)
	SetupMetrics(serviceName) // TODO is this really what we want/

	// traces
	logrus.Infof("setting up tracing for jaeger url %s", jaegerURL)
	if jaegerURL == "" {
		_ = SetUpNoopTracerProvider()
	} else {
		tp, err := SetUpJaegerTracerProvider(jaegerURL, serviceName)
		if err != nil {
			return err, cleanup
		}

		cleanup = func() {
			logrus.Infof("jaeger tracing cleanup")
			timedContext, timedCancel := context.WithTimeout(ctx, time.Second*5)
			defer timedCancel()
			_ = tp.Shutdown(timedContext)
		}
	}

	return nil, cleanup
}
