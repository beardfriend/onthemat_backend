package infrastructure

import (
	"context"
	"testing"

	"onthemat/internal/app/config"
)

func TestRedisConnect(t *testing.T) {
	ctx := context.Background()
	c := config.NewConfig()
	if err := c.Load("../../../configs"); err != nil {
		t.Error(err)
	}
	redisCli := NewRedis(c)

	if err := redisCli.Set(ctx, "key", "value", 0).Err(); err != nil {
		t.Error(err)
	}

	res, _ := redisCli.Get(ctx, "key").Result()
	if res != "value" {
		t.Error("에러")
	}

	redisCli.Set(ctx, "key", "value", 1)
}
