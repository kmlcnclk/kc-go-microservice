package common

import (
	"context"
	"net/http"

	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func SetGlobalTracer(ctx context.Context, serviceName, exporterEndpoint string) (*trace.TracerProvider, error) {
	headers := map[string]string{
		"content-type": "application/json",
	}

	exporter, err := otlptrace.New(
		ctx,
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(exporterEndpoint),
			otlptracehttp.WithHeaders(headers),
			otlptracehttp.WithInsecure(),
		))

	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}

func ExtractTraceFromAMQPHeaders(headers amqp.Table) context.Context {
	traceparent, ok := headers["traceparent"].(string)
	if !ok || traceparent == "" {
		return context.Background()
	}

	carrier := propagation.MapCarrier{
		"traceparent": traceparent,
	}

	propagator := propagation.TraceContext{}
	ctx := propagator.Extract(context.Background(), carrier)

	return ctx
}

func TraceparentHeaderFromContext(ctx context.Context) string {
	carrier := propagation.HeaderCarrier(http.Header{})
	propagator := propagation.TraceContext{}

	propagator.Inject(ctx, carrier)

	traceparent := carrier.Get("traceparent")
	return traceparent
}
