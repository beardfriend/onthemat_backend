package store

import (
	"context"
	"time"
)

type Store interface {
	Del(ctx context.Context, key string) error
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Get(ctx context.Context, key string) string
	Check(ctx context.Context, key string) (bool, error)
}
