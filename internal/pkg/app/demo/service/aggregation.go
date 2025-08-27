package service

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"

	"github.com/eko/gocache/lib/v4/store"
	api "github.com/michaljurecko/ch-demo/api/gen/go/demo/v1"
	"github.com/michaljurecko/ch-demo/internal/pkg/app/demo/service/mapper"
	"github.com/michaljurecko/ch-demo/internal/pkg/common/cache"
	"github.com/michaljurecko/ch-demo/internal/pkg/common/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) ListPlayersAndCharacters(
	ctx context.Context,
	_ *connect.Request[emptypb.Empty],
) (*connect.Response[api.ListPlayersAndCharactersResponse], error) {
	// Load response from cache
	response, err := s.playersCache.Get(ctx, PlayersAndCharactersCacheKey)
	if err == nil {
		return connect.NewResponse(response), nil
	}

	// If there is no cached response, obtain a distributed lock, only one request/node is building cache
	// lock, err := s.locker.Lock(ctx, PlayersAndCharactersCacheKey, DistLockTTL)
	// if err != nil {
	// 	s.logger.Error(ctx, "failed to obtain distributed lock", slog.Any(log.ErrorKey, err))
	// } else {
	// 	defer func() {
	// 		if err := lock.Unlock(ctx); err != nil {
	// 			s.logger.Error(ctx, "failed to unlock distributed lock", slog.Any(log.ErrorKey, err))
	// 		}
	// 	}()
	// }

	// Try to get the cached response again, it might have been built by another request/node
	response, err = s.playersCache.Get(ctx, PlayersAndCharactersCacheKey)
	if err == nil {
		return connect.NewResponse(response), nil
	} else if !cache.IsKeyNotFoundErr(err) {
		s.logger.Error(ctx, "failed to load cache", slog.Any(log.ErrorKey, err))
	}

	// Generate response
	players, err := s.repo.Player().All().Do(ctx)
	if err != nil {
		return nil, err
	}
	characters, err := s.repo.Character().All().Do(ctx)
	if err != nil {
		return nil, err
	}
	response = &api.ListPlayersAndCharactersResponse{Players: mapper.PlayersAndCharacters(players, characters)}

	// Save response to the cache
	tags := []string{PlayersCacheTag, CharactersCacheTag}
	err = s.playersCache.Set(
		ctx,
		PlayersAndCharactersCacheKey,
		response,
		store.WithTags(tags),
		store.WithExpiration(CacheTTL),
	)
	if err != nil {
		s.logger.Error(ctx, "failed to save cache", slog.Any(log.ErrorKey, err))
	}

	return connect.NewResponse(response), nil
}
