package main

import (
	"fmt"
	"github.com/tfbrother/tfcache"
)

var (
	cache *tfcache.Tfcache = tfcache.NewTfCache()
	err   error
	value interface{}
	key   string
)

func main() {
	key = "tfbrother"
	if err = cache.Set(key, 22); err != nil {
		fmt.Println(err, "111")
	}
	if value, err = cache.Get("tfbrother"); err != nil {
		fmt.Println(err, "222")
	} else {
		fmt.Println("key:", key, "value:", value)
	}
}
