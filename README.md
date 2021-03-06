# tfcache
key value cache, written by golang

## 版本
### V0.1(2018-11-12)
* 只支持GET/SET，不支持淘汰，不支持过期
* 采用map存储，性能很差
  - map里面加入新元素会导致动态的内存分配
  - 假设map里面有1亿+条数据的时候，频繁的加锁会导致性能差
  - 假设map里面1亿+条数据时，动态增加元素性能差
* ~~非线程安全的~~

### V0.11(2018-11-13)
* 缓存管理采用LRU算法，支持容量限制--参考Groupcache
* 支持相关统计信息(命中数，请求数，淘汰数，缓存总数)--参考Groupcache
* 提供http接口用于set/get cache

### V0.12(2018-11-14)
* 拆分Cache中的大map，引入一层Tfcache--参考bigcache的cacheshard
* key=>Cache对应，采用一致性hash算法，且引入虚拟节点--参考Groupcache

## 算法说明
### LRU算法
* 存储采用container/list中的双向链表
* 添加缓存时，放在链表的首部
* 获取缓存时，移动元素到链表首部
* 淘汰时就从链表的尾部开始来淘汰

### 一致性hash算法
* 所有的缓存节点以及其虚拟节点形成一个圆环，根据节点hash值顺序排列成圆环
* 对请求的key求hash值，找到在圆环中第一个大于该hash值的节点，就从该节点中取数据

## 参考设计

- [Bigcache](bigcache/bigcache.md)
- [Groupcache](https://github.com/golang/groupcache)