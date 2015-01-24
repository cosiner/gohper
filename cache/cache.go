// Package cache Implement a simple memory cache container for go,
// it has three type container, ordinarry cache, don't limit capaticy,
// random eliminate cache, based on ordinary cache, and lru eliminate cache
// random and lru cache can dynamic change  the cache capacity, ordinary cache's
// ChangeCap function make no difference
package cache

import (
	. "github.com/cosiner/golib/errors"
	"github.com/cosiner/golib/types"
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
	case REDIS:
		str = "Redis"
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
	// REDIS is redis cacher
	REDIS
)

// Cache is cache interface
// all method a safe for concurrent
type Cache interface {
	Init(config string) error
	// Get by key
	Get(key string) interface{}
	// Set key-value pair, if no remaining space, trigger a elimination
	// if key already exist, will be replaced
	Set(key string, val interface{})
	// Update only update exist key-value pair, if key not exist, return false
	Update(key string, val interface{}) bool
	// Remove key-value pair
	Remove(key string)
	// IsExist check whether item exist
	IsExist(key string) bool
	// Len return current cache count
	Len() int
	// Cap return cache capacity
	Cap() int
}

// NewCache return a actual cache container
// for cacher with elimination:Random and LRU, maxsize is the max capcity of cache
// for ordinary cache, it only used to initial cache space
func NewCache(typ CacheType, config string) (cache Cache, err error) {
	switch typ {
	case ORDINARY:
		cache = new(ordiCache)
	case RANDOM:
		cache = new(randCache)
	case LRU:
		cache = new(lruCache)
	case REDIS:
		cache = new(RedisCache)
	default:
		return nil, Err("Not supported cache type")
	}
	return cache, cache.Init(config)
}

func parseMaxSize(config string) (maxsize int, err error) {
	pair := types.ParsePair(config, "=")
	if pair.NoKey() || pair.NoValue() || pair.Key != "maxsize" {
		err = Err("Wrong format of config")
	} else {
		if maxsize, err = pair.IntValue(); err != nil || maxsize <= 0 {
			err = Err("Wrong format of config")
		}
	}
	return
}
