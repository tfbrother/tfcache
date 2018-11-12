package tfcache

import (
	"errors"
	"sync"
)

type Tfcache struct {
	cache map[string]interface{}
	mu    sync.RWMutex
}

// 设置缓存
func (c *Tfcache) Set(key string, value interface{}) (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// key已经存在
	if _, ok := c.cache[key]; ok {
		return errors.New("key has exist！！！")
	}
	c.cache[key] = value
	return nil
}

// 获取缓存
func (c *Tfcache) Get(key string) (value interface{}, err error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var ok bool
	if value, ok = c.cache[key]; ok {
		return value, nil
	}
	return nil, errors.New("key not exist！！！")
}

func NewTfCache() (tf *Tfcache) {
	tf = &Tfcache{
		cache: make(map[string]interface{}),
	}

	return
}
