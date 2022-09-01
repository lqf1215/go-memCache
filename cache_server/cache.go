package cache_server

import (
	"go-menCache/cache"
	"time"
)

type cacheServer struct {
	memCache cache.Cache
}

func NewMenCache() *cacheServer {
	return &cacheServer{memCache: cache.NewMenCache()}
}

var Cache = NewMenCache()

//SetMaxMemory size 是⼀个字符串。⽀持以下参数: 1KB，100KB，1MB，2MB，1GB 等
func (cs *cacheServer) SetMaxMemory(size string) bool {
	return cs.memCache.SetMaxMemory(size)
}

//Set 设置⼀个缓存项，并且在expire时间之后过期
func (cs *cacheServer) Set(key string, val interface{}, expire ...time.Duration) bool {
	duration := time.Second * 0
	if len(expire) > 0 {
		duration = expire[0]
	}
	return cs.memCache.Set(key, val, duration)
}

//Get 获取⼀个值
func (cs *cacheServer) Get(key string) (interface{}, bool) {

	return cs.memCache.Get(key)
}

//Del 删除⼀个值
func (cs *cacheServer) Del(key string) bool {
	return cs.memCache.Del(key)
}

//Exists 检测⼀个值 是否存在
func (cs *cacheServer) Exists(key string) bool {
	return cs.memCache.Exists(key)
}

//Flush 情况所有值
func (cs *cacheServer) Flush() bool {
	return cs.memCache.Flush()
}

//Keys 返回所有的key 多少
func (cs *cacheServer) Keys() int64 {
	return cs.memCache.Keys()
}
