package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tfbrother/tfcache"
	"strconv"
	"sync"
)

var (
	cache *tfcache.Cache = tfcache.NewCache(100) //限制缓存的最大数量
	err   error
	value interface{}
)

//输出结构体
func ToString(conf interface{}) string {
	b, err := json.Marshal(conf)
	if err != nil {
		return fmt.Sprintf("%+v", conf)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", conf)
	}
	return out.String()
}

func main() {
	var wait sync.WaitGroup
	wait.Add(2)
	go func() {
		for i := 0; i < 200; i++ {
			// 整型转字符串，不能直接使用string(i)，此时i会被当成ascii来对待
			cache.Set(strconv.Itoa(i), i)
		}
		wait.Done()
	}()

	go func() {
		for i := 0; i < 200; i++ {
			if value, err = cache.Get(strconv.Itoa(i)); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("key:", strconv.Itoa(i), "value:", value)
			}
		}
		wait.Done()
	}()
	wait.Wait()
	stats := cache.Stats()
	fmt.Println(ToString(&stats))
}
