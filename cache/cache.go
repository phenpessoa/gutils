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
	s := &Cache[K, V]{cache: make(map[K]item[V]), cacheFor: cacheFor}
	if cacheFor > 0 {
		go s.gc()
	}
	return s
}

// Cache is a simple Key/Value thread safe cache.
type Cache[K comparable, V any] struct {
	locker   sync.RWMutex
	cache    map[K]item[V]
	cacheFor time.Duration
}

type item[V any] struct {
	v V
	t time.Time
}

// Set sets a new value to the Cache cache.
func (s *Cache[K, V]) Set(k K, v V) {
	s.locker.Lock()
	defer s.locker.Unlock()
	s.cache[k] = item[V]{v, time.Now()}
}

// Get returns the value in the Cache cache of
// the passed key and if it was found or not.
func (s *Cache[K, V]) Get(k K) (V, bool) {
	s.locker.RLock()
	defer s.locker.RUnlock()
	item, ok := s.cache[k]
	return item.v, ok
}

// Delete deletes an entry from the Cache cache.
func (s *Cache[K, V]) Delete(k K) {
	s.locker.Lock()
	defer s.locker.Unlock()
	delete(s.cache, k)
}

// GetSet tries to find k in the cache.
// If found, GetSet will return the value found.
// If not found, GetSet will return the passed value and set it to the cache.
func (s *Cache[K, V]) GetSet(k K, v V) V {
	s.locker.Lock()
	defer s.locker.Unlock()
	val, ok := s.cache[k]
	if ok {
		return val.v
	}
	s.cache[k] = item[V]{v, time.Now()}
	return v
}

// TickGC runs the GC now.
func (s *Cache[K, V]) TickGC() {
	s.locker.Lock()
	for k, v := range s.cache {
		if time.Since(v.t) > s.cacheFor {
			delete(s.cache, k)
		}
	}
	s.locker.Unlock()
}

func (s *Cache[K, V]) gc() {
	ticker := time.NewTicker(s.cacheFor)
	for range ticker.C {
		s.TickGC()
	}
}
