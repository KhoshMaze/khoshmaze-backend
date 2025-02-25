package redis

import (
	"context"
	"errors"
	"time"

	c "github.com/KhoshMaze/khoshmaze-backend/internal/adapters/cache"

	"github.com/redis/go-redis/v9"
)

type redisCacheAdapter struct {
	client *redis.Client
}

func NewRedisProvider(redisAddr string, pass string, db int, clientName string) c.Provider {
	rdb := redis.NewClient(&redis.Options{
		Addr:       redisAddr,
		Password:   pass,
		DB:         db,
		ClientName: clientName,
	})

	return &redisCacheAdapter{
		client: rdb,
	}
}

func (r *redisCacheAdapter) Set(ctx context.Context, key string, ttl time.Duration, data []byte) error {
	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *redisCacheAdapter) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, c.ErrCacheMiss
		}
		return nil, err
	}

	return data, nil
}

func (r *redisCacheAdapter) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
