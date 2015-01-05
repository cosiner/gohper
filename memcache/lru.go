package memcache

import (
	"container/list"
	"sync"
)

// lruCacheEntry is a item of lru cache
type lruCacheEntry struct {
	accessCount int // access accessCount
	key         string
	val         interface{}
}

// init init lruCacheEntry
func (ci *lruCacheEntry) init(key string, val interface{}, accessCount int) {
	ci.accessCount = accessCount
	ci.key = key
	ci.val = val
}

// lruCache is a lru cacher
type lruCache struct {
	cacheData  *list.List
	cacheIndex map[string]*list.Element
	maxSize    int
	*sync.RWMutex
}

// init init lru cache
func (lc *lruCache) init(initSize int, maxSize ...int) error {
	if len(maxSize) <= 0 {
		return SizeNegativeError
	}
	if initSize < 0 {
		initSize = 0
	}
	if initSize > maxSize[0] {
		initSize = maxSize[0]
	}

	lc.cacheData = list.New()
	lc.cacheIndex = make(map[string]*list.Element, initSize)
	lc.maxSize = maxSize[0]
	lc.RWMutex = new(sync.RWMutex)
	return nil
}

// Len return current cache count
func (lc *lruCache) Len() int {
	return len(lc.cacheIndex)
}

// Cap return cache capacity
func (lc *lruCache) Cap() int {
	return lc.maxSize
}
func (lc *lruCache) ChangeCap(offset int) error {

	return nil
}
func (lc *lruCache) incrGet(elem *list.Element) interface{} {
	entry := elem.Value.(*lruCacheEntry)
	entry.accessCount++
	return entry.val
}

func (lc *lruCache) Get(key string) (val interface{}) {
	if len(lc.cacheIndex) != 0 {
		lc.RLock()
		elem, has := lc.cacheIndex[key]
		lc.RUnlock()
		if has {
			lc.Lock()
			lc.cacheData.MoveToFront(elem)
			val = lc.incrGet(elem)
			lc.Unlock()
		}
	}
	return
}

func (lc *lruCache) Remove(key string) {
	if len(lc.cacheIndex) != 0 {
		lc.RLock()
		elem, has := lc.cacheIndex[key]
		lc.RUnlock()
		if has {
			lc.Lock()
			lc.cacheData.Remove(elem)
			delete(lc.cacheIndex, key)
			lc.Unlock()
		}
	}
}

func (lc *lruCache) Set(key string, val interface{}) {
	lc.set(key, val, true)
}

func (lc *lruCache) Update(key string, val interface{}) bool {
	return lc.set(key, val, false)
}

func (lc *lruCache) set(key string, val interface{}, forceSet bool) bool {
	var entry *lruCacheEntry
	lc.RLock()
	elem := lc.cacheIndex[key]
	size := len(lc.cacheIndex)
	lc.RUnlock()
	lc.Lock() // lock big area only for convenience
	if elem == nil {
		if !forceSet {
			return false // don't exist and don't allow set
		} else if size == lc.maxSize {
			elem = lc.cacheData.Back() // remove last and reuse for new value
			entry = elem.Value.(*lruCacheEntry)
			lc.cacheData.Remove(elem)
			delete(lc.cacheIndex, entry.key)
		} else {
			elem = new(list.Element) // has remaining size
			entry = new(lruCacheEntry)
			elem.Value = entry
		}
	} else {
		lc.cacheData.Remove(elem) // if exist, remove it
	}
	entry.init(key, val, 1)      // setup value
	lc.cacheData.PushFront(elem) // insert to front
	lc.cacheIndex[key] = elem
	lc.Unlock()
	return true
}
