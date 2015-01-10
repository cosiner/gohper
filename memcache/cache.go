// Package memcache Implement a simple memory cache container for go,
// it has three type container, ordinarry cache, don't limit capaticy,
// random eliminate cache, based on ordinary cache, and lru eliminate cache
// random and lru cache can dynamic change  the cache capacity, ordinary cache's
// ChangeCap function make no difference
package memcache

import (
	. "github.com/cosiner/golib/errors"
)

// CacheType is implemented cache algorithm
type CacheType int8

func (ct CacheType) String() (str string) {
	switch ct {
	case ORDINARY:
		str = "Ordinary"
	case RANDOM:
		str = "Random-eliminate"
	case LRU:
		str = "LRU-eliminate"
	}
	return
}

const (
	// ORDINARY is normal cache, hs no elimination
	ORDINARY CacheType = 1 << iota
	// RANDOM is a random eliminate algorithm
	RANDOM
	// LRU is lru eliminate algorithm
	LRU
)

// MemCache is cache interface
// all method a safe for concurrent
type MemCache interface {
	Init(maxSize int)
	// Get by key
	Get(key string) interface{}
	// Set key-value pair, if no remaining space, trigger a elimination
	// if key already exist, will be replaced
	Set(key string, val interface{})
	// Update only update exist key-value pair, if key not exist, return false
	Update(key string, val interface{}) bool
	// Remove key-value pair
	Remove(key string)
	// Len return current cache count
	Len() int
	// Cap return cache capacity
	Cap() int
}

// Cacher return actual cache container
// for cacher with elimination:Random and LRU, maxsize is the max capcity of cache
// for ordinary cache, it only used to initial cache space
func Cacher(typ CacheType, maxSize int) (cache MemCache) {
	Assert(maxSize > 0, Errorf("Invalid  max size"))
	switch typ {
	case ORDINARY:
		cache = new(ordiCache)
	case RANDOM:
		cache = new(randCache)
	case LRU:
		cache = new(lruCache)
	default:
		panic("Unsupported Cache Type")
	}
	cache.Init(maxSize)
	return cache
}
