package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/eko/gocache/lib/v4/store"
	"google.golang.org/protobuf/proto"
)

type Marshaler[T proto.Message] struct {
	store store.StoreInterface
}

func NewMarshaler[T proto.Message](cache store.StoreInterface) *Marshaler[T] {
	return &Marshaler[T]{store: cache}
}

func (c *Marshaler[T]) Get(ctx context.Context, key string) (T, error) {
	result := *new(T)
	result = result.ProtoReflect().New().Interface().(T) //nolint:revive // type is T

	raw, err := c.store.Get(ctx, key)
	if err != nil {
		return result, err
	}

	rawBytes, ok := raw.(string)
	if !ok {
		return result, fmt.Errorf("unexpected type of cached value, expected []byte, given %T", raw)
	}

	if err := proto.Unmarshal([]byte(rawBytes), result); err != nil {
		return result, err
	}

	return result, nil
}

func (c *Marshaler[T]) GetWithTTL(ctx context.Context, key string) (T, time.Duration, error) {
	result := *new(T)
	result = result.ProtoReflect().New().Interface().(T) //nolint:revive // type is T

	raw, ttl, err := c.store.GetWithTTL(ctx, key)
	if err != nil {
		return *new(T), 0, err
	}

	rawBytes, ok := raw.(string)
	if !ok {
		return result, 0, fmt.Errorf("unexpected type of cached value, expected []byte, given %T", raw)
	}

	if err := proto.Unmarshal([]byte(rawBytes), result); err != nil {
		return result, 0, err
	}

	return result, ttl, nil
}

func (c *Marshaler[T]) Set(ctx context.Context, key string, object T, options ...store.Option) error {
	bytes, err := proto.Marshal(object)
	if err != nil {
		return err
	}

	return c.store.Set(ctx, key, string(bytes), options...)
}

func (c *Marshaler[T]) Delete(ctx context.Context, key string) error {
	return c.store.Delete(ctx, key)
}

func (c *Marshaler[T]) Invalidate(ctx context.Context, options ...store.InvalidateOption) error {
	return c.store.Invalidate(ctx, options...)
}

func (c *Marshaler[T]) Clear(ctx context.Context) error {
	return c.store.Clear(ctx)
}

func IsKeyNotFoundErr(err error) bool {
	var target *store.NotFound
	return errors.As(err, &target)
}
