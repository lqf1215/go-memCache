package main

import (
	"fmt"
	"go-menCache/cache_server"
	"time"
)

func main() {

	/**
	  使⽤示例
	*/
	cache := cache_server.NewMenCache()
	cache.SetMaxMemory("300GB")

	cache.Set("int", 1, time.Second)
	cache.Set("bool", false, time.Second)
	cache.Set("data", map[string]interface{}{"a": 1}, time.Second)

	cache.Set("int", 1)
	cache.Set("bool", false)
	cache.Set("data", map[string]interface{}{"a": 1})
	cache.Get("int")
	cache.Del("int")
	cache.Flush()
	fmt.Println(cache.Keys())

	//cache := cache_server.NewMenCache()
	cache.SetMaxMemory("300GB")

	cache.Set("int", 1, time.Second)
	cache.Set("bool", false, time.Second)
	cache.Set("data", map[string]interface{}{"a": 1}, time.Second)

	cache.Set("int", 1)
	cache.Set("bool", false)
	cache.Set("data", map[string]interface{}{"a": 1})
	cache.Get("int")
	fmt.Println(cache.Keys())
	fmt.Println(cache.Get("int"))
	fmt.Println(cache.Get("bool"))
}
