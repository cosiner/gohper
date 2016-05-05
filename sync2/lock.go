package sync2

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Spinlock int32

func (s *Spinlock) Lock() {
	for {
		for atomic.LoadInt32((*int32)(s)) != 0 {
		}
		if atomic.CompareAndSwapInt32((*int32)(s), 0, 1) {
			return
		}
	}
}

func (s *Spinlock) Unlock() {
	if !atomic.CompareAndSwapInt32((*int32)(s), 1, 0) {
		panic("unlock unlocked spinlock")
	}
}

type rwLocker interface {
	sync.Locker
	RLock()
	RUnlock()
	Ref() int
	SetRef(int)
}

type mutexRef int

func (m *mutexRef) Ref() int {
	return int(*m)
}

func (m *mutexRef) SetRef(ref int) {
	*m = mutexRef(ref)
}

type refRWMutex struct {
	mutexRef
	sync.RWMutex
}

type refMutex struct {
	mutexRef
	sync.Mutex
}

func (r *refMutex) RLock() {
	r.Mutex.Lock()
}

func (r *refMutex) RUnlock() {
	r.Mutex.Unlock()
}

type AutorefMutex struct {
	mu   sync.Mutex
	pool sync.Pool

	mus map[string]rwLocker
}

func NewAutorefMutex(rw bool) *AutorefMutex {
	var new func() interface{}
	if rw {
		new = func() interface{} {
			return &refRWMutex{}
		}
	} else {
		new = func() interface{} {
			return &refMutex{}
		}
	}
	return &AutorefMutex{
		pool: sync.Pool{
			New: new,
		},
		mus: make(map[string]rwLocker),
	}
}

func (m *AutorefMutex) newLocker(key string) rwLocker {
	mu := m.pool.Get().(rwLocker)
	mu.SetRef(0)
	m.mus[key] = mu
	return mu
}

func (m *AutorefMutex) freeLocker(mu rwLocker) {
	m.pool.Put(mu)
}

func (m *AutorefMutex) locker(key string, unlock bool) (recyle bool, mu rwLocker) {
	m.mu.Lock()
	mu, has := m.mus[key]
	if unlock {
		if !has {
			m.mu.Unlock()
			panic("unlock unexisted key")
		}
		ref := mu.Ref()

		if ref < 1 {
			m.mu.Unlock()
			panic(fmt.Sprintf("impossible reference count: %d for key: %s", ref, key))
		}

		mu.SetRef(ref - 1)
		recyle = ref == 1
		if recyle {
			delete(m.mus, key)
		}
	} else {
		if !has {
			mu = m.newLocker(key)
		}
		mu.SetRef(mu.Ref() + 1)
	}
	m.mu.Unlock()

	return recyle, mu
}

func (m *AutorefMutex) Lock(key string) {
	_, mu := m.locker(key, false)
	mu.Lock()
}

func (m *AutorefMutex) Unlock(key string) {
	// unlock mean works done, it's safe to create another mutex with same key
	// even if current mutex hasn't been unlocked really.
	recyle, mu := m.locker(key, true)
	mu.Unlock()
	if recyle {
		m.freeLocker(mu)
	}
}

func (m *AutorefMutex) RLock(key string) {
	_, mu := m.locker(key, false)
	mu.RLock()
}

func (m *AutorefMutex) RUnlock(key string) {
	recyle, mu := m.locker(key, true)
	mu.RUnlock()
	if recyle {
		m.freeLocker(mu)
	}
}
