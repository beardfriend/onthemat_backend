package infrastructor

import (
	"fmt"

	"onthemat/internal/app/config"

	"github.com/go-redis/redis/v9"
)

func NewRedis() *redis.Client {
	config := config.Info
	addr := fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	return rdb
}
