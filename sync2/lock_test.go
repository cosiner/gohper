package sync2

import (
	"runtime"
	"testing"
	"time"

	"github.com/cosiner/gohper/testing2"
)

func TestSpinLock(t *testing.T) {
	tt := testing2.Wrap(t)
	if runtime.NumCPU() == 1 {
		return
	}

	runtime.GOMAXPROCS(2)

	var lock SpinLock
	go func() {
		lock.Lock()
		lock.Unlock()
	}()
	lock.Lock()
	time.Sleep(1 * time.Millisecond)
	lock.Unlock()

	defer tt.Recover()
	lock.Unlock()
}
