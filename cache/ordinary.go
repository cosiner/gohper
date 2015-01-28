package cache

import (
	"math"
	"sync"
)

// capUnlimit means no limit size
const capUnlimit int = math.MaxInt32

// OrdinaryCache is ordinary cache, there is no limit of cache item count
type OrdinaryCache struct {
	values map[string]interface{}
	lock   *sync.RWMutex
}

// Init init ordinary cache, no config need, no error returned
func (oc *OrdinaryCache) Init(string) error {
	oc.values = make(map[string]interface{})
	oc.lock = new(sync.RWMutex)
	return nil
}

// InitVals init ordinary cache with given initial value
func (oc *OrdinaryCache) InitVals(_ string, values map[string]interface{}) error {
	if values != nil {
		oc.values = values
	} else {
		oc.values = make(map[string]interface{})
	}
	oc.lock = new(sync.RWMutex)
	return nil
}

// Len return elements count of ordinary cahce
func (oc *OrdinaryCache) Len() int {
	oc.lock.RLock()
	length := oc.len()
	oc.lock.RUnlock()
	return length
}

func (oc *OrdinaryCache) len() int {
	return len(oc.values)
}

func (oc *OrdinaryCache) Cap() int {
	return capUnlimit
}

func (oc *OrdinaryCache) Get(key string) (val interface{}) {
	oc.lock.RLock()
	val = oc.values[key]
	oc.lock.RUnlock()
	return
}

func (oc *OrdinaryCache) IsExist(key string) bool {
	oc.lock.RLock()
	_, has := oc.values[key]
	oc.lock.RUnlock()
	return has
}

func (oc *OrdinaryCache) Remove(key string) {
	oc.lock.Lock()
	delete(oc.values, key)
	oc.lock.Unlock()
}

func (oc *OrdinaryCache) Set(key string, val interface{}) {
	oc.set(key, val, true)
}

func (oc *OrdinaryCache) Update(key string, val interface{}) bool {
	return oc.set(key, val, false)
}

// set bind value to key
// allowAdd means if or not allow add new value when key don't exist
func (oc *OrdinaryCache) set(key string, val interface{}, allowAdd bool) (ret bool) {
	oc.lock.Lock()
	v := oc.values[key]
	if v != nil || allowAdd {
		oc.values[key] = val
		ret = true
	}
	oc.lock.Unlock()
	return
}

// AccessAllValues access all values exist in cacher
// for no copy and safety access, so need an function parameter
// rather than return values map's reference
func (oc *OrdinaryCache) AccessAllValues(fn func(map[string]interface{})) {
	oc.lock.RLock()
	fn(oc.values)
	oc.lock.RUnlock()
}
