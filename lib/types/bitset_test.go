package types

import (
	"testing"

	"github.com/cosiner/gohper/lib/test"
)

func testInOrNot(t *testing.T) {
	test.Eq(t, uint(1<<2), BitIn(2, uint((1<<0)|(1<<1)|(1<<2)|(1<<3)|(1<<4))))
	test.Eq(t, uint(1<<2), BitNotIn(2, uint((1<<0)|(1<<1)|(1<<3)|(1<<4))))
}

func TestBitSet(t *testing.T) {
	tt := test.Wrap(t)
	bs := NewBitSet(23)
	tt.False(bs.IsSet(1))

	bs.Set(23) //23
	tt.True(bs.IsSet(23))
	tt.False(bs.IsSet(24))
	tt.Eq(uint(64), bs.Cap())
	tt.Eq(uint(24), bs.Len())

	b := bs.Clone() //23
	tt.True(b.IsSet(23))
	tt.False(b.IsSet(24))
	tt.Eq(uint(64), b.Cap())
	tt.Eq(uint(24), b.Len())

	b.Set(20)   //23,20
	bs.Union(b) // 23,20
	bs.Set(21)  // 23 21 20
	tt.True(bs.IsSet(20))

	bs.Diff(b) //21
	tt.False(bs.IsSet(20))
	tt.True(bs.IsSet(21))

	tt.False(bs.IsSet(23))
}
