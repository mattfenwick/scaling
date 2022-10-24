package telemetry

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

var keyValCounter *prometheus.CounterVec
var apiDurationHistogram *prometheus.HistogramVec
var eventLoopDurationHistogram *prometheus.HistogramVec

func RecordKeyValEvent(name string, value string) {
	labels := prometheus.Labels{"name": name, "value": value}
	keyValCounter.With(labels).Inc()
}

func RecordAPIDuration(path string, method string, code int, start time.Time) {
	duration := time.Since(start)
	labels := prometheus.Labels{"path": path, "method": method, "code": fmt.Sprintf("%d", code)}
	apiDurationHistogram.With(labels).Observe(float64(duration / time.Millisecond))
}

func RecordEventLoopDuration(name string, err error, start time.Time) {
	duration := time.Since(start)
	labels := prometheus.Labels{"name": name, "isError": fmt.Sprintf("%t", err != nil)}
	eventLoopDurationHistogram.With(labels).Observe(float64(duration / time.Millisecond))
}

func CreateMetrics(namespace string) {
	apiDurationHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: "api",
		Name:      "duration_histogram_milliseconds",
		Help:      "record duration of API endpoints in milliseconds",
		Buckets:   prometheus.ExponentialBuckets(1, 2, 20),
	}, []string{"path", "method", "code"})
	prometheus.MustRegister(apiDurationHistogram)

	eventLoopDurationHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: "api",
		Name:      "event_loop_duration_histogram_milliseconds",
		Help:      "record duration of event loop handler",
		Buckets:   prometheus.ExponentialBuckets(1, 2, 20),
	}, []string{"name", "isError"})
	prometheus.MustRegister(eventLoopDurationHistogram)

	keyValCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "api",
		Name:      "keyval_counter",
		Help:      "event counts by keyval",
	}, []string{"name", "value"})
	prometheus.MustRegister(keyValCounter)
}
