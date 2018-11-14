package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tfbrother/tfcache"
	"net/http"
	"strconv"
)

var (
	cache *tfcache.Tfcache = tfcache.NewTfcache(10, 20000000) //限制数量
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

// 初始化cache
func init() {
	for i := 0; i < 20000000; i++ {
		// 整型转字符串，不能直接使用string(i)，此时i会被当成ascii来对待
		cache.Set(strconv.Itoa(i)+"tfbrother", i)
	}
}

func main() {
	go func() {
		for i := 0; i < 20000000; i++ {
			cache.Get(strconv.Itoa(i) + "tfbrother")
		}
	}()

	http.HandleFunc("/stats", handleStats)
	http.HandleFunc("/setCache", handleSet)
	http.HandleFunc("/getCache", handleGet)

	http.ListenAndServe("0.0.0.0:7777", nil)
}

// http://127.0.0.1:7777/stats?index=1
func handleStats(resp http.ResponseWriter, req *http.Request) {
	var (
		index int
		err   error
	)

	if index, err = strconv.Atoi(req.URL.Query().Get("index")); err == nil {
		stats := cache.Stats(index)
		resp.Write([]byte(ToString(&stats)))
		return
	}

	resp.Write([]byte("统计信息不存在"))
}

// 获取缓存
// curl http://localhost:7777/getCache?key=tfbrother
func handleGet(resp http.ResponseWriter, req *http.Request) {
	var (
		err   error
		key   string
		value interface{}
	)

	key = req.URL.Query().Get("key")
	if value, err = cache.Get(key); err == nil {
		//类型断言
		resp.Write([]byte(value.(string)))
	} else {
		resp.Write([]byte(err.Error()))
	}

}

// 设置缓存
// curl http://localhost:7777/setCache?key=tfbrother&value=22
func handleSet(resp http.ResponseWriter, req *http.Request) {
	var (
		err        error
		key, value string
	)

	key = req.URL.Query().Get("key")
	value = req.URL.Query().Get("value")
	if err = cache.Set(key, value); err == nil {
		resp.Write([]byte("sucess!"))
	} else {
		resp.Write([]byte(err.Error()))
	}
}
