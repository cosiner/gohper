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

// CacherType is implemented cache algorithm
type CacherType int8

func (ct CacherType) String() (str string) {
	switch ct {
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
	// RANDOM is a random eliminate algorithm
	RANDOM = iota
	// LRU is lru eliminate algorithm
	LRU
	// REDIS is redis cacher
	REDIS

	ErrUnsupportedType = Err("Not supported cache type")
	ErrWrongFormat     = Err("Wrong format of config")
)

// Cache is cache interface
// all method a safe for concurrent
type Cache interface {
	Init(config string) error
	InitVals(config string, values map[string]interface{}) error
	// Get by key
	Get(key string) interface{}
	// if key already exist, will be replaced
	Set(key string, val interface{})
	// SetValues replace all values
	// Update only update exist key-value pair, if key not exist, return false
	Update(key string, val interface{}) bool
	// Remove key-value pair
	Remove(key string)
	// IsExist check whether item exist
	IsExist(key string) bool
	// Size return current cache count
	Size() int
	// Cap return cache capacity
	Cap() int
}

// NewCache return a actual cache container
// for cacher with elimination:Random and LRU, maxsize is the max capcity of cache
// for ordinary cache, no config need, no error returned
func NewCache(typ CacherType, config string) (cache Cache, err error) {
	switch typ {
	case RANDOM:
		cache = new(randCache)
	case LRU:
		cache = new(lruCache)
	case REDIS:
		cache = new(RedisCache)
	default:
		return nil, ErrUnsupportedType
	}
	return cache, cache.Init(config)
}

// fixSize fix values's size by random remove the rest elemtents
func fixSize(values map[string]interface{}, size int) {
	for k := range values {
		if len(values) > size {
			delete(values, k)
		}
	}
}

// parseMaxSize parse maxsize  from config string
func parseMaxSize(config string) (maxsize int, err error) {
	pair := types.ParsePair(config, "=")
	if pair.NoKey() || pair.NoValue() || pair.Key != "maxsize" {
		err = ErrWrongFormat
	} else if maxsize, err = pair.IntValue(); err != nil || maxsize <= 0 {
		err = ErrWrongFormat
	}
	return
}
