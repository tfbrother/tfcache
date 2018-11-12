# tfcache
key value cache, written by golang

## 版本问题
### V0.1(2018-11-12)
* 只支持GET/SET，不支持淘汰，不支持过期
* 采用map存储，性能很差
  - map里面加入新元素会导致动态的内存分配
  - 假设map里面有100万+条数据的时候，频繁的加锁会导致性能差
  - 假设map里面100万+条数据时，动态增加元素性能差
* 非线程安全的

[参考设计](bigcache/bigcache.md)