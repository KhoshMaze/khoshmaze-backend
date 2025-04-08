package redis

import (
	"context"
	"errors"
	"time"


	"github.com/redis/go-redis/v9"
)

var (
	ErrCacheMiss = errors.New("redis cache miss")
	ErrCacheFull = errors.New("redis cache full")
	ErrCacheSet  = errors.New("redis cache set error")
)

type redisCacheAdapter struct {
	client *redis.Client
}

func NewRedisProvider(redisAddr string, pass string, db int, clientName string) *redisCacheAdapter {
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
			return nil, ErrCacheMiss
		}
		return nil, err
	}

	return data, nil
}

func (r *redisCacheAdapter) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *redisCacheAdapter) Exists(ctx context.Context, key string) (bool, error) {
	return r.client.Exists(ctx, key).Val() == 1, nil
}

func (r *redisCacheAdapter) Purge(ctx context.Context){
	r.client.FlushDB(ctx)
}
