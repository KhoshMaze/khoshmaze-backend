package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/KhoshMaze/khoshmaze-backend/config"
	appCtx "github.com/KhoshMaze/khoshmaze-backend/internal/adapters/context"
	memcache "github.com/KhoshMaze/khoshmaze-backend/internal/adapters/memory-cache"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/redis"
)

type FallbackCache struct {
	primary   Provider
	secondary Provider
	mu        sync.RWMutex
	isHealthy bool

	healthCheckInterval time.Duration
	maxErrors           int
}

func NewFallbackCache(r *config.RedisConfig, m *config.MemoryCacheConfig) Provider {
	fc := FallbackCache{
		primary:             redis.NewRedisProvider(fmt.Sprintf("%s:%d", r.Host, r.Port), r.Password, r.DB, r.ClientName),
		secondary:           memcache.NewMemoryCache(m.Size),
		healthCheckInterval: time.Second * m.HealthCheckInterval,
		maxErrors:           3,
	}

	go fc.startHealthCheck()

	return &fc
}

func (fc *FallbackCache) startHealthCheck() {
	ticker := time.NewTicker(fc.healthCheckInterval)
	defer ticker.Stop()
	errors := 0
	for range ticker.C {
		ctx := context.Background()
		logger := appCtx.GetLogger(ctx)
		testKey := "health-check"
		testValue := []byte("test")

		err := fc.primary.Set(ctx, testKey, time.Second*fc.healthCheckInterval, testValue)
		if err != nil {
			if errors < fc.maxErrors {
				logger.Error("REDIS IS UNHEALTHY", "error", err)
				fc.mu.Lock()
				fc.isHealthy = false
				fc.mu.Unlock()
			}
			errors++
			continue
		}

		_, err = fc.primary.Get(ctx, testKey)
		if err != nil {
			if errors < fc.maxErrors {
				logger.Error("REDIS IS UNHEALTHY", "error", err)
				fc.mu.Lock()
				fc.isHealthy = false
				fc.mu.Unlock()
			}
			errors++
			continue
		}

		// Redis is back, purge in-memory cache to release memory
		fc.mu.Lock()
		fc.isHealthy = true
		fc.secondary.Purge(ctx)
		fc.mu.Unlock()
		errors = 0
	}
}

func (fc *FallbackCache) currentProvider() Provider {
	// TODO: Read about if you should move in-memory cache to redis after it's healthy
	fc.mu.RLock()
	defer fc.mu.RUnlock()

	if fc.isHealthy {
		return fc.primary
	}
	return fc.secondary
}

func (fc *FallbackCache) Set(ctx context.Context, key string, ttl time.Duration, data []byte) error {
	return fc.currentProvider().Set(ctx, key, ttl, data)
}

func (fc *FallbackCache) Get(ctx context.Context, key string) ([]byte, error) {
	return fc.currentProvider().Get(ctx, key)
}

func (fc *FallbackCache) Del(ctx context.Context, key string) error {
	return fc.currentProvider().Del(ctx, key)
}

func (fc *FallbackCache) Exists(ctx context.Context, key string) (bool, error) {
	return fc.currentProvider().Exists(ctx, key)
}

func (fc *FallbackCache) Purge(ctx context.Context) {
	fc.currentProvider().Purge(ctx)
}
