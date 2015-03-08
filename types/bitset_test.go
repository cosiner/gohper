package types

import (
	"testing"

	"github.com/cosiner/golib/test"
)

func testInOrNot(t *testing.T) {
	test.AssertEq(t, uint(1<<2), BitIn(2, uint((1<<0)|(1<<1)|(1<<2)|(1<<3)|(1<<4))))
	test.AssertEq(t, uint(1<<2), BitNotIn(2, uint((1<<0)|(1<<1)|(1<<3)|(1<<4))))
}

func TestBitSet(t *testing.T) {
	tt := test.WrapTest(t)
	bs := NewBitSet(23)
	tt.AssertFalse(bs.IsSet(1))

	bs.Set(23) //23
	tt.AssertTrue(bs.IsSet(23))
	tt.AssertFalse(bs.IsSet(24))
	tt.AssertEq(uint(64), bs.Cap())
	tt.AssertEq(uint(24), bs.Len())

	b := bs.Clone() //23
	tt.AssertTrue(b.IsSet(23))
	tt.AssertFalse(b.IsSet(24))
	tt.AssertEq(uint(64), b.Cap())
	tt.AssertEq(uint(24), b.Len())

	b.Set(20)   //23,20
	bs.Union(b) // 23,20
	bs.Set(21)  // 23 21 20
	tt.AssertTrue(bs.IsSet(20))

	bs.Diff(b) //21
	tt.AssertFalse(bs.IsSet(20))
	tt.AssertTrue(bs.IsSet(21))

	tt.AssertFalse(bs.IsSet(23))
}
