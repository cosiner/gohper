package sync2

import (
	"sync"
	"sync/atomic"

	"github.com/cosiner/gohper/ds/set"
)

type Flag int32

const (
	FLAG_FALSE Flag = 0
	FLAG_TRUE  Flag = 1
)

func (f *Flag) IsTrue() bool {
	return atomic.LoadInt32((*int32)(f)) == int32(FLAG_TRUE)
}

func (f *Flag) MakeTrue() bool {
	return atomic.CompareAndSwapInt32((*int32)(f), int32(FLAG_FALSE), int32(FLAG_TRUE))
}

func (f *Flag) MakeFalse() bool {
	return atomic.CompareAndSwapInt32((*int32)(f), int32(FLAG_TRUE), int32(FLAG_FALSE))
}

type Flags struct {
	flags set.Strings
	mu    sync.Mutex
}

func (f *Flags) IsTrue(s string) bool {
	f.mu.Lock()
	isTrue := f.flags != nil && f.flags.HasKey(s)
	f.mu.Unlock()
	return isTrue
}

func (f *Flags) MakeTrue(s string) (done bool) {
	return f.switchTo(s, true)
}

func (f *Flags) MakeFalse(s string) (done bool) {
	return f.switchTo(s, false)
}

func (f *Flags) switchTo(s string, val bool) (done bool) {
	f.mu.Lock()
	if f.flags == nil {
		f.flags = set.NewStrings()
	}
	if val != f.flags.HasKey(s) {
		done = true
		if val {
			f.flags.Put(s)
		} else {
			f.flags.Remove(s)
		}
	}
	f.mu.Unlock()
	return done
}
