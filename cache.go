package go_cache

import (
	"go_cache/lru"
	"sync"
)

type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	//如果没有实例化对象,创建一个实例化对象
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	//否则添加
	c.add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.GET(key); ok {
		return v.(ByteView), ok
	}
	return
}
