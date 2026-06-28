package cache

import (
	"sync"
	"time"
)

// item 缓存条目，包含值和过期时间
type item struct {
	value     any
	expiresAt time.Time
}

// LocalCache 本地内存缓存，基于 sync.Map 实现并发安全
type LocalCache struct {
	m sync.Map
}

// NewLocalCache 创建本地缓存实例
func NewLocalCache() *LocalCache {
	return &LocalCache{}
}

// Set 写入缓存，ttl 为过期时间
func (c *LocalCache) Set(key string, value any, ttl time.Duration) {
	c.m.Store(key, item{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	})
}

// Get 读取缓存，未命中或已过期返回 nil
func (c *LocalCache) Get(key string) (any, bool) {
	val, ok := c.m.Load(key)
	if !ok {
		return nil, false
	}
	it := val.(item)
	if time.Now().After(it.expiresAt) {
		c.m.Delete(key) // 惰性删除：过期了就删掉
		return nil, false
	}
	return it.value, true
}

// Delete 删除缓存
func (c *LocalCache) Delete(key string) {
	c.m.Delete(key)
}
