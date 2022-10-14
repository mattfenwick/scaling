package telemetry

import (
	"context"
	"fmt"
	"github.com/mattfenwick/scaling/pkg/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func Setup(ctx context.Context, logLevel string, serviceName string, prometheusPort int, jaegerURL string) (error, func()) {
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
	CreateMetrics(serviceName) // TODO is this really what we want/
	SetupPrometheus(prometheusPort)

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

func SetupPrometheus(port int) {
	addr := fmt.Sprintf(":%d", port)

	serveMux := http.NewServeMux()
	serveMux.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars. -- only on go 1.17
			//EnableOpenMetrics: true,
			Timeout: 10 * time.Second,
		},
	))
	go func() {
		utils.DoOrDie(http.ListenAndServe(addr, serveMux))
	}()
}
