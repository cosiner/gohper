package server

import (
	"sync"
	"time"

	"github.com/cosiner/gomodule/config"
)

type (
	// memStoreNode represent a store node of memStore
	// contains created time and it's lifetime
	memStoreNode struct {
		time     uint64
		lifetime int64 // if lifetime < 0, means never expired
		value    Values
	}

	// memStore is a store in memory with lifetime manage
	memStore struct {
		values      map[string]*memStoreNode
		rmChan      chan string
		destroyChan chan bool
		lock        *sync.RWMutex
	}
)

// unixNow return now time as unix time seconds
func unixNow() uint64 {
	return uint64(time.Now().Unix())
}

// newMemStoreNode create a new store node for memStore
func newMemStoreNode(values Values, lifetime int64) (msn *memStoreNode) {
	if lifetime != 0 {
		msn = &memStoreNode{
			time:     unixNow(),
			lifetime: lifetime,
			value:    values,
		}
	}
	return
}

// isExpired check whether store node is expired now
func (msn *memStoreNode) isExpired() bool {
	return msn.isExpiredTill(unixNow())
}

// isExpiredTill check whether store node is expired till given time
func (msn *memStoreNode) isExpiredTill(time uint64) (expired bool) {
	if msn.lifetime < 0 {
		expired = false
	} else if msn.time+uint64(msn.lifetime) <= time {
		expired = true
	}
	return
}

// NewMemStore create a session store in memory
func NewMemStore() SessionStore {
	return new(memStore)
}

// Init init memStore,  config is like "gcinterval=*&rmbacklog=*"
func (ms *memStore) Init(conf string) (err error) {
	c := config.NewConfig(config.LINE)
	if err = c.ParseString(conf); err == nil {
		ms.values = make(map[string]*memStoreNode)
		ms.rmChan = make(chan string, c.IntValDef(SESSION_MEM_RMBACKLOG, DEF_SESSION_MEM_RMBACKLOG))
		ms.destroyChan = make(chan bool, 1)
		ms.lock = new(sync.RWMutex)
		go ms.gc(c.IntValDef(SESSION_MEM_GCINTERVAL, DEF_SESSION_MEM_GCINTERVAL))
	}
	return
}

// Destroy destroy memory store, release resources
func (ms *memStore) Destroy() {
	ms.destroyChan <- true
	ms.lock.Lock()
	ms.values = nil
	ms.lock.Unlock()
	<-ms.destroyChan
	close(ms.rmChan)
	close(ms.destroyChan)
	ms.lock = nil
}

// IsExist check whether given id of node is exist
func (ms *memStore) IsExist(id string) bool {
	return ms.Get(id) != nil
}

// Get return values bind with given id, if id already expired, then remove it
func (ms *memStore) Get(id string) (values Values) {
	ms.lock.RLock()
	msn := ms.values[id]
	ms.lock.RUnlock()
	if msn != nil && !ms.expiredRemove(msn, id) {
		values = msn.value
	}
	return
}

// Save save values with given id and lifetime time
func (ms *memStore) Save(id string, values Values, lifetime int64) {
	if msn := newMemStoreNode(values, lifetime); msn != nil {
		ms.lock.Lock()
		ms.values[id] = msn
		ms.lock.Unlock()
	}
}

// Rename perform rename operation that move all values of old id to new id
// and delete old id
func (ms *memStore) Rename(oldId, newId string) {
	values := ms.values
	ms.lock.RLock()
	msn := values[oldId]
	ms.lock.RUnlock()
	if msn != nil && !ms.expiredRemove(msn, oldId) {
		ms.lock.Lock()
		delete(values, oldId)
		values[newId] = msn
		ms.lock.Unlock()
	}
}

// expiredRemove check whether store node is expired, if true, remove it
func (ms *memStore) expiredRemove(msn *memStoreNode, id string) (expired bool) {
	if expired = msn.isExpired(); expired {
		ms.rmChan <- id
	}
	return
}

// gc perform expired store node collection for memStore
func (ms *memStore) gc(inteval int) {
	if inteval <= 0 {
		return
	}
	values, ticker := ms.values, time.NewTicker(time.Duration(inteval)*time.Second)
	for {
		select {
		case <-ticker.C:
			now := unixNow()
			ms.lock.Lock()
			for id, v := range values {
				if v.isExpiredTill(now) {
					delete(values, id)
				}
			}
			ms.lock.Unlock()
		case id := <-ms.rmChan:
			ms.lock.Lock()
			delete(values, id)
			ms.lock.Unlock()
		case <-ms.destroyChan:
			ms.destroyChan <- true
			return
		}
	}
}
