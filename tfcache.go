package tfcache

import "errors"

type Tfcache struct {
	cache map[string]interface{}
}

// 设置缓存
func (c *Tfcache) Set(key string, value interface{}) (err error) {
	// key已经存在
	if _, ok := c.cache[key]; ok {
		return errors.New("key has exist！！！")
	}
	c.cache[key] = value
	return nil
}

// 获取缓存
func (c *Tfcache) Get(key string) (value interface{}, err error) {
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
