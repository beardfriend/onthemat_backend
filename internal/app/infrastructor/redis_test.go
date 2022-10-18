package infrastructor

import (
	"context"
	"testing"

	"onthemat/pkg/test"
)

func TestRedisConnect(t *testing.T) {
	ctx := context.Background()

	test.BeforeStart()
	redisCli := NewRedis()

	if err := redisCli.Set(ctx, "key", "value", 0).Err(); err != nil {
		t.Error(err)
	}

	res, _ := redisCli.Get(ctx, "key").Result()
	if res != "value" {
		t.Error("에러")
	}

	redisCli.Set(ctx, "key", "value", 1)
}
