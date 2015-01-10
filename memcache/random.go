package memcache

type randCache struct {
	maxSize int
	ordiCache
}

func (rc *randCache) Init(maxSize int) {
	rc.maxSize = maxSize
	rc.ordiCache.Init(maxSize)
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
	v := rc.cache[key]
	if v != nil || forceSet {
		if rc.len() == rc.cap() {
			for k, _ := range rc.cache {
				delete(rc.cache, k)
				break
			}
		}
		rc.cache[key] = val
		ret = true
	}
	rc.Unlock()
	return
}
