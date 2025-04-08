package memcache

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	lru "github.com/hashicorp/golang-lru"
)

var (
	ErrCacheMiss = errors.New("memory cache miss")
	ErrCacheFull = errors.New("memory cache full")
	ErrCacheSet  = errors.New("memory cache set error")
)

type cachedItem struct {
	value     []byte
	ExpiresAt time.Time
}

type memoryCache struct {
	cache     *lru.Cache
	mu        sync.RWMutex
	ticker    *time.Ticker
	isActive  atomic.Bool
	cleanerMu sync.Mutex
	stopChan  chan struct{}
}

func NewMemoryCache(size int) *memoryCache {
	cache, err := lru.New(size)
	if err != nil {
		panic(err)
	}

	mc := &memoryCache{
		cache:    cache,
		stopChan: make(chan struct{}),
		ticker:   time.NewTicker(time.Second * 10),
	}
	go mc.startCleanup()

	return mc
}

func (m *memoryCache) startCleanup() {
	m.cleanerMu.Lock()
	defer m.cleanerMu.Unlock()

	if m.isActive.Load() {
		return
	}

	m.isActive.Store(true)
	m.stopChan = make(chan struct{})
	m.ticker = time.NewTicker(time.Second * 10)
	go m.checkCleanup()
}

func (m *memoryCache) stopCleanup() {
	m.cleanerMu.Lock()
	defer m.cleanerMu.Unlock()

	if !m.isActive.Load() {
		return
	}

	m.isActive.Store(false)
	close(m.stopChan)
	m.ticker.Stop()
}

func (m *memoryCache) checkCleanup() {
	select {
	case <-m.ticker.C:
		m.cleanup()
	case <-m.stopChan:
		return
	}
}

func (m *memoryCache) cleanup() {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	keys := m.cache.Keys()

	for _, key := range keys {
		item, ok := m.cache.Peek(key)
		if !ok {
			continue
		}

		if item.(cachedItem).ExpiresAt.Before(now) {
			m.cache.Remove(key)
		}
	}
	if keys != nil {
		if len(keys) == 0 {
			go m.stopCleanup()
		}
	} else {
		go m.stopCleanup()
	}
}

func (m *memoryCache) Set(ctx context.Context, key string, ttl time.Duration, data []byte) error {

	m.mu.Lock()
	defer m.mu.Unlock()
	m.cache.Add(key, cachedItem{
		value:     data,
		ExpiresAt: time.Now().Add(ttl),
	})

	go m.startCleanup()

	return nil
}

func (m *memoryCache) Get(ctx context.Context, key string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	item, ok := m.cache.Get(key)
	if !ok {
		return nil, ErrCacheMiss
	}
	return item.(cachedItem).value, nil
}

func (m *memoryCache) Del(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	ok := m.cache.Remove(key)
	if !ok {
		return ErrCacheMiss
	}
	return nil
}

func (m *memoryCache) Exists(ctx context.Context, key string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	exists := m.cache.Contains(key)
	if !exists {
		return false, ErrCacheMiss
	}
	return true, nil
}

func (m *memoryCache) Purge(ctx context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cache.Purge()
	go m.stopCleanup()
}
