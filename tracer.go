package tracer

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const baseVersion = "1.0.0"

type (
	Config struct {
		Enable   bool   `env:"TRACE_ENABLE" envDefault:"1"`
		Endpoint string `env:"TRACE_ENDPOINT" envDefault:"simplest-collector.monitoring.svc.cluster.local:4318"`
	}
)

func New(ctx context.Context, cfg Config, serviceName, appName string) (func(context.Context) error, error) {
	if !cfg.Enable {
		return ShutdownNull, nil
	}

	exporter, err := otlptrace.New(ctx,
		otlptracehttp.NewClient(
			otlptracehttp.WithInsecure(),
			otlptracehttp.WithEndpoint(cfg.Endpoint),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("Otel tracer: Could not setup exporter: %w", err)
	}

	resources, err := resource.New(ctx,
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("application", appName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("Otel tracer: Could not setup resources: %w", err)
	}

	tracer := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exporter)),
		sdktrace.WithResource(resources),
	)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tracer)

	return exporter.Shutdown, nil
}

func BaseTracer(ctx context.Context, tracerName, tracerVersion, spanName string) (context.Context, trace.Span) {
	if tracerVersion == "" {
		tracerVersion = baseVersion
	}

	tracer := otel.GetTracerProvider().Tracer(
		tracerName,
		trace.WithInstrumentationVersion("semver:"+tracerVersion),
	)

	return tracer.Start(ctx, spanName)
}

func ShutdownNull(ctx context.Context) error {
	return nil
}
