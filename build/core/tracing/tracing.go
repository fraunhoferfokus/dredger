package tracing

import (
	"build/core"
	"build/core/log"

	"go.opentelemetry.io/otel"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// Create tracer for the service
var Tracer = otel.Tracer(core.AppConfig.Sid)

// Create trace provider
func InitTracer() (*sdktrace.TracerProvider, error) {
	var exporter sdktrace.SpanExporter
	var err error
	// output traces locally on stdout
	exporter, err = stdout.New(stdout.WithPrettyPrint())
	if err != nil {
		log.Error().Err(err).Msg("creating stdout tracing exporter failed")
		return nil, err
	}
	log.Info().Msg("created stdout tracing exporter")

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.AlwaysSample())),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceVersion(core.Version), semconv.ServiceInstanceID(core.AppConfig.Sid), semconv.ServiceName(core.AppConfig.Service))),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}
