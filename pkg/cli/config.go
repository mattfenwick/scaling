package cli

type Config struct {
	LogLevel       string
	JaegerURL      string
	PrometheusPort int

	Webserver struct {
		Host          string
		ContainerPort int
		ServicePort   int
	}

	LoadGen struct {
	}
}
