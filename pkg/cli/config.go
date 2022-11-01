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

	Postgres *PostgresConfig

	LoadGen loadgen.Config
}

type PostgresConfig struct {
	Host     string
	User     string
	Password string
	Database string
}
