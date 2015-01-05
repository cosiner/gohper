package memcache

import (
	"sync"
)

type ordiCache struct {
	cache map[string]interface{}
	*sync.RWMutex
}

func (oc *ordiCache) init(initSize int, maxsize ...int) error {
	if initSize < 0 {
		initSize = 0
	}
	oc.cache = make(map[string]interface{}, initSize)
	oc.RWMutex = new(sync.RWMutex)
	return nil
}

func (oc *ordiCache) Len() int {
	return len(oc.cache)
}

func (oc *ordiCache) Cap() int {
	return -1
}

func (oc *ordiCache) ChangeCap(offset int) error {
	return nil
}

func (oc *ordiCache) Get(key string) (val interface{}) {
	if len(oc.cache) != 0 {
		oc.RLock()
		val = oc.cache[key]
		oc.RUnlock()
	}
	return
}

func (oc *ordiCache) Remove(key string) {
	if len(oc.cache) != 0 {
		oc.Lock()
		delete(oc.cache, key)
		oc.Unlock()
	}
}

func (oc *ordiCache) Set(key string, val interface{}) {
	oc.set(key, val, true)
}

func (oc *ordiCache) Update(key string, val interface{}) bool {
	return oc.set(key, val, false)
}

func (oc *ordiCache) set(key string, val interface{}, forceSet bool) bool {
	oc.RLock()
	v := oc.cache[key]
	oc.RUnlock()
	if v == nil && !forceSet {
		return false
	}
	oc.Lock()
	oc.cache[key] = val
	oc.Unlock()
	return true
}
