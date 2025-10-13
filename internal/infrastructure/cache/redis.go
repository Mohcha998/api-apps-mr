package cache

import (
	"apps/config"
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	conn *redis.Client
}

func NewRedisConnection(cfg *config.Config) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	return &Client{conn: rdb}
}

// ===================================================
// Basic string cache operations
// ===================================================

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.conn.Get(ctx, key).Result()
}

func (c *Client) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return c.conn.Set(ctx, key, value, ttl).Err()
}

func (c *Client) Delete(ctx context.Context, key string) error {
	return c.conn.Del(ctx, key).Err()
}

// ===================================================
// JSON helper methods
// ===================================================

// GetJSON mengambil data JSON dari Redis dan unmarshal ke target struct
func (c *Client) GetJSON(ctx context.Context, key string, target interface{}) error {
	val, err := c.conn.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), target)
}

// SetJSON menyimpan struct ke Redis dalam format JSON
func (c *Client) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.conn.Set(ctx, key, data, ttl).Err()
}
