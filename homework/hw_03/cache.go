package cache

import (
	"sync"
)

type Cache[K comparable, V any] struct {
	mu sync.RWMutex

	s map[K]*val[V]
}

func New[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{
		s: make(map[K]*val[V]),
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

// SetWithTTL устанавливает значение в кэше с заданным TTL (в секундах).
// Если TTL <= 0, значение будет храниться в кэше бесконечно.
func (c *Cache[K, V]) Set(key K, value V, ttl int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	v := newVal(value, ttl)
	v.deleteFunc = c.delete(key)

	c.s[key] = v
	go v.run()
}

func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.s, key)
}

func (c *Cache[K, V]) delete(key K) func() {
	return func() {
		delete(c.s, key)
	}
}
