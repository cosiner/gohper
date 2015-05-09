package bytes2

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestSyncPool(t *testing.T) {
	poolTest(t, NewSyncPool(4096, true), false)
}

func TestFakePool(t *testing.T) {
	poolTest(t, FakePool{}, true)
}

func TestListPool(t *testing.T) {
	lp := NewListPool(4096, true).(*ListPool)
	poolTest(t, lp, false)
	for i := 0; i < 10; i++ {
		lp.Put(make([]byte, 2048))
	}
	t.Log(lp.Count)
	c := lp.Count
	c--
	lp.ShrinkTo(-1)
	lp.ShrinkTo(c)
}

func TestListSlotPool(t *testing.T) {
	poolTest(t, ListSlotPool(8, 4096, true), false)
}

func TestSyncSlotPool(t *testing.T) {
	poolTest(t, SyncSlotPool(1024, 4096, true), false)
}

func poolTest(tt testing.TB, pool Pool, fake bool) {
	t := testing2.Wrap(tt)

	buf := pool.Get(1024, true)
	t.Eq(1024, len(buf))
	if !fake {
		t.Eq(4096, cap(buf))
	}
	pool.Put(buf)

	buf = pool.Get(10240, false)
	t.Eq(0, len(buf))
	t.Eq(10240, cap(buf))
	pool.Put(buf)

	t.False(pool.TryPut(nil))
	t.False(pool.TryPut(make([]byte, 1024)))

	if !fake {
		t.True(pool.TryPut(make([]byte, 2048)))
	}

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

func TestPanicSlot(t *testing.T) {
	tt := testing2.Wrap(t)
	defer tt.Recover()
	NewSlotPool(nil)
}
