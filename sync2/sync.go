package sync2

import "sync/atomic"

// Once is an object that will perform exactly one action unless call Undo.
type Once uint32

// Do will do f only once no matter it's successful or not
// if f is blocked, Do will also be
func (o *Once) Do(fn ...func()) bool {
	if !atomic.CompareAndSwapUint32((*uint32)(o), 0, 1) {
		return false
	}
	for _, f := range fn {
		if f != nil {
			f()
		}
	}
	return true
}
