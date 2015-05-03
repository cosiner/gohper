package strings2

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestSnakeCase(t *testing.T) {
	testing2.Eq(t, "_xy_xy", ToSnake("_xy_xy"))
	testing2.Eq(t, "_xy_xy", ToSnake("_xyXy"))
	testing2.Eq(t, "_xy xy", ToSnake("_Xy Xy"))
	testing2.Eq(t, "_xy_xy", ToSnake("_Xy_Xy"))
}

func TestCamelString(t *testing.T) {
	testing2.Eq(t, "XyXy", ToCamel("xy_xy"))
	testing2.Eq(t, "Xy__Xy", ToCamel("xy__Xy"))
	testing2.Eq(t, "Xy Xy", ToCamel("xy Xy"))
	testing2.Eq(t, "XY Xy", ToCamel("x_y Xy"))
	testing2.Eq(t, "X_Y XY", ToCamel("x__y XY"))
	testing2.Eq(t, "XY XY", ToCamel("x_y xY"))
	testing2.Eq(t, "XY XY", ToCamel("x_y _x_y"))
	testing2.Eq(t, "  XY", ToCamel("  x_y"))
}

func TestAbridgeString(t *testing.T) {
	tt := testing2.Wrap(t)

	tt.Eq("ABC", ToAbridge("AaaBbbCcc"))
	tt.Eq("ABC", ToAbridge("AaaBbbCcc"))
}

func TestTrimQuote(t *testing.T) {
	tt := testing2.Wrap(t)
	s, err := TrimQuote("\"aaa\"")
	tt.Eq("aaa", s)
	tt.Eq(err, nil)
}

func TestSplitAtN(t *testing.T) {
	tt := testing2.Wrap(t)
	tt.Eq(3, SplitAtN("123123123", "12", 2))
	tt.Eq(6, SplitAtN("123123123", "12", 3))
	tt.Eq(-1, SplitAtN("123123123", "12", 4))
}

func TestSplitAtLastN(t *testing.T) {
	tt := testing2.Wrap(t)
	tt.Eq(6, SplitAtLastN("123123123", "12", 1))
	tt.Eq(3, SplitAtLastN("123123123", "12", 2))
	tt.Eq(0, SplitAtLastN("123123123", "12", 3))
	tt.Eq(-1, SplitAtLastN("123123123", "12", 4))
}

func TestRepeatJoin(t *testing.T) {
	tt := testing2.Wrap(t)
	tt.Eq("abc=?,abc=?,abc", RepeatJoin("abc", "=?,", 3))
}

func TestValid(t *testing.T) {
	tt := testing2.Wrap(t)
	tt.True(IsAllCharsIn("", "abcdefghijklmn"))
	tt.True(IsAllCharsIn("abc", "abcdefghijklmn"))
	tt.False(IsAllCharsIn("ao", "abcdefghijklmn"))
}

func TestRemoveSpace(t *testing.T) {
	tt := testing2.Wrap(t)
	tt.Eq("abcdefg", RemoveSpace(`a b
    	c d 	e
    	 	f g`))
}
