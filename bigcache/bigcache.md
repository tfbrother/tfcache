### bigcache研究
#### 类图
![uml](https://github.com/tfbrother/tfcache/blob/master/bigcache/uml.png?raw=true)
#### 调用图
![call](https://github.com/tfbrother/tfcache/blob/master/bigcache/call.png?raw=true)
#### 主要功能
- Get/Set/淘汰
#### 缓存项存储格式
采用字节队列（数组）来存储缓存数据。一个缓存项在字节队列中对应的存储格式如下：
|headerEntrySize|timestamp|hashedKey|keyLength|key|entry|
| --- | --- | --- | --- | --- | --- |
| 4个字节 | 8个字节 | 8个字节 | 2个字节 | x个字节 | y个字节 |

* headerEntrySize
    存储的是后面几项的总长度
* timestamp
    缓存放入时的时间戳
* hashedKey
    key的hash值
* keyLength
    key的长度
* key
    key的内容
* entry
    value的内容
* 公式
    len(value)=总长度-8-8-2-len(key)
* 思考：为何没有单独存储len(value)？

### 字节队列BytesQueue
``` go
type BytesQueue struct {
	array        []byte
	capacity     int    //当前容量
	maxCapacity  int    //最大容量
	head         int    //头部位置
	tail         int    //尾部位置
	count        int    //缓存项数目
	rightMargin  int    //右边距，这个值以后面的字节才能使用来存储数据。
	headerBuffer []byte //4个字节的头部信息缓冲，方便重复使用这部分内存。
	verbose      bool
}
```
* tail >= head
此时rightMargin==tail，

| 0 | 1 | . | X | . | . | . | Y | . | . | . |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 0 | leftMarginIndex | . | head | . | . | . | tail/rightMargin | . | . | . |

* tail < head

| 0 | 1 | . | X | . | . | . | Y | . | Z | . |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 0 | leftMarginIndex | . | tail | . | . | . | head | . | rightMargin | . |

* rightMargin的含义
  作为字节队列，默认是从头开始依次插入缓存。假设总容量是100，当前tail是80，当插入一个长度为21的缓存项目时，没没法在尾部插入的（此时假设不能扩容），就需要从头部插入（假设头部剩余容量够），那么rightMargin就等于80。这部分时没法利用的，类似除不尽的余数。
  
### 字节队列的内存管理
* 初始化时申请一片内存，不够用再成本增长申请。
* 缓存淘汰算法
  根据最久最先淘汰算法，所有从head位置开始淘汰，淘汰后这部分内存会被后面添加的缓存继续使用，所以才会存在tail小于head的情况。
* 重复设置key时，对应的算法是？
  用resetKeyFromEntry函数先将旧内容全部字节设置为0，然后重新设置hashkey指向新内容所在的索引。因此这部分旧内容所占的内存就浪费了。