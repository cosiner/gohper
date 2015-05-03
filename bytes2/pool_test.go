package bytes2

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
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

func poolTest(tt testing.TB, pool Pool) {
	t := testing2.Wrap(tt)

	buf := pool.Get(1024, true)
	t.Eq(1024, len(buf))
	t.Eq(4096, cap(buf))
	pool.Put(buf)

	buf = pool.Get(10240, false)
	t.Eq(0, len(buf))
	t.Eq(10240, cap(buf))
	pool.Put(buf)

	t.False(pool.TryPut(nil))
	t.False(pool.TryPut(make([]byte, 1024)))
	t.True(pool.TryPut(make([]byte, 2048)))

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
