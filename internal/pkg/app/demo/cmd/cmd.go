package cmd

import (
	"context"
	"fmt"
)

func Run(ctx context.Context) error {
	if srv, err := NewServer(ctx); err != nil {
		return fmt.Errorf("cannot create server: %w", err)
	} else if err := srv.Serve(ctx); err != nil {
		return fmt.Errorf("fatal server error: %w", err)
	}
	return nil
}
