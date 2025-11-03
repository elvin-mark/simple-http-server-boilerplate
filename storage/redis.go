package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"http-server/config"
)

// RedisClient represents the Redis client.
type RedisClient struct {
	*redis.Client
}

// NewRedisClient initializes and returns a new Redis client.
func NewRedisClient(cfg *config.RedisConfig) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("could not connect to Redis: %w", err)
	}

	return &RedisClient{rdb}, nil
}
