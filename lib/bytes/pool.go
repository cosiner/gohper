package bytes

import (
	"sync"
	"sync/atomic"

	"github.com/cosiner/gohper/lib/defval"
)

// Usable is the default checker to check whether a buffer is reuseable
var Usable = func(cap, bufsize int, allowsmall bool) bool {
	return cap >= bufsize || (allowsmall && cap >= bufsize>>1)
}

const DEF_BUFSIZE = 1024

type Pool interface {
	// Init check buffer size and init sync.Pool, it must be called
	//
	// all bufer with small capacity will be dropped,
	// unless Pool.AllowSmall is true and the size is at least half of bufsize
	Init() Pool
	Get(size int, asLen bool) []byte
	Put(buf []byte)
	CheckPut(buf []byte) bool
}

type FakePool struct{}

func (p FakePool) Init() Pool {
	return p
}

func (p FakePool) Get(size int, asLen bool) []byte {
	bs := make([]byte, 0, size)
	if asLen {
		bs = bs[:size]
	}
	return bs
}

func (p FakePool) Put([]byte) {}

func (p FakePool) CheckPut([]byte) bool {
	return false
}

// SyncPool is a sync.Pool's wrapper with same interface
type SyncPool struct {
	pool       sync.Pool
	Bufsize    int
	AllowSmall bool
}

func NewSyncPool(bufsize int, allowsmall bool) Pool {
	return &SyncPool{
		Bufsize:    bufsize,
		AllowSmall: allowsmall,
	}
}

func (p *SyncPool) Init() Pool {
	defval.Int(&p.Bufsize, DEF_BUFSIZE)
	p.pool.New = func() interface{} {
		return make([]byte, p.Bufsize)
	}
	return p
}

func (p *SyncPool) Put(buf []byte) {
	p.CheckPut(buf)
}

func (p *SyncPool) CheckPut(buf []byte) bool {
	succ := Usable(cap(buf), p.Bufsize, p.AllowSmall)
	if succ {
		p.pool.Put(buf)
	}
	return succ
}

// Get a buffer from pool, if the buffer size is not enough, create a new one,
// put old to pool again. if asLen, the buffer length will be set to size, otherwise 0
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

type ListPool struct {
	head  *poolNode
	Count int
	sync.Mutex

	Bufsize    int
	AllowSmall bool
}

func NewListPool(bufsize int, allowsmall bool) Pool {
	return &ListPool{
		Bufsize:    bufsize,
		AllowSmall: allowsmall,
	}
}

func (p *ListPool) Init() Pool {
	defval.Int(&p.Bufsize, DEF_BUFSIZE)
	return p
}

func (p *ListPool) Get(size int, asLen bool) []byte {
	p.Lock()
	var buf []byte
	if p.head == nil || cap(p.head.buffers) < size {
		if size < p.Bufsize { // allocate at least Bufsize
			buf = make([]byte, p.Bufsize)
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
	p.CheckPut(buf)
}

func (p *ListPool) CheckPut(buf []byte) bool {
	succ := Usable(cap(buf), p.Bufsize, p.AllowSmall)
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

// SlotPool will truncate pools size to power of 2
type SlotPool struct {
	Pools []Pool
	curr  int32
	mask  int32
}

func NewSlotPool(pools []Pool) Pool {
	return &SlotPool{
		Pools: pools,
	}
}

func (p *SlotPool) Init() Pool {
	count := int32(len(p.Pools)) & (^1)
	if count == 0 {
		panic("at least 2 pools")
	}
	p.Pools = p.Pools[:count]
	for i := int32(0); i < count; i++ {
		p.Pools[i].Init()
	}
	p.mask = count - 1
	return p
}

func (p *SlotPool) Get(size int, asLen bool) []byte {
	c := atomic.AddInt32(&p.curr, 1) & p.mask
	return p.Pools[c].Get(size, asLen)

}

func (p *SlotPool) Put(buf []byte) {
	p.CheckPut(buf)
}

func (p *SlotPool) CheckPut(buf []byte) bool {
	c := atomic.LoadInt32(&p.curr) & p.mask
	return p.Pools[c].CheckPut(buf)
}

func SyncSlotPool(slot, bufsize int, allowsmall bool) Pool {
	var pools = make([]Pool, slot)
	for i := 0; i < slot; i++ {
		pools[i] = NewSyncPool(bufsize, allowsmall)
	}
	return NewSlotPool(pools)
}

func ListSlotPool(slot, bufsize int, allowsmall bool) Pool {
	var pools = make([]Pool, slot)
	for i := 0; i < slot; i++ {
		pools[i] = NewListPool(bufsize, allowsmall)
	}
	return NewSlotPool(pools)
}
