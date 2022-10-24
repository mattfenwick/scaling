package cli

import "github.com/mattfenwick/scaling/pkg/loadgen"

type Config struct {
	LogLevel       string
	JaegerURL      string
	PrometheusPort int

	Webserver struct {
		Host          string
		ContainerPort int
		ServicePort   int
	}

	LoadGen loadgen.Config
}
