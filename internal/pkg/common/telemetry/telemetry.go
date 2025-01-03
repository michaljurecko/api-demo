// Package telemetry provides auxiliary code for OpenTelemetry.
package telemetry

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/michaljurecko/ch-demo/internal/pkg/common/shutdown"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"

	"github.com/michaljurecko/ch-demo/internal/pkg/common/log"
	"github.com/michaljurecko/ch-demo/internal/pkg/common/telemetry/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/metric"
	metricnnoop "go.opentelemetry.io/otel/metric/noop"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
)

func NewTracerProvider(ctx context.Context, logger *log.Logger, down *shutdown.Stack, cfg config.TraceConfig) (trace.TracerProvider, error) {
	setGlobalErrorHandler(ctx, logger)

	switch cfg.Exporter {
	case config.NoneTraceExporter:
		return tracenoop.NewTracerProvider(), nil
	case config.HTTPTraceExporter:
		exporter, err := otlptracehttp.New(ctx,
			otlptracehttp.WithEndpointURL(cfg.HTTPEndpointURL),
			otlptracehttp.WithHeaders(map[string]string{
				"Authorization": cfg.HTTPAuthorization,
			}),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create 'otlptracehttp' trace exporter: %w", err)
		}
		return traceProviderWith(ctx, logger, down, exporter)
	default:
		return nil, fmt.Errorf("unexpected trace exporter '%s'", cfg.Exporter)
	}
}

func NewMeterProvider(ctx context.Context, logger *log.Logger, down *shutdown.Stack, cfg config.MetricConfig) (metric.MeterProvider, error) {
	setGlobalErrorHandler(ctx, logger)

	switch cfg.Exporter {
	case config.NoneMetricExporter:
		return metricnnoop.NewMeterProvider(), nil
	case config.HTTPMetricExporter:
		exporter, err := otlpmetrichttp.New(ctx,
			otlpmetrichttp.WithEndpointURL(cfg.HTTPEndpointURL),
			otlpmetrichttp.WithHeaders(map[string]string{
				"Authorization": cfg.HTTPAuthorization,
				"stream-name":   "default",
			}),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create 'otlpmetrichttp' metric exporter: %w", err)
		}
		return meterProviderWith(ctx, logger, down, exporter)
	default:
		return nil, fmt.Errorf("unexpected metric exporter '%s'", cfg.Exporter)
	}
}

func traceProviderWith(ctx context.Context, logger *log.Logger, down *shutdown.Stack, exporter tracesdk.SpanExporter) (trace.TracerProvider, error) {
	res, err := newResource(ctx)
	if err != nil {
		return nil, err
	}

	provider := tracesdk.NewTracerProvider(
		tracesdk.WithResource(res),
		tracesdk.WithSyncer(
			exporter,
		),
	)

	down.OnShutdown(func() {
		ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
		defer cancel()
		if err := provider.Shutdown(ctx); err != nil {
			err = fmt.Errorf("cannot shutdown trace provider: %w", err)
			logger.Error(ctx, err.Error(), slog.Any(log.ErrorKey, err))
		}
	})

	return provider, nil
}

func meterProviderWith(ctx context.Context, logger *log.Logger, down *shutdown.Stack, exporter metricsdk.Exporter) (metric.MeterProvider, error) {
	res, err := newResource(ctx)
	if err != nil {
		return nil, err
	}

	provider := metricsdk.NewMeterProvider(
		metricsdk.WithResource(res),
		metricsdk.WithReader(
			metricsdk.NewPeriodicReader(
				exporter,
				metricsdk.WithInterval(10*time.Second),
			),
		),
	)

	down.OnShutdown(func() {
		ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
		defer cancel()
		if err := provider.Shutdown(ctx); err != nil {
			err = fmt.Errorf("cannot shutdown metric provider: %w", err)
			logger.Error(ctx, err.Error(), slog.Any(log.ErrorKey, err))
		}
	})

	return provider, nil
}

func newResource(ctx context.Context) (*resource.Resource, error) {
	res, err := resource.New(ctx, resource.WithAttributes(
		attribute.String("service_name", "demo"),
	))
	if err != nil {
		return nil, fmt.Errorf("failed to create telemetry resource: %w", err)
	}
	return res, nil
}

// setGlobalErrorHandler - there is a global OpenTelemetry logger.
// It is used for  errors during traces/metrics processing.
// Connect it to our logger.
func setGlobalErrorHandler(ctx context.Context, logger *log.Logger) {
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		logger.Error(ctx, err.Error(), slog.String(log.LoggerKey, "otel"), slog.Any(log.ErrorKey, err))
	}))
}
