package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"threadly-api-v1/internal/config"
)

// NewRedisClient creates a Redis client and pings it so a bad REDIS_ADDR
// or an unreachable server fails at startup.
func NewRedisClient(ctx context.Context, redisConfig config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         redisConfig.Addr,
		Password:     redisConfig.Password,
		DB:           redisConfig.DB,
		PoolSize:     redisConfig.PoolSize,
		DialTimeout:  redisConfig.DialTimeout,
		ReadTimeout:  redisConfig.ReadTimeout,
		WriteTimeout: redisConfig.WriteTimeout,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return client, nil
}
