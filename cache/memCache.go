package cache

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type menCache struct {
	// 最大内存
	maxMemorySize int64
	// 最大内存字符串表示
	maxMemorySizeStr string
	// 当前已使用内存
	currMemorySize int64
	//values m 缓存键值对
	values map[string]*memCacheValue
	// 读写锁
	locker sync.RWMutex
	// 清除 过期缓存时间间隔
	clearExpiredItemTimeInterval time.Duration
}

type memCacheValue struct {
	//value 值
	val interface{}
	// 过期时间
	expireTime time.Time
	// 有效时长
	expire time.Duration
	// value 大小
	size int64
}

func NewMenCache() *menCache {
	mc := &menCache{
		values:                       make(map[string]*memCacheValue, 0),
		clearExpiredItemTimeInterval: time.Second * 10,
	}
	go mc.clearExpiredTime()
	return mc
}

// SetMaxMemory size : 1KB 100KB 1MB 2MB 1GB
func (mc *menCache) SetMaxMemory(size string) bool {
	mc.maxMemorySize, mc.maxMemorySizeStr = ParseSize(size)
	fmt.Println(mc.maxMemorySize, mc.maxMemorySizeStr)
	//fmt.Println(" called set max memory")
	return false
}

// Set  将value 写入 缓存
func (mc *menCache) Set(key string, val interface{}, expire time.Duration) bool {
	mc.locker.Lock()
	defer mc.locker.Unlock()
	//fmt.Println("called set")

	v := &memCacheValue{
		val:        val,
		expireTime: time.Now().Add(expire),
		size:       GetValSize(val),
		expire:     expire,
	}
	mc.del(key)
	mc.add(key, v)

	if mc.currMemorySize > mc.maxMemorySize {
		mc.del(key)
		log.Println(fmt.Sprintf("max memory size %s", mc.maxMemorySize))
		panic(fmt.Sprintf("max memory size %s", mc.maxMemorySize))
	}
	return true
}

func (mc menCache) get(key string) (*memCacheValue, bool) {

	val, ok := mc.values[key]
	return val, ok
}

func (mc *menCache) del(key string) {
	tmp, ok := mc.get(key)
	if ok && tmp != nil {
		mc.currMemorySize -= tmp.size
		delete(mc.values, key)
	}
}

func (mc *menCache) add(key string, val *memCacheValue) {

	mc.values[key] = val
	mc.currMemorySize += val.size
}

// Get   根据Key值获取value
func (mc *menCache) Get(key string) (interface{}, bool) {
	mc.locker.RLock()
	defer mc.locker.RUnlock()
	//fmt.Println("called get")
	mcv, ok := mc.get(key)
	if ok {
		//判断缓存舒服过期
		if mcv.expire != 0 && mcv.expireTime.Before(time.Now()) {
			mc.del(key)
			return nil, false
		}
		return mcv.val, ok
	}
	return nil, false
}

// Del   删除Key值
func (mc *menCache) Del(key string) bool {
	mc.locker.Lock()
	defer mc.locker.Unlock()
	//fmt.Println("called del")
	mc.del(key)
	return true
}

// Exists   判断key是否存在
func (mc *menCache) Exists(key string) bool {
	mc.locker.RLock()
	defer mc.locker.RUnlock()
	//fmt.Println("called exists")

	_, ok := mc.values[key]
	return ok
}

// Flush   清空所有key
func (mc *menCache) Flush() bool {
	mc.locker.Lock()
	defer mc.locker.Unlock()

	mc.values = make(map[string]*memCacheValue, 9)
	mc.currMemorySize = 0
	return true
}

// Keys   获取缓存中的所有key数量
func (mc *menCache) Keys() int64 {
	mc.locker.RLock()
	defer mc.locker.RUnlock()
	return int64(len(mc.values))
}

func (mc *menCache) clearExpiredTime() {
	ticker := time.NewTicker(mc.clearExpiredItemTimeInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for key, item := range mc.values {
				if item.expire != 0 && time.Now().After(item.expireTime) {
					mc.locker.Lock()
					mc.del(key)
					mc.locker.Unlock()
				}
			}
		}
	}
}
