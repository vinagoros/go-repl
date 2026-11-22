package pokego

import (
	"sync"
	"time"
)

type Cache struct {
	ttl     time.Duration
	entries map[string]CacheEntry
	lock    sync.Mutex
}

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	c := Cache{ttl: interval}
	c.entries = make(map[string]CacheEntry)
	go c.reapLoop()
	return &c
}

func (c *Cache) Add(key string, val []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.entries[key] = CacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	entry, found := c.entries[key]
	if found {
		entry.createdAt = time.Now()
		c.entries[key] = entry
		return entry.val, found
	}
	return nil, found
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.ttl)
	defer ticker.Stop()
	for range ticker.C {
		c.lock.Lock()
		for key, entry := range c.entries {
			if time.Since(entry.createdAt) > c.ttl {
				delete(c.entries, key)
			}
		}
		c.lock.Unlock()
	}
}
