package sync2

import (
	"sync"
)

// LockCond is a wrapper of sync.Cond, each method will be called with lock
type LockCond sync.Cond

func WrapCond(cond *sync.Cond) *LockCond {
	return (*LockCond)(cond)
}

func NewLockCond(l sync.Locker) *LockCond {
	if l == nil {
		l = new(sync.Mutex)
	}

	return (*LockCond)(sync.NewCond(l))
}

func (c *LockCond) Signal() {
	cond := (*sync.Cond)(c)

	cond.L.Lock()
	cond.Signal()
	cond.L.Unlock()
}

func (c *LockCond) Wait() {
	cond := (*sync.Cond)(c)

	cond.L.Lock()
	cond.Wait()
	cond.L.Unlock()
}

func (c *LockCond) Broadcast() {
	cond := (*sync.Cond)(c)

	cond.L.Lock()
	cond.Broadcast()
	cond.L.Unlock()
}

func (c *LockCond) Cond() *sync.Cond {
	return (*sync.Cond)(c)
}
