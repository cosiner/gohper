package bitset

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestBitSet(t *testing.T) {
	tt := testing2.Wrap(t)
	bs := NewBitset(23)
	tt.False(bs.IsSet(1))

	bs.Set(23) //23
	tt.True(bs.IsSet(23))
	tt.False(bs.IsSet(24))

	b := bs.Clone() //23
	tt.True(b.IsSet(23))
	tt.False(b.IsSet(24))
	tt.Eq(uint(24), b.Length(0))

	b.Set(20)   //23,20
	bs.Union(b) // 23,20
	bs.Set(21)  // 23 21 20
	tt.True(bs.IsSet(20))

	bs.Diff(b) //21
	tt.False(bs.IsSet(20))
	tt.True(bs.IsSet(21))

	tt.False(bs.IsSet(23))

	tt.Eq(uint(30), b.Length(30))
}
