package memcache

import (
	"math"
	"sync"
)

const capUnlimit int = math.MaxInt32

// ordiCache is ordinary cache, there is no limit of cache item count
type ordiCache struct {
	cache map[string]interface{}
	*sync.RWMutex
}

func (oc *ordiCache) Init(maxsize int) {
	oc.cache = make(map[string]interface{}, maxsize)
	oc.RWMutex = new(sync.RWMutex)
}

func (oc *ordiCache) Len() int {
	oc.RLock()
	length := oc.len()
	oc.RUnlock()
	return length
}

func (oc *ordiCache) len() int {
	return len(oc.cache)
}

func (oc *ordiCache) Cap() int {
	return capUnlimit
}

func (oc *ordiCache) Get(key string) (val interface{}) {
	oc.RLock()
	val = oc.cache[key]
	oc.RUnlock()
	return
}

func (oc *ordiCache) Remove(key string) {
	oc.Lock()
	delete(oc.cache, key)
	oc.Unlock()
}

func (oc *ordiCache) Set(key string, val interface{}) {
	oc.set(key, val, true)
}

func (oc *ordiCache) Update(key string, val interface{}) bool {
	return oc.set(key, val, false)
}

// set bind value to key
// allowAdd means if or not allow add new value when key don't exist
func (oc *ordiCache) set(key string, val interface{}, allowAdd bool) (ret bool) {
	oc.Lock()
	v := oc.cache[key]
	if v != nil || allowAdd {
		oc.cache[key] = val
		ret = true
	}
	oc.Unlock()
	return
}
