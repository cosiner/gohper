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

const (
	// ORDINARY is normal cache, hs no elimination
	ORDINARY CacheType = 1 << iota
	// RANDOM is a random eliminate algorithm
	RANDOM
	// LRU is lru eliminate algorithm
	LRU
)

// MemCache is cache interface
type MemCache interface {
	init(initSize int, maxSize ...int) error
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
	// ChangeCap will change cache capacity, if offset > 0, increate,
	// else, delete cache data by offse, for ordinary and random cache, random delete,
	// for lru cache, delete lru cache.
	// offset is not allowed that -offset >= Cap() when offset < 0,
	// in other words, cache's capacity must larger than 0
	ChangeCap(offset int) error
}

// SizeNegativeError
SizeNegativeError := Err("Cache Size must > 0")

// Cacher return actual cache container, if initSize is larger than maxSize,
// initSize will setup tp maxSize, for ordinary cache, there is no maxSize
func Cacher(typ CacheType, initSize int, maxSize ...int) MemCache {
	var cache MemCache
	switch typ {
	case ORDINARY:
		cache = new(ordiCache)
	case RANDOM:
		cache = new(randCache)
	case LRU:
		cache = new(lruCache)
	default:
		return cache
	}
	cache.init(initSize, maxSize...)
	return cache
}
