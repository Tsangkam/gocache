package gocache

import (
	"gocache/gocache/lru"
	"sync"
)

// cache block
type block struct {
	mu         sync.RWMutex
	lru        *lru.Cache
	cacheBytes int64
}

func (c *block) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *block) get(key string) (value ByteView, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
