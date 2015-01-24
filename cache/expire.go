package cache

// import (
// 	"time"

// 	. "github.com/cosiner/golib/errors"
// 	"github.com/cosiner/golib/types"

// 	"sync"
// )

// type nodeList struct {
// 	node *node
// 	next *nodeList
// }

// type node struct {
// 	time, expire int64
// 	key          string
// 	val          []interface{}
// }

// func (n *node) isExpired() bool {
// 	return n.isExpiredUntil(timeNow())
// }

// func (n *node) isExpiredUntil(t int64) bool {
// 	return n.expire != -1 && n.time+n.expire < t
// }

// type expireCache struct {
// 	lock   *sync.RWMutex
// 	values map[string]*node
// 	head   *nodeList
// 	tail   *nodeList
// }

// func (ec *expireCache) Init(conf string) error {
// 	pair := types.ParsePair(conf, "=")
// 	if pair.HasValue() && pair.Key == "expire" {
// 		if expire, err := pair.IntValue(); err == nil {
// 			ec.expire = expire
// 			return nil
// 		}
// 	}
// 	return Errorf("Wrong format:%s", conf)
// }

// func (ec *expireCache) Get(key string) (val interface{}) {
// 	ec.lock.RLock()
// 	node := ec.values[key]
// 	ec.lock.RUnlock()
// 	if !node.isExpired() {
// 		val = node.val
// 	}
// 	return
// }

// func timeNow() int64 {
// 	return time.Now().UnixNano()
// }

// func (ec *expireCache) Set(key string, val interface{}, expire int64) {
// 	var time int64
// 	if expire == 0 {
// 		return
// 	} else if expire < 0 {
// 		time, expire = 0, -1
// 	} else {
// 		time = timeNow()
// 	}
// 	node := &node{key: key, val: val, time: time, expire: expire}
// 	listNode := &listNode{node: node}
// 	ec.lock.Lock()
// 	ec.values[key] = node
// 	if ec.tail == nil {
// 		ec.tail = listNode
// 	} else {
// 		ec.tail.next = listNode
// 	}
// 	ec.lock.Unlock()
// }

// // Update only update exist key-value pair, if key not exist, return false
// func (ec *expireCache) Update(key string, val interface{}) bool {}

// // Remove key-value pair
// func (ec *expireCache) Remove(key string) {}

// // IsExist check whether item exist
// func (ec *expireCache) IsExist(key string) bool {}

// // Len return current cache count
// func (ec *expireCache) Len() int {}

// // Cap return cache capacity
// func (ec *expireCache) Cap() int {}
