# tracer

The `tracer` package provides a simple and configurable way to set up distributed tracing for Go applications using [OpenTelemetry](https://opentelemetry.io/). It enables you to instrument your services, export traces to an OTLP-compatible collector, and manage tracer lifecycle with minimal setup.

---

## Features

- **Configurable Tracing**: Enable or disable tracing and set the OTLP endpoint via environment variables.
- **OpenTelemetry Integration**: Uses OpenTelemetry SDK and OTLP HTTP exporter for trace collection.
- **Resource Attribution**: Automatically attaches service and application metadata to traces.
- **Graceful Shutdown**: Provides a shutdown function to flush and close the exporter cleanly.
- **Simple API**: Minimal functions to initialize and use tracing in your application.

---

## Usage

### Configuration

The package uses the following configuration struct:

```go
type Config struct {
    Enable   bool   `env:"TRACE_ENABLE" envDefault:"1"`
    Endpoint string `env:"TRACE_ENDPOINT" envDefault:"simplest-collector.monitoring.svc.cluster.local:4318"`
}
```

### Initialization

Initialize the tracer in your application startup:

```go
import "github.com/philippe-berto/tracer"

shutdown, err := tracer.New(ctx, cfg, "my-service", "my-app", log)
if err != nil {
    // handle error
}
defer shutdown(ctx)
```

### Creating Spans

To create a new span:

```go
ctx, span := tracer.BaseTracer(ctx, "my-tracer", "1.0.0", "operation-name")
defer span.End()
```

### No-op Shutdown

If tracing is disabled, tracer.New returns a no-op shutdown function.

## API Reference

- **Config**: Tracing configuration struct.
- **New**: Initializes the tracer provider and exporter.
- **BaseTracer**: Starts a new span with the given tracer name and version.
- **ShutdownNull**: No-op shutdown function for disabled tracing.

## Dependencies

[OpenTelemetry Go SDK](https://pkg.go.dev/go.opentelemetry.io/otel)  
[logger](https://github.com/philippe-berto/logger)
