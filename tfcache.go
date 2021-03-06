package tfcache

import (
	"container/list"
	"errors"
	"github.com/tfbrother/tfcache/consistenthash"
	"sync"
)

type Tfcache struct {
	caches      []*Cache
	maxShardNum int
	peers       *consistenthash.Map
}

// 设置缓存
func (b *Tfcache) Set(key string, value interface{}) (err error) {
	index := b.peers.Get(key)
	// debug
	//log.Println("key=", key, ",index=", index)
	return b.caches[index].Set(key, value)
}

// 获取缓存
func (b *Tfcache) Get(key string) (value interface{}, err error) {
	index := b.peers.Get(key)
	//log.Println("key=", key, ",index=", index)
	return b.caches[index].Get(key)
}
func (b *Tfcache) Stats(index int) CacheStats {
	if index < len(b.caches) && index >= 0 {
		return b.caches[index].Stats()
	}

	return CacheStats{}
}

func NewTfcache(shardNum int, maxNum int64) (b *Tfcache) {
	b = &Tfcache{
		peers:       consistenthash.New(10),
		maxShardNum: shardNum,                 //默认设置容量限制在10，设置得比较小，是为了测试方便
		caches:      make([]*Cache, shardNum), //slice必须要初始化
	}

	var peers []int
	for i := 0; i < shardNum; i++ {
		b.caches[i] = NewCache(maxNum / int64(shardNum))
		peers = append(peers, i)
	}

	b.peers.Add(peers)
	return
}

type item struct {
	key   string
	value interface{}
}

// 采用LRU算法
type Cache struct {
	cache              map[string]*list.Element
	mu                 sync.RWMutex
	num                int64 // 当前缓存中的key数量
	maxNum             int64 // 设置缓存中的key最大数量
	ll                 *list.List
	nhit, nget, nevict int64
}

// 设置缓存
func (c *Cache) Set(key string, value interface{}) (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// key已经存在，移动到链表头部
	if aa, ok := c.cache[key]; ok {
		c.ll.MoveToFront(aa)
		return
	}

	// 超过了容量限制，则删除该元素，同时淘汰掉链表末尾的元素
	if c.num >= c.maxNum {
		c.nevict++
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
	c.nget++
	if ele, ok := c.cache[key]; ok {
		c.nhit++
		c.ll.MoveToFront(ele)
		// 类型断言
		value = ele.Value.(*item).value
		return value, nil
	}
	return nil, errors.New("key not exist！！！")
}

func (c *Cache) Stats() CacheStats {
	c.mu.Lock()
	defer c.mu.Unlock()

	return CacheStats{
		Bytes:     0,
		Items:     c.num,
		Gets:      c.nget,
		Hits:      c.nhit,
		Evictions: c.nevict,
	}
}

func NewCache(maxNum int64) (tf *Cache) {
	tf = &Cache{
		cache:  make(map[string]*list.Element),
		ll:     list.New(),
		maxNum: maxNum, //默认设置容量限制在10，设置得比较小，是为了测试方便
	}

	return
}

// 缓存相关的统计信息
type CacheStats struct {
	Bytes     int64 // 内存占用
	Items     int64 // 缓存项目数量
	Gets      int64 // 请求数量
	Hits      int64 // 命中数量
	Evictions int64 // 淘汰数量
}
