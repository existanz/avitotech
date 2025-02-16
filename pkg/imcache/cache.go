package imcache

import (
	"sync"
	"time"
)

type Cache interface {
	// Set sets the value for the given key.
	Set(key string, value interface{})
	// Get returns the value for the given key.
	Get(key string) (interface{}, bool)
	// Delete deletes the value for the given key.
	Delete(key string)
}

type CacheItem struct {
	Value      interface{}
	Expiration int64
}

type InMemoryCache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
	ttl   time.Duration
}

func NewInMemoryCache(defaultTTL time.Duration) *InMemoryCache {
	cache := &InMemoryCache{
		items: make(map[string]CacheItem),
		ttl:   defaultTTL,
	}
	go cache.startCleanupJob()
	return cache
}

func (c *InMemoryCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	expiration := time.Now().Add(c.ttl).UnixNano()
	c.items[key] = CacheItem{
		Value:      value,
		Expiration: expiration,
	}
}

func (c *InMemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.items[key]
	if !found || time.Now().UnixNano() > item.Expiration {
		return nil, false
	}
	return item.Value, true
}

func (c *InMemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *InMemoryCache) startCleanupJob() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		c.cleanup()
	}
}

func (c *InMemoryCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now().UnixNano()
	for key, item := range c.items {
		if now > item.Expiration {
			delete(c.items, key)
		}
	}
}
