package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/michaljurecko/ch-demo/internal/pkg/app/demo/cmd"

	"github.com/michaljurecko/ch-demo/internal/pkg/common/log"
)

func main() {
	// Listen for SIGTERM/SIGINT to gracefully shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// Unexpected application termination is logged with a fallback logger
	if err := cmd.Run(ctx); err != nil {
		log.NewFallbackLogger().Error(ctx, err.Error())
	}
}
