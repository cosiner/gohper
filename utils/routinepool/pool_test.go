package routinepool

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/cosiner/gohper/testing2"
)

func TestRoutinePool(t *testing.T) {
	tt := testing2.Wrap(t)

	var done int64
	var jobs int64
	pool := New(func(job Job) {
		time.Sleep(1 * time.Millisecond)
		atomic.AddInt64(&done, 1)

	}, 20, 20, 0)
	for i := 0; i < 10; i++ {
		go func(i int) {
			for j := 0; j < 1000; j++ {
				atomic.AddInt64(&jobs, 1)
				tt.True(pool.Do(i*10 + j))
			}
		}(i)
	}

	time.Sleep(2 * time.Second)
	pool.Close()
	time.Sleep(2 * time.Second)
	t.Log(pool.Info())
	t.Log(atomic.LoadInt64(&jobs) - atomic.LoadInt64(&done))

	tt.False(pool.Do(123))
}
