package cache

import (
	"sync"
)

const (
	// DefaultTTL is the default time to live for cache entries in seconds.
	defaultTTL = 0
)

type Cache[K comparable, V any] struct {
	mu sync.RWMutex

	s map[K]*val[K, V]
}

func New[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{
		s: make(map[K]*val[K, V]),
	}
}

func (c *Cache[K, V]) Get(key K) V {
	c.mu.RLock()
	defer c.mu.RUnlock()

	val := c.s[key]
	if val == nil {
		return *new(V)
	}
	return val.v
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	v := newVal(key, value, defaultTTL)

	c.s[key] = v
}

func (c *Cache[K, V]) SetWithTTL(key K, value V, ttl int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	v := newVal(key, value, ttl)
	v.deleteFunc = c.delete

	c.s[key] = v
	go v.run()
}

func (c *Cache[K, V]) delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.s, key)
}
