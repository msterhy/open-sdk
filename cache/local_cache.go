package cache

import (
	"sync"
	"time"
)

type CacheItem struct {
	Value    interface{}
	ExpireAt time.Time
}

type LocalCache struct {
	data  map[string]*CacheItem
	mutex sync.RWMutex
	ttl   time.Duration
}

// NewLocalCache 创建一个新的本地缓存
func NewLocalCache(ttl time.Duration) *LocalCache {
	return &LocalCache{
		data: make(map[string]*CacheItem),
		ttl:  ttl,
	}
}

// Set 设置缓存值
func (c *LocalCache) Set(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = &CacheItem{
		Value:    value,
		ExpireAt: time.Now().Add(c.ttl),
	}
}

// Get 获取缓存值
func (c *LocalCache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.data[key]
	if !exists || time.Now().After(item.ExpireAt) {
		return nil, false
	}
	return item.Value, true
}

// Delete 删除缓存
func (c *LocalCache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.data, key)
}

// Cleanup 定期清理过期缓存
func (c *LocalCache) Cleanup() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for key, item := range c.data {
		if time.Now().After(item.ExpireAt) {
			delete(c.data, key)
		}
	}
}
