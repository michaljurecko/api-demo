package service

import (
	"context"
	"log/slog"

	"github.com/eko/gocache/lib/v4/store"
	"github.com/michaljurecko/ch-demo/internal/pkg/common/log"
)

func (s *Service) invalidateCache(ctx context.Context, tags ...string) {
	if err := s.cacheStore.Invalidate(ctx, store.WithInvalidateTags(tags)); err != nil {
		s.logger.Error(ctx, "failed to invalidate cache", slog.Any(log.ErrorKey, err))
	}
}
