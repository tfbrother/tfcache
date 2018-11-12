# tfcache
key value cache, written by golang

## 版本
### V0.1(2018-11-12)
* 只支持GET/SET，不支持淘汰，不支持过期
* 采用map存储，性能很差
  - map里面加入新元素会导致动态的内存分配
  - 假设map里面有100万+条数据的时候，频繁的加锁会导致性能差
  - 假设map里面100万+条数据时，动态增加元素性能差
* ~~非线程安全的~~

### V0.11(2018-11-12)
* 缓存管理采用LRU算法，支持容量限制


## 算法说明
### LRU算法
* 存储采用container/list中的双向链表
* 添加缓存时，放在链表的首部
* 获取缓存时，移动元素到链表首部
* 淘汰时就从链表的尾部开始来淘汰

## 参考设计

- [Bigcache](bigcache/bigcache.md)
- [Groupcache](https://github.com/golang/groupcache)