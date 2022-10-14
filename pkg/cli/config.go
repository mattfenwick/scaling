package cli

type Config struct {
	LogLevel       string
	JaegerURL      string
	Mode           string
	Port           int
	PrometheusPort int
}
