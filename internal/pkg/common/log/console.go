package log

import (
	"log/slog"
	"os"
)

func newConsoleHandler() slog.Handler {
	return slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})
}
