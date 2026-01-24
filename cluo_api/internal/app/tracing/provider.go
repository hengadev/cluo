package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

// Config holds tracing configuration.
type Config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	Endpoint       string
	SampleRate     float64
}

// Provider wraps the OpenTelemetry TracerProvider.
type Provider struct {
	tracerProvider *sdktrace.TracerProvider
}

// Setup initializes OpenTelemetry tracing and returns a shutdown function.
func Setup(ctx context.Context, cfg Config) (*Provider, func(context.Context) error, error) {
	// Create resource with service information
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			semconv.DeploymentEnvironment(cfg.Environment),
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("create resource: %w", err)
	}

	// Create OTLP exporter
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(cfg.Endpoint),
		otlptracegrpc.WithInsecure(), // Use WithTLS in production
	)
	if err != nil {
		return nil, nil, fmt.Errorf("create exporter: %w", err)
	}

	// Create sampler
	var sampler sdktrace.Sampler
	if cfg.SampleRate >= 1.0 {
		sampler = sdktrace.AlwaysSample()
	} else if cfg.SampleRate <= 0 {
		sampler = sdktrace.NeverSample()
	} else {
		sampler = sdktrace.TraceIDRatioBased(cfg.SampleRate)
	}

	// Create TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
	)

	// Set global TracerProvider
	otel.SetTracerProvider(tp)

	// Set global propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	provider := &Provider{
		tracerProvider: tp,
	}

	// Return shutdown function
	shutdown := func(ctx context.Context) error {
		return tp.Shutdown(ctx)
	}

	return provider, shutdown, nil
}

// Tracer returns a tracer for the given name.
func (p *Provider) Tracer(name string) trace.Tracer {
	return p.tracerProvider.Tracer(name)
}

// SetupNoop sets up a no-op tracer provider for when tracing is disabled.
func SetupNoop() {
	otel.SetTracerProvider(trace.NewNoopTracerProvider())
}
