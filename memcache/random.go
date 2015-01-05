package memcache

import (
	"fmt"
)

type randCache struct {
	maxSize int
	ordiCache
}

func (rc *randCache) init(initSize int, maxSize ...int) error {
	if len(maxSize) <= 0 {
		return SizeNegativeError
	}
	rc.maxSize = maxSize[0]
	return rc.ordiCache.init(initSize, 0)
}

func (rc *randCache) Cap() int {
	return rc.maxSize
}

func (rc *randCache) ChangeCap(offset int) error {
	if offset == 0 {
		return nil
	} else if offset > 0 {
		rc.Lock()
		rc.maxSize += offset
		rc.Unlock()
	} else {
		// if -offsett >= r
	}
	return nil
}

func (rc *randCache) Set(key string, val interface{}) {
	rc.set(key, val, true)
}

func (rc *randCache) Update(key string, val interface{}) bool {
	return rc.set(key, val, false)
}

func (rc *randCache) set(key string, val interface{}, forceSet bool) bool {
	rc.RLock()
	v := rc.cache[key]
	size := len(rc.cache)
	rc.RUnlock()
	if v == nil && !forceSet {
		return false
	}
	rc.Lock()
	if size == rc.maxSize {
		for k, _ := range rc.cache {
			delete(rc.cache, k)
			break
		}
	}
	rc.cache[key] = val
	rc.Unlock()
	return true
}
