package bytes2

import (
	"sync"
	"sync/atomic"

	"github.com/cosiner/gohper/utils/defval"
)

// CanPool is the default checker to check whether a buffer is reuseable, if
// the buffer's capacity is less than half of default bufsize or small
// buffer was not allowed, it will be dropped
var CanPool = func(cap, bufsize int, allowsmall bool) bool {
	return cap >= bufsize || (allowsmall && cap >= bufsize>>1)
}

// DEF_BUFSIZE is default buffer size
const DEF_BUFSIZE = 1024

type Pool interface {
	// Get a buffer, capacity is at least the given size if asLen, length will
	// be set to size, else 0
	Get(size int, asLen bool) []byte
	Put(buf []byte)
	// TryPut will pool the buffer, return whether buffer was pooled by CanPool function
	TryPut(buf []byte) bool
}

// FakePool is a fake pool, it isn't a actual pool
type FakePool struct{}

func (p FakePool) Get(size int, asLen bool) []byte {
	if asLen {
		return make([]byte, size)
	}

	return make([]byte, 0, size)
}

func (p FakePool) Put([]byte) {}

func (p FakePool) TryPut([]byte) bool {
	return false
}

func NewFakePool() Pool {
	return FakePool{}
}

// SyncPool is based on sync.Pool
type SyncPool struct {
	pool       sync.Pool
	bufsize    int
	allowSmall bool
}

func NewSyncPool(bufsize int, allowsmall bool) Pool {
	p := &SyncPool{
		bufsize:    bufsize,
		allowSmall: allowsmall,
	}

	defval.Int(&p.bufsize, DEF_BUFSIZE)
	p.pool.New = func() interface{} {
		return make([]byte, p.bufsize)
	}

	return p
}

func (p *SyncPool) Put(buf []byte) {
	p.TryPut(buf)
}

func (p *SyncPool) TryPut(buf []byte) bool {
	succ := CanPool(cap(buf), p.bufsize, p.allowSmall)
	if succ {
		p.pool.Put(buf)
	}

	return succ
}

func (p *SyncPool) Get(size int, asLen bool) []byte {
	buf := p.pool.Get().([]byte)
	if cap(buf) < size {
		p.pool.Put(buf) // reuse
		buf = make([]byte, size)
	}

	if asLen {
		buf = buf[:size]
	} else {
		buf = buf[:0]
	}

	return buf
}

type poolNode struct {
	buffers []byte
	next    *poolNode
}

// ListPool is based on a linked list, will not expires, shrink by call ShrinkTo
type ListPool struct {
	head  *poolNode
	Count int
	sync.Mutex

	bufsize    int
	allowSmall bool
}

func NewListPool(bufsize int, allowsmall bool) Pool {
	p := &ListPool{
		bufsize:    bufsize,
		allowSmall: allowsmall,
	}

	defval.Int(&p.bufsize, DEF_BUFSIZE)

	return p
}

func (p *ListPool) Get(size int, asLen bool) []byte {
	p.Lock()
	var buf []byte
	if p.head == nil || cap(p.head.buffers) < size {
		if size < p.bufsize { // allocate at least bufsize
			buf = make([]byte, p.bufsize)
		} else {
			buf = make([]byte, size)
		}
	} else {
		buf = p.head.buffers
		p.head = p.head.next
		p.Count--
	}
	p.Unlock()

	if asLen {
		buf = buf[:size]
	} else {
		buf = buf[:0]
	}

	return buf
}

func (p *ListPool) Put(buf []byte) {
	p.TryPut(buf)
}

func (p *ListPool) TryPut(buf []byte) bool {
	succ := CanPool(cap(buf), p.bufsize, p.allowSmall)
	if succ {
		p.Lock()
		p.Count++
		p.head = &poolNode{
			buffers: buf,
			next:    p.head,
		}
		p.Unlock()
	}

	return succ
}

// ShrinkTo cause the ListPool's buffer count reduce to count
func (p *ListPool) ShrinkTo(count int) {
	if count < 0 {
		return
	}

	p.Lock()
	c := p.Count - count
	if c > 0 {
		p.Count = count
	}
	for ; c > 0; c-- {
		p.head = p.head.next
	}
	p.Unlock()
}

// SlotPool is a slice of Pool, it will truncate pools size to power of 2,
// each Get/Put will cycle through the pools
type SlotPool struct {
	pools []Pool
	curr  int32
	mask  int32
}

func NewSlotPool(pools []Pool) Pool {
	count := int32(len(pools)) & (^1)
	if count == 0 {
		panic("at least 2 pools")
	}

	return &SlotPool{
		pools: pools[:count],
		mask:  count - 1,
	}
}

func (p *SlotPool) Get(size int, asLen bool) []byte {
	c := atomic.AddInt32(&p.curr, 1) & p.mask

	return p.pools[c].Get(size, asLen)

}

func (p *SlotPool) Put(buf []byte) {
	p.TryPut(buf)
}

func (p *SlotPool) TryPut(buf []byte) bool {
	c := atomic.LoadInt32(&p.curr) & p.mask

	return p.pools[c].TryPut(buf)
}

// SyncSlotPool create a SlotPool based on SyncPool
func SyncSlotPool(slot, bufsize int, allowsmall bool) Pool {
	var pools = make([]Pool, slot)
	for i := 0; i < slot; i++ {
		pools[i] = NewSyncPool(bufsize, allowsmall)
	}

	return NewSlotPool(pools)
}

// SyncSlotPool create a SlotPool based on ListPool
func ListSlotPool(slot, bufsize int, allowsmall bool) Pool {
	var pools = make([]Pool, slot)
	for i := 0; i < slot; i++ {
		pools[i] = NewListPool(bufsize, allowsmall)
	}

	return NewSlotPool(pools)
}
