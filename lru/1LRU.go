package lru

import (
	"container/list"
)

type Cache struct {
	MaxBytes  int64
	NowBytes  int64
	List      *list.List
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func (c *Cache) Len() int {
	return c.List.Len()
}

// 创建实例
func New(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		MaxBytes:  maxBytes,
		NowBytes:  0,
		List:      list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// 查找
func (c *Cache) GET(key string) (value Value, ok bool) {
	//如果查到了
	if ele, ok := c.cache[key]; ok {
		c.List.MoveToFront(ele)  //将该元素移到队尾
		kv := ele.Value.(*entry) //获取该元素的value,断言一下
		return kv.value, ok
	}
	return
}

func (c *Cache) RemoveOldest() {
	ele := c.List.Back() //获得对头元素
	if ele != nil {
		c.List.Remove(ele) //移除对头元素
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.NowBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	//如果当前key存在就修改值
	if ele, ok := c.cache[key]; ok {
		c.List.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.NowBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.List.PushFront(&entry{key: key, value: value})
		c.cache[key] = ele
		c.NowBytes += int64(len(key)) + int64(value.Len()) //加上map和链表的内存占值
	}
	//如果超出了最大内存,循环淘汰对头元素
	for c.MaxBytes != 0 && c.MaxBytes < c.NowBytes {
		c.RemoveOldest()
	}
}
