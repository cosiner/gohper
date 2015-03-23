package cache

import (
	"github.com/cosiner/gohper/lib/types"
)

type randCache struct {
	maxSize int
	*types.LockedValues
}

func (rc *randCache) Init(config string) (err error) {
	var maxsize int
	if maxsize, err = parseMaxSize(config); err == nil {
		rc.LockedValues = types.NewLockedValues()
		rc.maxSize = maxsize
	}
	return
}

func (rc *randCache) InitVals(config string, values map[string]interface{}) (err error) {
	var maxsize int
	if maxsize, err = parseMaxSize(config); err == nil {
		fixSize(values, maxsize)
		rc.LockedValues = types.NewLockedValuesWith(values)
		rc.maxSize = maxsize
	}
	return
}

func (rc *randCache) Cap() int {
	// rc.RLock() // currently don't need lock for maxSize will not be modified
	c := rc.cap()
	// rc.RUnlock()
	return c
}

func (rc *randCache) cap() int {
	return rc.maxSize
}

func (rc *randCache) Set(key string, val interface{}) {
	rc.set(key, val, true)
}

func (rc *randCache) Update(key string, val interface{}) bool {
	return rc.set(key, val, false)
}

func (rc *randCache) set(key string, val interface{}, forceSet bool) (ret bool) {
	rc.Lock()
	values := rc.Values
	if values.IsExist(key) || forceSet {
		if values.Size() == rc.cap() {
			values.RandomRemove()
		}
		values.Set(key, val)
		ret = true
	}
	rc.Unlock()
	return
}
