package bytes

import (
	"testing"

	"github.com/cosiner/gohper/lib/test"
)

func TestSyncPool(t *testing.T) {
	poolTest(t, NewSyncPool(4096, true))
}

func TestListPool(t *testing.T) {
	poolTest(t, NewListPool(4096, true))
}

func TestListSlotPool(t *testing.T) {
	poolTest(t, ListSlotPool(8, 4096, true))
}

func TestSyncSlotPool(t *testing.T) {
	poolTest(t, SyncSlotPool(1024, 4096, true))
}

func poolTest(ts testing.TB, pool Pool) {
	t := test.Wrap(ts)
	pool.Init()

	buf := pool.Get(1024, true)
	t.Eq(1024, len(buf))
	t.Eq(4096, cap(buf))
	pool.Put(buf)

	buf = pool.Get(10240, false)
	t.Eq(0, len(buf))
	t.Eq(10240, cap(buf))
	pool.Put(buf)

	t.False(pool.CheckPut(nil))
	t.False(pool.CheckPut(make([]byte, 1024)))
	t.True(pool.CheckPut(make([]byte, 2048)))

	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
	t.Eq(10, len(pool.Get(10, true)))
}
