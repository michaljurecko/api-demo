package test

import (
	"context"
	"net/http"
	"time"

	"github.com/michaljurecko/ch-demo/internal/pkg/app/demo/cmd"
	"github.com/michaljurecko/ch-demo/internal/pkg/app/demo/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func StartTestServer(ctx context.Context) string {
	cfgFn := func(cfg config.Config) (config.Config, error) {
		cfg.Server.ListenAddress = "localhost:0" // listen on a random port
		return cfg, nil
	}

	// Create new background context
	ctx, cancel := context.WithCancel(context.WithoutCancel(ctx))
	DeferCleanup(func() {
		cancel()
	})

	// Create server
	srv, err := cmd.NewServerForTest(ctx, cfgFn)
	Expect(err).NotTo(HaveOccurred())

	// Listen on a random port
	listener, err := srv.Listen()
	Expect(err).NotTo(HaveOccurred())

	// Serve in background
	go func() {
		err := srv.Serve(ctx)
		Expect(err).NotTo(HaveOccurred())
	}()

	// Compose base URL
	baseURL := "http://" + listener.Addr().String()

	// Wait for server to start
	err = WaitForHTTP(ctx, baseURL, 10*time.Second)
	Expect(err).NotTo(HaveOccurred())

	return baseURL
}

func WaitForHTTP(ctx context.Context, url string, timeout time.Duration) (err error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			if err == nil {
				err = ctx.Err()
			}
			return err
		case <-ticker.C:
			var req *http.Request
			req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				return err
			}

			var resp *http.Response
			resp, err = http.DefaultClient.Do(req)
			if err == nil {
				_ = resp.Body.Close()
				return nil
			}
		}
	}
}
