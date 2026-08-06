package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/km-saifullah/go-erp/internal/config"
)

type Redis struct {
	Client *redis.Client
}

func NewRedis(cfg config.Config) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx := context.Background()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Redis{
		Client: client,
	}, nil
}

func (r *Redis) Close() error {
	return r.Client.Close()
}
