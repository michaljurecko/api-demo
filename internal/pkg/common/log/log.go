// Package log provides logging interface based on slog.
// The interface requires context to be passed to each logging method.
// Arguments for structured logging can be passed only as slog.Attr type.
package log

import (
	"context"
	"log/slog"

	"github.com/michaljurecko/ch-demo/internal/pkg/common/log/config"
	"github.com/michaljurecko/ch-demo/internal/pkg/common/shutdown"
	slogmulti "github.com/samber/slog-multi"
)

type Logger struct {
	logger *slog.Logger
}

func NewLogger(ctx context.Context, down *shutdown.Stack, cfg config.Config) (*Logger, error) {
	var handlers []slog.Handler

	// Always log to console
	handlers = append(handlers, newConsoleHandler())

	// Use OpenTelemetry HTTP exporter, if enabled
	if cfg.Exporter == config.HTTPExporter {
		h, err := newOTELHandler(ctx, down, cfg)
		if err != nil {
			return nil, err
		}

		handlers = append(handlers, h)
	}

	var handler slog.Handler
	if len(handlers) == 1 {
		handler = handlers[0]
	} else {
		handler = slogmulti.Fanout(handlers...)
	}

	logger := wrapLogger(handler)

	return logger, nil
}

// NewFallbackLogger without dependencies, it can be used if the main logger fails to initialize.
func NewFallbackLogger() *Logger {
	return wrapLogger(newConsoleHandler())
}

func wrapLogger(handler slog.Handler) *Logger {
	return &Logger{logger: slog.New(handler)}
}

// Enabled reports whether l emits log records at the given context and level.
func (l *Logger) Enabled(ctx context.Context, level slog.Level) bool {
	return l.logger.Enabled(ctx, level)
}

// With returns a Logger that includes the given attributes
// in each output operation. Arguments are converted to
// attributes as if by [Logger.Log].
func (l *Logger) With(args ...slog.Attr) *Logger {
	var argsAny []any
	for _, a := range args {
		argsAny = append(argsAny, a)
	}
	return &Logger{logger: l.logger.With(argsAny...)}
}

// Debug logs at [LevelDebug] with the given context.
func (l *Logger) Debug(ctx context.Context, msg string, args ...slog.Attr) {
	l.logger.LogAttrs(ctx, slog.LevelDebug, msg, args...) //nolint:sloglint
}

// Info logs at [LevelInfo] with the given context.
func (l *Logger) Info(ctx context.Context, msg string, args ...slog.Attr) {
	l.logger.LogAttrs(ctx, slog.LevelInfo, msg, args...) //nolint:sloglint
}

// Warn logs at [LevelWarn] with the given context.
func (l *Logger) Warn(ctx context.Context, msg string, args ...slog.Attr) {
	l.logger.LogAttrs(ctx, slog.LevelWarn, msg, args...) //nolint:sloglint
}

// ErrorContext logs at [LevelError] with the given context.
func (l *Logger) Error(ctx context.Context, msg string, args ...slog.Attr) {
	l.logger.LogAttrs(ctx, slog.LevelError, msg, args...) //nolint:sloglint
}
