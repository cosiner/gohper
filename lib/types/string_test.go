package types

import (
	"testing"

	"github.com/cosiner/gohper/lib/test"
)

func TestSnakeString(t *testing.T) {
	test.AssertEq(t, "_xy_xy", SnakeString("_xy_xy"))
	test.AssertEq(t, "_xy_xy", SnakeString("_xyXy"))
	test.AssertEq(t, "_xy xy", SnakeString("_Xy Xy"))
	test.AssertEq(t, "_xy_xy", SnakeString("_Xy_Xy"))
}

func TestCamelString(t *testing.T) {
	test.AssertEq(t, "XyXy", CamelString("xy_xy"))
	test.AssertEq(t, "Xy__Xy", CamelString("xy__Xy"))
	test.AssertEq(t, "Xy Xy", CamelString("xy Xy"))
	test.AssertEq(t, "XY Xy", CamelString("x_y Xy"))
	test.AssertEq(t, "X_Y XY", CamelString("x__y XY"))
	test.AssertEq(t, "XY XY", CamelString("x_y xY"))
	test.AssertEq(t, "XY XY", CamelString("x_y _x_y"))
	test.AssertEq(t, "  XY", CamelString("  x_y"))
}

func TestAbridgeString(t *testing.T) {
	tt := test.WrapTest(t)

	tt.AssertEq("ABC", AbridgeString("AaaBbbCcc"))
	tt.AssertEq("ABC", AbridgeString("AaaBbbCcc"))
}

func TestTrimQuote(t *testing.T) {
	tt := test.WrapTest(t)
	s, err := TrimQuote("\"aaa\"")
	tt.AssertEq("aaa", s)
	tt.AssertEq(err, nil)
}

func TestStringIndexN(t *testing.T) {
	tt := test.WrapTest(t)
	tt.AssertEq(3, StrIndexN("123123123", "12", 2))
	tt.AssertEq(6, StrIndexN("123123123", "12", 3))
	tt.AssertEq(-1, StrIndexN("123123123", "12", 4))
}

func TestStringLastIndexN(t *testing.T) {
	tt := test.WrapTest(t)
	tt.AssertEq(6, StrLastIndexN("123123123", "12", 1))
	tt.AssertEq(3, StrLastIndexN("123123123", "12", 2))
	tt.AssertEq(0, StrLastIndexN("123123123", "12", 3))
	tt.AssertEq(-1, StrLastIndexN("123123123", "12", 4))
}

func TestRepeatJoin(t *testing.T) {
	tt := test.WrapTest(t)
	tt.Log(RepeatJoin("abc", "=?", 10))
}
