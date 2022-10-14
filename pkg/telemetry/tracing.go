package telemetry

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
)

func SetUpJaegerTracerProvider(aggregatorURL string, service string) (*tracesdk.TracerProvider, error) {
	logrus.Infof("setting up jaeger tracer provider at %s for service %s", aggregatorURL, service)

	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(aggregatorURL)))
	if err != nil {
		return nil, errors.Wrapf(err, "unable to instantiate jaeger tracer provider")
	}
	tracerProvider := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
		)),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tracerProvider, nil
}

func SetUpNoopTracerProvider() trace.TracerProvider {
	logrus.Infof("setting up noop tracer provider")
	tracerProvider := trace.NewNoopTracerProvider()
	otel.SetTracerProvider(tracerProvider)
	return tracerProvider
}
