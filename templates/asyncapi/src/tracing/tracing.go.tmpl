// Package tracing initialisiert OpenTelemetry Tracing für Async-Service
package tracing

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// InitTracer startet Jaeger-Exporter mit ServiceName.
func InitTracer(serviceName string) (*sdktrace.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint())
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}
