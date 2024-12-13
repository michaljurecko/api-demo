package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/michaljurecko/ch-demo/internal/pkg/common/dataverse/metadata"

	"github.com/michaljurecko/ch-demo/internal/pkg/common/dataverse/entitygen"
	webapi2 "github.com/michaljurecko/ch-demo/internal/pkg/common/dataverse/webapi"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %s\n", err) //nolint:forbidigo
		os.Exit(1)
	}
}

func run() error {
	// Close application on Ctrl+C
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

	if len(os.Args) != 3 {
		return errors.New("specify two arguments: <package name> <target directory>")
	}
	pkgName := os.Args[1]
	targetDir := os.Args[2]

	config := webapi2.ConfigFromENV()
	if config.IsEmpty() {
		return errors.New("missing ENVs: 'DEMO_MODEL_{TENANT_ID,CLIENT_ID,CLIENT_SECRET,API_HOST}'")
	}

	client, err := webapi2.NewClient(ctx, config, http.DefaultClient)
	if err != nil {
		return fmt.Errorf("cannot create client: %w", err)
	}

	api := metadata.NewAPI(client)

	return entitygen.Generate(ctx, logger, api, pkgName, targetDir)
}
