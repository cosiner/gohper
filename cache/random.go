package cache

type randCache struct {
	maxSize int
	OrdinaryCache
}

func (rc *randCache) Init(config string) (err error) {
	var maxsize int
	if maxsize, err = parseMaxSize(config); err == nil {
		rc.OrdinaryCache.Init("")
		rc.maxSize = maxsize
	}
	return
}

func (rc *randCache) InitVals(config string, values map[string]interface{}) (err error) {
	var maxsize int
	if maxsize, err = parseMaxSize(config); err == nil {
		fixSize(values, maxsize)
		rc.OrdinaryCache.InitVals("", values)
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
	rc.lock.Lock()
	values := rc.values
	v := values[key]
	if v != nil || forceSet {
		if rc.len() == rc.cap() {
			for k, _ := range values {
				delete(values, k)
				break
			}
		}
		values[key] = val
		ret = true
	}
	rc.lock.Unlock()
	return
}
