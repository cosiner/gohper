package sync2

import (
	"math/rand"
	"runtime"
	"sync"
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

	var lock Spinlock
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

func TestAutorefMutex(t *testing.T) {
	t.Log("start testing.")
	mu := NewAutorefMutex(false)

	keys := []string{"1", "2", "3"}
	keyNum := len(keys)
	wg := sync.WaitGroup{}

	routine := 21
	wg.Add(routine)

	for i := 0; i < routine; i++ {
		n := i + 1
		go func() {
			time.Sleep(time.Duration(n%3) * time.Millisecond)
			key := keys[rand.Intn(keyNum)]
			mu.Lock(key)
			t.Logf("routine %d locked %s", n, key)
			time.Sleep(time.Duration(n%3) * time.Millisecond)
			mu.Unlock(key)
			t.Logf("routine %d unlocked %s", n, key)

			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkAutorefMutex(b *testing.B) {
	mu := NewAutorefMutex(false)

	for i := 0; i < b.N; i++ {
		mu.Lock("a")
		mu.Unlock("a")
	}
}
