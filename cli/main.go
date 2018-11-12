package main

import (
	"fmt"
	"github.com/tfbrother/tfcache"
)

var (
	cache *tfcache.Cache = tfcache.NewCache()
	err   error
	value interface{}
)

func main() {
	for i := 0; i < 20; i++ {
		cache.Set(string(i), i)
	}
	for i := 0; i < 20; i++ {
		if value, err = cache.Get(string(i)); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("key:", string(i), "value:", value)
		}
	}
}
