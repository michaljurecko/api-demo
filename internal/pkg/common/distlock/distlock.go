package distlock

import (
	"context"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

type Locker struct {
	locker *redislock.Client
}

type Lock struct {
	lock *redislock.Lock
}

func NewLocker(client *redis.Client) *Locker {
	return &Locker{locker: redislock.New(client)}
}

func (l *Locker) Lock(ctx context.Context, key string, ttl time.Duration) (*Lock, error) {
	key = "lock-" + key
	lock, err := l.locker.Obtain(ctx, key, ttl, &redislock.Options{})
	if err != nil {
		return nil, err
	}

	return &Lock{lock: lock}, nil
}

func (l *Lock) Refresh(ctx context.Context, ttl time.Duration) error {
	return l.lock.Refresh(ctx, ttl, &redislock.Options{})
}

func (l *Lock) Unlock(ctx context.Context) error {
	return l.lock.Release(ctx)
}
