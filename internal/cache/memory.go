package cache

import (
	"sync"
	"time"

	"go-exercise/internal/models"
)

type entry struct {
	value     float64
	expiresAt time.Time
}

type MemoryCache struct {
	mu   sync.RWMutex
	data map[models.Pair]entry
	ttl  time.Duration
	now  func() time.Time
}

func NewMemoryCache(ttl time.Duration) *MemoryCache {
	return &MemoryCache{
		data: make(map[models.Pair]entry),
		ttl:  ttl,
		now:  time.Now,
	}
}

func (c *MemoryCache) Get(p models.Pair) (float64, bool) {
	c.mu.RLock()
	e, ok := c.data[p]
	c.mu.RUnlock()
	if !ok {
		return 0, false
	}
	if c.now().After(e.expiresAt) {
		// expired
		c.mu.Lock()
		delete(c.data, p)
		c.mu.Unlock()
		return 0, false
	}
	return e.value, true
}

func (c *MemoryCache) Set(p models.Pair, v float64) {
	c.mu.Lock()
	c.data[p] = entry{
		value:     v,
		expiresAt: c.now().Add(c.ttl),
	}
	c.mu.Unlock()
}
