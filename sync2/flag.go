package sync2

import "sync/atomic"

type Flag int32

func (f *Flag) IsTrue() bool {
	return atomic.LoadInt32((*int32)(f)) == 1
}

func (f *Flag) IsFalse() bool {
	return atomic.LoadInt32((*int32)(f)) == 0
}

func (f *Flag) MakeTrue() bool {
	return atomic.CompareAndSwapInt32((*int32)(f), 0, 1)
}

func (f *Flag) MakeFalse() bool {
	return atomic.CompareAndSwapInt32((*int32)(f), 1, 0)
}
