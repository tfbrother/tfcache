package main

import (
	"fmt"
	"github.com/tfbrother/tfcache"
	"strconv"
)

var (
	cache *tfcache.Cache = tfcache.NewCache(10) //限制缓存的最大数量
	err   error
	value interface{}
)

func main() {
	for i := 0; i < 20; i++ {
		// 整型转字符串，不能直接使用string(i)，此时i会被当成ascii来对待
		cache.Set(strconv.Itoa(i), i)
	}
	for i := 0; i < 20; i++ {
		if value, err = cache.Get(strconv.Itoa(i)); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("key:", strconv.Itoa(i), "value:", value)
		}
	}
}
