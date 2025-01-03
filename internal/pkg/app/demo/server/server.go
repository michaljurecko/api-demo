package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/michaljurecko/ch-demo/internal/pkg/app/demo/service"

	"github.com/michaljurecko/ch-demo/api"
	"github.com/michaljurecko/ch-demo/internal/pkg/app/demo/server/static"

	"connectrpc.com/connect"
	connectcors "connectrpc.com/cors"
	"connectrpc.com/otelconnect"
	"github.com/michaljurecko/ch-demo/internal/pkg/common/protovalidate"

	"github.com/michaljurecko/ch-demo/api/gen/go/demo/v1/apiconnect"
	"github.com/michaljurecko/ch-demo/internal/pkg/app/demo/server/config"
	"github.com/michaljurecko/ch-demo/internal/pkg/common/log"
	"github.com/michaljurecko/ch-demo/internal/pkg/common/shutdown"
	"github.com/rs/cors"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const ReadHeaderTimeout = 30 * time.Second

type Server struct {
	logger   *log.Logger
	down     *shutdown.Stack
	listener net.Listener
	server   *http.Server
}

func New(
	cfg config.Config,
	svc *service.Service,
	logger *log.Logger,
	down *shutdown.Stack,
	tracerProvider trace.TracerProvider,
	meterProvider metric.MeterProvider,
) (*Server, error) {
	srv := &Server{
		logger: logger.With(slog.String(log.LoggerKey, "http.server")),
		down:   down,
	}

	// Setup telemetry
	telemetryHandler, err := otelconnect.NewInterceptor(
		otelconnect.WithMeterProvider(meterProvider),
		otelconnect.WithTracerProvider(tracerProvider),
	)
	if err != nil {
		err = fmt.Errorf("failed to create telemetry interceptor: %w", err)
		return nil, err
	}

	// Setup validation
	validator, err := protovalidate.NewValidator()
	if err != nil {
		return nil, fmt.Errorf("failed to create proto validator: %w", err)
	}
	validatorHandler := protovalidate.NewInterceptor(validator, protovalidate.UpdateMaskPlugin)

	// Create router
	mux := http.NewServeMux()

	// Create generated handler
	path, handler := apiconnect.NewApiServiceHandler(
		svc,
		connect.WithInterceptors(telemetryHandler, validatorHandler),
	)

	// Add middlewares
	handler = withCORS(handler)

	// Connect handler to the router
	mux.Handle(path, handler)

	// Static routers
	mux.Handle("/", static.FS())
	mux.Handle("/openapi.yaml", api.GenFS())

	// Create HTTP server
	srv.server = &http.Server{
		Addr:              cfg.ListenAddress,
		ReadHeaderTimeout: ReadHeaderTimeout,
		Handler:           h2c.NewHandler(mux, &http2.Server{}), // HTTP2 without TLS,
	}

	return srv, nil
}

func (s *Server) Addr() net.Addr {
	return s.listener.Addr()
}

func (s *Server) Serve(ctx context.Context) error {
	// Connect request contexts to the parent context
	// There is a graceful shutdown, so the context shouldn't be canceled with the parent context.
	s.server.BaseContext = func(net.Listener) context.Context {
		return context.WithoutCancel(ctx)
	}

	// Setup graceful shutdown
	var wg sync.WaitGroup
	wg.Add(1)
	s.down.OnShutdown(func() {
		defer wg.Done()

		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 30*time.Second)
		defer cancel()

		if err := s.server.Shutdown(shutdownCtx); err != nil {
			err = fmt.Errorf("shutdown error: %w", err)
			s.logger.Error(ctx, err.Error(), slog.Any(log.ErrorKey, err))
			return
		}

		s.logger.Info(ctx, "shutdown complete")
	})

	// Start listener
	listener, err := s.Listen()
	if err != nil {
		return err
	}

	// Start server
	s.logger.Info(ctx, "starting HTTP server", slog.String("addr", s.listener.Addr().String()))
	if err = s.server.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
		err = fmt.Errorf("unexpected server error: %w", err)
		s.logger.Error(ctx, err.Error(), slog.Any(log.ErrorKey, err))
		return err
	}

	// ListenAndServe method ends when the server shutdown starts
	s.logger.Info(ctx, "stopped serving new connections")

	// Terminating the HTTP server triggers shutdown of the entire app (if it hasn't already happened)
	s.down.Shutdown()

	// Wait for graceful shutdown
	wg.Wait()
	return nil
}

// Listen creates a listener for the server, but does not start it.
// It is useful for test, when listening on a random port is needed.
// In the production, call only the Serve method.
func (s *Server) Listen() (listener net.Listener, err error) {
	if s.listener == nil {
		if s.listener, err = net.Listen("tcp", s.server.Addr); err != nil {
			return nil, fmt.Errorf("cannot create listener for HTTP server: %w", err)
		}
	}
	return s.listener, nil
}

func withCORS(next http.Handler) http.Handler {
	return cors.
		New(cors.Options{
			AllowedMethods: connectcors.AllowedMethods(),
			AllowedHeaders: connectcors.AllowedHeaders(),
			ExposedHeaders: connectcors.ExposedHeaders(),
		}).
		Handler(next)
}
