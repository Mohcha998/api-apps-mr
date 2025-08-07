package cache

import (
	"apps/config"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisConnection(cfg *config.Config) *redis.Client {
	dsn := fmt.Sprintf("%s:%s",
		cfg.Redis.Host,
		cfg.Redis.Port,
	)

	client := redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: cfg.Redis.Password,
	})

	return client
}
