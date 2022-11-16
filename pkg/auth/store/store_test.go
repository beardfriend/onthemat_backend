package store

import (
	"context"
	"fmt"
	"testing"
	"time"

	"onthemat/internal/app/config"
	"onthemat/internal/app/infrastructure"
	"onthemat/pkg/auth/store/redis"

	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	c := config.NewConfig()
	err := c.Load("../../../configs")
	assert.NoError(t, err)
	redisClient := infrastructure.NewRedis(c)
	store := redis.NewStore(redisClient)

	if err := store.Set(context.Background(), "email", "value", time.Duration(time.Second*24)); err != nil {
		t.Error(err)
	}
}

func TestGet(t *testing.T) {
	c := config.NewConfig()
	err := c.Load("../../../configs")
	assert.NoError(t, err)
	redisClient := infrastructure.NewRedis(c)
	store := redis.NewStore(redisClient)

	asd := store.Get(context.Background(), "email")
	fmt.Println(asd)
}

func TestHSet(t *testing.T) {
	c := config.NewConfig()
	err := c.Load("../../../configs")
	assert.NoError(t, err)
	redisClient := infrastructure.NewRedis(c)
	store := redis.NewStore(redisClient)
	if err := store.HSet(context.Background(), "1", "asd", "asd", time.Duration(time.Second*24)); err != nil {
		t.Error(err)
	}
}

func TestHGet(t *testing.T) {
	c := config.NewConfig()
	err := c.Load("../../../configs")
	assert.NoError(t, err)
	redisClient := infrastructure.NewRedis(c)
	store := redis.NewStore(redisClient)
	val, _ := store.HGet(context.Background(), "1", "asd")
	assert.Equal(t, val, "asd")
}
