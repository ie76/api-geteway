package cache

import (
	"assignment/config"
	"assignment/internal/errors"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
}

type RedisCache struct {
	client *redis.Client
}

func NewCacheProvider() (*RedisCache, *errors.Error) {
	cfg := config.GetConfig()

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.REDIS_HOST, cfg.REDIS_PORT),
		Password: cfg.REDIS_PASSWORD,
		DB:       cfg.REDIs_DATABASE,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, errors.New(errors.ErrRedisConnect)
	}

	return &RedisCache{
		client: client,
	}, nil
}

func (rc *RedisCache) Get(ctx context.Context, key string) ([]byte, *errors.Error) {
	data, err := rc.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, errors.New(errors.ErrRedisConnect)
	}
	return data, nil
}

func (rc *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) *errors.Error {
	err := rc.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return errors.New(errors.ErrRedisConnect)
	}
	return nil
}
