package types

import (
	"testing"

	"github.com/cosiner/golib/test"
)

func testInOrNot(t *testing.T) {
	test.AssertEq(t, "BitIn", uint(1<<2), BitIn(2, uint((1<<0)|(1<<1)|(1<<2)|(1<<3)|(1<<4))))
	test.AssertEq(t, "BitNotIn", uint(1<<2), BitNotIn(2, uint((1<<0)|(1<<1)|(1<<3)|(1<<4))))
}

func TestBitSet(t *testing.T) {
	tt := test.WrapTest(t)
	bs := NewBitSet(23)
	tt.AssertFalse("BitSet", bs.IsSet(1))

	bs.Set(23) //23
	tt.AssertTrue("BitSet", bs.IsSet(23))
	tt.AssertFalse("BitSet", bs.IsSet(24))
	tt.AssertEq("BitSetCap", uint(64), bs.Cap())
	tt.AssertEq("BitSetLen", uint(24), bs.Len())

	b := bs.Clone() //23
	tt.AssertTrue("BitSet", b.IsSet(23))
	tt.AssertFalse("BitSet", b.IsSet(24))
	tt.AssertEq("BitSetCap", uint(64), b.Cap())
	tt.AssertEq("BitSetLen", uint(24), b.Len())

	b.Set(20)   //23,20
	bs.Union(b) // 23,20
	bs.Set(21)  // 23 21 20
	tt.AssertTrue("BitSetUnion", bs.IsSet(20))

	bs.Diff(b) //21
	tt.AssertFalse("BitSetDiff", bs.IsSet(20))
	tt.AssertTrue("BitSetDiff", bs.IsSet(21))

	tt.AssertFalse("BitSetIntersection", bs.IsSet(23))
}
