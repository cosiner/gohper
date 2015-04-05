package signal

import (
	"sync"
)

// WARNING: UNTESTED
type Signal struct {
	s    chan byte
	lock sync.Mutex
}

func New() *Signal {
	return &Signal{
		s: make(chan byte, 1),
	}
}

func (s *Signal) Wait() {
	<-s.s
}

func (s *Signal) Notify() {
	s.lock.Lock()
	if len(s.s) == 0 {
		s.s <- 0
	}
	s.lock.Unlock()
}
