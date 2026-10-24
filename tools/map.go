package tools

import "sync"

// SyncMap is a concurrent safe map that uses sync.Map internally.
type SyncMap[K comparable, V any] struct {
	m sync.Map
}

func (c *SyncMap[K, V]) Store(key K, value V) {
	c.m.Store(key, value)
}

func (c *SyncMap[K, V]) Load(key K) (V, bool) {
	v, ok := c.m.Load(key)
	if !ok {
		var zero V
		return zero, false
	}
	return v.(V), true
}

func (c *SyncMap[K, V]) Delete(key K) {
	c.m.Delete(key)
}

func (c *SyncMap[K, V]) Range(f func(key K, value V) bool) {
	c.m.Range(func(k, v any) bool {
		return f(k.(K), v.(V))
	})
}
