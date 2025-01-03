package redis

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/michaljurecko/ch-demo/internal/pkg/common/shutdown"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"github.com/michaljurecko/ch-demo/internal/pkg/common/log"
	"github.com/michaljurecko/ch-demo/internal/pkg/common/redis/config"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

const AddressKey = "address"

func NewClient(
	ctx context.Context,
	down *shutdown.Stack,
	logger *log.Logger,
	cfg config.Config,
	traceProvider trace.TracerProvider,
	meterProvider metric.MeterProvider,
) (*redis.Client, error) {
	logger = logger.With(slog.String(log.LoggerKey, "redis"))
	client := redis.NewClient(&redis.Options{
		Network:     "tcp",
		Addr:        cfg.Address,
		Username:    cfg.Password,
		Password:    cfg.Password,
		DB:          cfg.DB,
		ReadTimeout: 5 * time.Second,
	})

	if err := redisotel.InstrumentTracing(client, redisotel.WithTracerProvider(traceProvider)); err != nil {
		return nil, err
	}

	if err := redisotel.InstrumentMetrics(client, redisotel.WithMeterProvider(meterProvider)); err != nil {
		return nil, err
	}

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("cannot create redis client: ping failed: %w", err)
	}

	down.OnShutdown(func() {
		if err := client.Close(); err != nil {
			logger.Error(ctx, "failed to close redis client", slog.Any(log.ErrorKey, err))
			return
		}
		logger.Info(ctx, "redis client closed", slog.String(AddressKey, cfg.Address))
	})

	logger.Info(ctx, "connected to redis", slog.String(AddressKey, cfg.Address))

	return client, nil
}
