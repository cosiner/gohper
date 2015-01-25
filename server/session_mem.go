package server

import (
	"sync"
	"time"

	"github.com/cosiner/gomodule/config"
)

func unixNow() uint64 {
	return uint64(time.Now().Unix())
}

type memStoreNode struct {
	time, expire uint64
	value        map[string]interface{}
}

func newMemStoreNode(val map[string]interface{}, expire uint64) *memStoreNode {
	return &memStoreNode{
		time:   unixNow(),
		expire: expire,
		value:  val,
	}
}

func (msn *memStoreNode) isExpired() bool {
	return msn.isExpiredTill(unixNow())
}

func (msn *memStoreNode) isExpiredTill(t uint64) (expired bool) {
	if msn.time+msn.expire <= t {
		expired = true
	}
	return
}

type MemStore struct {
	values map[string]*memStoreNode
	lock   *sync.RWMutex
}

func (ms *MemStore) Init(conf string) (err error) {
	c := config.NewConfig(config.LINE)
	if err = c.ParseString(conf); err == nil {
		ms.values = make(map[string]*memStoreNode)
		ms.lock = new(sync.RWMutex)
		go ms.gc(c.IntValDef("gcinterval", 600))
	}
	return
}

func (ms *MemStore) IsExist(id string) bool {
	ms.lock.RLock()
	_, has := ms.values[id]
	ms.lock.RUnlock()
	return has
}

func (ms *MemStore) Get(id string) (val map[string]interface{}) {
	ms.lock.RLock()
	msn := ms.values[id]
	ms.lock.RUnlock()
	if msn != nil && !msn.isExpired() {
		val = msn.value
	}
	return
}

func (ms *MemStore) Set(id string, val map[string]interface{}, expire uint64) {
	ms.lock.Lock()
	ms.values[id] = newMemStoreNode(val, expire)
	ms.lock.Unlock()
}

func (ms *MemStore) Rename(oldId, newId string) {
	values := ms.values
	ms.lock.RLock()
	msn := values[oldId]
	ms.lock.RUnlock()
	if msn != nil {
		ms.lock.Lock()
		delete(values, oldId)
		values[newId] = msn
		ms.lock.Unlock()
	}
}

func (ms *MemStore) gc(inteval int) {
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
		}
	}
}
