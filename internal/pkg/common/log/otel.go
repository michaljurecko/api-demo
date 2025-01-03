package log

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/michaljurecko/ch-demo/internal/pkg/common/log/config"
	"github.com/michaljurecko/ch-demo/internal/pkg/common/shutdown"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	logsdk "go.opentelemetry.io/otel/sdk/log"
)

func newOTELHandler(ctx context.Context, down *shutdown.Stack, cfg config.Config) (slog.Handler, error) {
	exporter, err := otlploghttp.New(ctx,
		otlploghttp.WithEndpointURL(cfg.HTTPEndpointURL),
		otlploghttp.WithHeaders(map[string]string{
			"Authorization": cfg.HTTPAuthorization,
			"stream-name":   "default",
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create 'otlploghttp' log exporter: %w", err)
	}

	provider := logsdk.NewLoggerProvider(
		logsdk.WithProcessor(
			logsdk.NewBatchProcessor(exporter),
		),
	)

	down.OnShutdown(func() {
		ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
		defer cancel()
		if err := provider.Shutdown(ctx); err != nil {
			err = fmt.Errorf("cannot shutdown log provider: %w", err)
			NewFallbackLogger().Error(ctx, err.Error(), slog.Any(ErrorKey, err))
		}
	})

	// Create slog handler
	return otelslog.NewHandler("demo", otelslog.WithLoggerProvider(provider)), nil
}
