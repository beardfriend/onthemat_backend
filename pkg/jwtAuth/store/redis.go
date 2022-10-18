package store

import (
	"context"
	"time"

	auth "onthemat/pkg/jwtAuth"

	"github.com/go-redis/redis/v9"
)

type store struct {
	cli *redis.Client
}

func NewStore(client *redis.Client) auth.Store {
	return &store{
		cli: client,
	}
}

func (s *store) Set(ctx context.Context, tokenString string, expiration time.Duration) error {
	return s.cli.Set(ctx, tokenString, 1, expiration).Err()
}

func (s *store) Check(ctx context.Context, tokenString string) (bool, error) {
	val, err := s.cli.Exists(ctx, tokenString).Result()
	if err != nil {
		return false, err
	}

	if val > 0 {
		return true, err
	}

	return false, err
}
