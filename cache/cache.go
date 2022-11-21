package cache

import (
	"sync"
	"time"
)

// NewCache creates a new Cache.
// cacheFor is the amount of time the value will be cached, set to 0 to cache forever.
// cacheFor is also the tick period of the GC.
// The caller can manually call the GC at any time using the TickGC method.
func NewCache[K comparable, V any](cacheFor time.Duration) *Cache[K, V] {
	c := &Cache[K, V]{cache: make(map[K]item[V]), cacheFor: cacheFor}
	if cacheFor > 0 {
		go c.gc()
	}
	return c
}

// Cache is a simple Key/Value thread safe cache.
type Cache[K comparable, V any] struct {
	locker   sync.RWMutex
	cache    map[K]item[V]
	cacheFor time.Duration
}

type item[V any] struct {
	t time.Time
	v V
}

// Set sets a new value to the Cache cache.
func (c *Cache[K, V]) Set(k K, v V) {
	c.locker.Lock()
	defer c.locker.Unlock()
	c.cache[k] = item[V]{time.Now(), v}
}

// Get returns the value in the Cache cache of
// the passed key and if it was found or not.
func (c *Cache[K, V]) Get(k K) (V, bool) {
	c.locker.RLock()
	defer c.locker.RUnlock()
	item, ok := c.cache[k]
	return item.v, ok
}

// Contains reports whether k is present in the cache.
func (c *Cache[K, V]) Contains(k K) bool {
	_, ok := c.Get(k)
	return ok
}

// Delete deletes an entry from the Cache cache.
func (c *Cache[K, V]) Delete(k K) {
	c.locker.Lock()
	defer c.locker.Unlock()
	delete(c.cache, k)
}

// GetSet tries to find k in the cache.
// If found, GetSet will return the value found.
// If not found, GetSet will return the passed value and set it to the cache.
func (c *Cache[K, V]) GetSet(k K, v V) V {
	c.locker.Lock()
	defer c.locker.Unlock()
	val, ok := c.cache[k]
	if ok {
		return val.v
	}
	c.cache[k] = item[V]{time.Now(), v}
	return v
}

// Wipe deletes all entries from the cache.
func (c *Cache[K, V]) Wipe() {
	c.locker.Lock()
	defer c.locker.Unlock()
	c.cache = make(map[K]item[V])
}

// Len returns the len of the cache.
func (c *Cache[K, V]) Len() int {
	c.locker.RLock()
	defer c.locker.RUnlock()
	return len(c.cache)
}

// TickGC runs the GC now.
// It will delete all expired entries
// from the cache.
func (c *Cache[K, V]) TickGC() {
	c.locker.Lock()
	for k, v := range c.cache {
		if time.Since(v.t) > c.cacheFor {
			delete(c.cache, k)
		}
	}
	c.locker.Unlock()
}

func (c *Cache[K, V]) gc() {
	ticker := time.NewTicker(c.cacheFor)
	for range ticker.C {
		c.TickGC()
	}
}
