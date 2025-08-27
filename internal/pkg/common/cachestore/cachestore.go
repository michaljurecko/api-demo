package cachestore

import (
	"context"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/eko/gocache/lib/v4/store"
	ristrettostore "github.com/eko/gocache/store/ristretto/v4"
)

type Store struct {
	store store.StoreInterface
}

func New() (*Store, error) {
	// configure Ristretto cache
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of
		MaxCost:     1 << 30, // maximum cost of cache (1GB)
		BufferItems: 64,      // number of keys per Get buffer
	})
	if err != nil {
		return nil, err
	}

	return &Store{
		store: ristrettostore.NewRistretto(cache, store.WithExpiration(15*time.Second)),
	}, nil
}

func (s *Store) Store() store.StoreInterface {
	return s.store
}

func (s *Store) Invalidate(ctx context.Context, options ...store.InvalidateOption) error {
	return s.store.Invalidate(ctx, options...)
}
