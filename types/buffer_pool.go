package types

import (
	"bytes"
	"sync"
)

type (
	// BufferReleaser is bytes.Buffer + Release
	BufferReleaser struct {
		bp  *BufferPool
		buf []byte
		*bytes.Buffer
	}
	// poolNode is node type of BufferPool
	poolNode struct {
		data []byte
		next *poolNode
	}
	// BufferPool represent a buffer pool
	BufferPool struct {
		buffers *poolNode
		lock    *sync.Mutex
	}
)

// Release release buffer to BufferPool
func (bc *BufferReleaser) Release() {
	bc.bp.Collect(bc.buf)
}

// NewBufferPool create a new BufferPool
func NewBufferPool() *BufferPool {
	return &BufferPool{
		buffers: new(poolNode),
		lock:    new(sync.Mutex),
	}
}

// Buffer() acquire a new BufferRelaser from BufferPool
func (bp *BufferPool) Buffer() *BufferReleaser {
	var data []byte
	bp.lock.Lock()
	if buffers := bp.buffers; buffers == nil {
		data = make([]byte, 2048)
	} else {
		data = buffers.data
		bp.buffers = buffers.next
	}
	bp.lock.Unlock()
	return &BufferReleaser{
		bp:     bp,
		buf:    data,
		Buffer: bytes.NewBuffer(data),
	}
}

// Collect collect a buffer to BufferPool for later acquire
func (bp *BufferPool) Collect(buf []byte) {
	bn := &poolNode{
		data: buf[0:0],
		next: nil,
	}
	bp.lock.Lock()
	bn.next = bp.buffers
	bp.buffers = bn
	bp.lock.Unlock()
}
