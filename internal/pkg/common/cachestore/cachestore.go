package cachestore

import (
	"context"

	"github.com/eko/gocache/lib/v4/store"
	redisstore "github.com/eko/gocache/store/redis/v4"
	"github.com/redis/go-redis/v9"
)

type Store struct {
	store store.StoreInterface
}

func New(client *redis.Client) *Store {
	return &Store{
		store: redisstore.NewRedis(client),
	}
}

func (s *Store) Store() store.StoreInterface {
	return s.store
}

func (s *Store) Invalidate(ctx context.Context, options ...store.InvalidateOption) error {
	return s.store.Invalidate(ctx, options...)
}
