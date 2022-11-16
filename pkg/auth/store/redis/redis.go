package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
)

type store struct {
	cli *redis.Client
}

func NewStore(client *redis.Client) *store {
	return &store{
		cli: client,
	}
}

func (s *store) Del(ctx context.Context, key string) error {
	return s.cli.Del(ctx, key).Err()
}

func (s *store) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return s.cli.SetEx(ctx, key, value, expiration).Err()
}

func (s *store) Get(ctx context.Context, key string) string {
	return s.cli.Get(ctx, key).Val()
}

func (s *store) Check(ctx context.Context, key string) (bool, error) {
	val, err := s.cli.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if val > 0 {
		return true, err
	}

	return false, err
}
