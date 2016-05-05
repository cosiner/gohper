package sync2

import "sync/atomic"

type Barrier struct {
	curr uint32
	n    uint32

	c chan struct{}
}

func NewBarrier(n uint32) *Barrier {
	return &Barrier{
		n: n,
		c: make(chan struct{}),
	}
}

func (b *Barrier) Wait() {
	if atomic.AddUint32(&b.curr, 1) >= b.n {
		close(b.c)
	} else {
		<-b.c
	}
}
