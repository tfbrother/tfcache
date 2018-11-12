package tfcache

import (
	"container/list"
	"errors"
	"sync"
)

type item struct {
	key   string
	value interface{}
}

// 采用LRU算法
type Cache struct {
	cache  map[string]*list.Element
	mu     sync.RWMutex
	num    int // 当前缓存中的key数量
	maxNum int // 设置缓存中的key最大数量
	ll     *list.List
}

// 设置缓存
func (c *Cache) Set(key string, value interface{}) (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// key已经存在，移动到链表头部
	if _, ok := c.cache[key]; ok {
		var aa *list.Element
		c.ll.MoveToFront(aa)
	}

	// 超过了容量限制，则删除该元素，同时淘汰掉链表末尾的元素
	if c.num >= c.maxNum {
		c.num--
		ele := c.ll.Remove(c.ll.Back())
		k := ele.(*item).key
		// 链表中必须要把缓存的key存下来，否则无法通过list中的元素找到对应缓存的key，实现删除map cache中对应key的功能
		delete(c.cache, k)
	}

	// 放入头部
	ele := c.ll.PushFront(&item{key, value})
	c.num++
	c.cache[key] = ele
	return nil
}

// 获取缓存
func (c *Cache) Get(key string) (value interface{}, err error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		// 类型断言
		value = ele.Value.(*item).value
		return value, nil
	}
	return nil, errors.New("key not exist！！！")
}

func NewCache(maxNum int) (tf *Cache) {
	tf = &Cache{
		cache:  make(map[string]*list.Element),
		ll:     list.New(),
		maxNum: maxNum, //默认设置容量限制在10，设置得比较小，是为了测试方便
	}

	return
}
