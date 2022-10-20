package infrastructor

import (
	"fmt"

	"onthemat/internal/app/config"

	"github.com/go-redis/redis/v9"
)

func NewRedis(c *config.Config) *redis.Client {
	addr := fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	return rdb
}
