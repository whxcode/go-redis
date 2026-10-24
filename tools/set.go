package tools

import "sync"

// SyncSet is a concurrent safe set that uses sync.Map internally.
type SyncSet[K comparable] struct {
	m sync.Map
}

func (c *SyncSet[K]) Has(key K) bool {
	_, ok := c.m.Load(key)
	return ok
}

func (c *SyncSet[K]) Store(key K) {
	c.m.Store(key, struct{}{})
}

func (c *SyncSet[K]) Delete(key K) {
	c.m.Delete(key)
}

func (c *SyncSet[K]) Range(f func(key K) bool) {
	c.m.Range(func(k, v any) bool {
		return f(k.(K))
	})
}

func (c *SyncSet[K]) Keys() (result []string) {
	c.m.Range(func(k, v any) bool {
		result = append(result, k.(string))
		return true
	})

	return result
}
