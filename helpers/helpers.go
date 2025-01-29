package helpers

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// initOpenTelemetry initializes OpenTelemetry for traces, metrics, and logs.
// initOpenTelemetry(
func InitOpenTelemetry(
	ctx context.Context,
) (*trace.TracerProvider, *metric.MeterProvider, *log.LoggerProvider, error) {
	// Create a resource describing the service
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("rabbitmq-worker"),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	// Initialize Traces
	traceExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, nil, nil, err
	}
	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tracerProvider)

	// Initialize Metrics
	metricExporter, err := stdoutmetric.New()
	if err != nil {
		return nil, nil, nil, err
	}
	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter)),
		metric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)

	// Initialize Logs
	logExporter, err := stdoutlog.New()
	if err != nil {
		return nil, nil, nil, err
	}
	loggerProvider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(logExporter)),
		log.WithResource(res),
	)
	global.SetLoggerProvider(loggerProvider)

	return tracerProvider, meterProvider, loggerProvider, nil
}

// shutdownOpenTelemetry shuts down the OpenTelemetry providers.
func ShutdownOpenTelemetry(
	ctx context.Context,
	tp *trace.TracerProvider,
	mp *metric.MeterProvider,
	lp *log.LoggerProvider,
) {
	if tp != nil {
		if err := tp.Shutdown(ctx); err != nil {
			fmt.Printf("Error shutting down tracer provider: %v", err)
		}
	}
	if mp != nil {
		if err := mp.Shutdown(ctx); err != nil {
			fmt.Printf("Error shutting down meter provider: %v", err)
		}
	}
	if lp != nil {
		if err := lp.Shutdown(ctx); err != nil {
			fmt.Printf("Error shutting down logger provider: %v", err)
		}
	}
}

func Log() {
	fmt.Println("Logging.....")
}
