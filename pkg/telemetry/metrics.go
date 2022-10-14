package telemetry

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

var keyValCounter *prometheus.CounterVec
var durationHistogram *prometheus.HistogramVec

func RecordKeyValEvent(name string, value string) {
	labels := prometheus.Labels{"name": name, "value": value}
	keyValCounter.With(labels).Inc()
}

func RecordEventDuration(name string, code int, start time.Time) {
	duration := time.Since(start)
	labels := prometheus.Labels{"name": name, "code": fmt.Sprintf("%d", code)}
	durationHistogram.With(labels).Observe(float64(duration / time.Microsecond))
}

func SetupMetrics(namespace string) {
	durationHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: "api",
		Name:      "duration_histogram_microseconds",
		Help:      "record duration of API endpoints in microseconds",
		Buckets:   prometheus.ExponentialBuckets(1, 3, 20),
	}, []string{"name", "code"})
	prometheus.MustRegister(durationHistogram)

	keyValCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "api",
		Name:      "keyval_counter",
		Help:      "event counts by keyval",
	}, []string{"name", "value"})
	prometheus.MustRegister(keyValCounter)
}
