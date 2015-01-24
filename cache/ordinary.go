package cache

import (
	"math"
	"sync"
)

const capUnlimit int = math.MaxInt32

// ordiCache is ordinary cache, there is no limit of cache item count
type ordiCache struct {
	cache map[string]interface{}
	lock  *sync.RWMutex
}

func (oc *ordiCache) Init(config string) (err error) {
	var maxsize int
	if maxsize, err = parseMaxSize(config); err == nil {
		oc.cache = make(map[string]interface{}, maxsize)
		oc.lock = new(sync.RWMutex)
	}
	return
}

func (oc *ordiCache) Len() int {
	oc.lock.RLock()
	length := oc.len()
	oc.lock.RUnlock()
	return length
}

func (oc *ordiCache) len() int {
	return len(oc.cache)
}

func (oc *ordiCache) Cap() int {
	return capUnlimit
}

func (oc *ordiCache) Get(key string) (val interface{}) {
	oc.lock.RLock()
	val = oc.cache[key]
	oc.lock.RUnlock()
	return
}

func (oc *ordiCache) IsExist(key string) bool {
	oc.lock.RLock()
	_, has := oc.cache[key]
	oc.lock.RUnlock()
	return has
}

func (oc *ordiCache) Remove(key string) {
	oc.lock.Lock()
	delete(oc.cache, key)
	oc.lock.Unlock()
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
	oc.lock.Lock()
	v := oc.cache[key]
	if v != nil || allowAdd {
		oc.cache[key] = val
		ret = true
	}
	oc.lock.Unlock()
	return
}
