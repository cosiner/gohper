package types

import (
	"testing"

	"github.com/cosiner/gohper/lib/test"
)

func TestSnakeString(t *testing.T) {
	test.Eq(t, "_xy_xy", SnakeString("_xy_xy"))
	test.Eq(t, "_xy_xy", SnakeString("_xyXy"))
	test.Eq(t, "_xy xy", SnakeString("_Xy Xy"))
	test.Eq(t, "_xy_xy", SnakeString("_Xy_Xy"))
}

func TestCamelString(t *testing.T) {
	test.Eq(t, "XyXy", CamelString("xy_xy"))
	test.Eq(t, "Xy__Xy", CamelString("xy__Xy"))
	test.Eq(t, "Xy Xy", CamelString("xy Xy"))
	test.Eq(t, "XY Xy", CamelString("x_y Xy"))
	test.Eq(t, "X_Y XY", CamelString("x__y XY"))
	test.Eq(t, "XY XY", CamelString("x_y xY"))
	test.Eq(t, "XY XY", CamelString("x_y _x_y"))
	test.Eq(t, "  XY", CamelString("  x_y"))
}

func TestAbridgeString(t *testing.T) {
	tt := test.Wrap(t)

	tt.Eq("ABC", AbridgeString("AaaBbbCcc"))
	tt.Eq("ABC", AbridgeString("AaaBbbCcc"))
}

func TestTrimQuote(t *testing.T) {
	tt := test.Wrap(t)
	s, err := TrimQuote("\"aaa\"")
	tt.Eq("aaa", s)
	tt.Eq(err, nil)
}

func TestStringIndexN(t *testing.T) {
	tt := test.Wrap(t)
	tt.Eq(3, StrIndexN("123123123", "12", 2))
	tt.Eq(6, StrIndexN("123123123", "12", 3))
	tt.Eq(-1, StrIndexN("123123123", "12", 4))
}

func TestStringLastIndexN(t *testing.T) {
	tt := test.Wrap(t)
	tt.Eq(6, StrLastIndexN("123123123", "12", 1))
	tt.Eq(3, StrLastIndexN("123123123", "12", 2))
	tt.Eq(0, StrLastIndexN("123123123", "12", 3))
	tt.Eq(-1, StrLastIndexN("123123123", "12", 4))
}

func TestRepeatJoin(t *testing.T) {
	tt := test.Wrap(t)
	tt.Log(RepeatJoin("abc", "=?", 10))
}

func TestValid(t *testing.T) {
	tt := test.Wrap(t)
	tt.True(AllCharsIn("", "abcdefghijklmn"))
	tt.True(AllCharsIn("abc", "abcdefghijklmn"))
	tt.False(AllCharsIn("ao", "abcdefghijklmn"))
}

func TestRemoveSpace(t *testing.T) {
	tt := test.Wrap(t)
	tt.Eq("abcdefg", RemoveSpace(`a b
    	c d 	e
    	 	f g`))
}
