package types

import (
	"github.com/cosiner/golib/test"
	"testing"
)

func TestLowerCase(t *testing.T) {

}

func TestUpperCase(t *testing.T) {

}

func TestTrimSpace(t *testing.T) {

}

func TestTrimUpper(t *testing.T) {

}

func TestTrimLower(t *testing.T) {

}

func TestBytesTrim2Str(t *testing.T) {

}

func TestStr2Bool(t *testing.T) {

}

func TestStr2BoolDef(t *testing.T) {

}

func TestTrimSplit(t *testing.T) {

}

func TestBytesTrimSplit(t *testing.T) {

}

func TestIsSpaceQuote(t *testing.T) {

}

func TestIsSpace(t *testing.T) {

}

func TestEndwith(t *testing.T) {

}

func TestStartWith(t *testing.T) {

}

func TestRepeatJoin(t *testing.T) {

}

func TestFindStringIn(t *testing.T) {
	test.AssertEq(t, 1, FindStringIn([]string{"a", "b", "c"}, "b"), "FindStringIn")
}

func TestSnakeString(t *testing.T) {
	test.AssertEq(t, "_xy_xy", SnakeString("_xy_xy"), "SnakeString1")
	test.AssertEq(t, "_xy_xy", SnakeString("_xyXy"), "SnakeString2")
	test.AssertEq(t, "_xy xy", SnakeString("_Xy Xy"), "SnakeString3")
	test.AssertEq(t, "_xy_xy", SnakeString("_Xy_Xy"), "SnakeString4")
}

func TestCamelString(t *testing.T) {
	test.AssertEq(t, "XyXy", CamelString("xy_xy"), "CamelString1")
	test.AssertEq(t, "Xy__Xy", CamelString("xy__Xy"), "CamelString2")
	test.AssertEq(t, "Xy Xy", CamelString("xy Xy"), "CamelString3")
	test.AssertEq(t, "XY Xy", CamelString("x_y Xy"), "CamelString4")
	test.AssertEq(t, "X_Y XY", CamelString("x__y XY"), "CamelString5")
	test.AssertEq(t, "XY XY", CamelString("x_y xY"), "CamelString6")
	test.AssertEq(t, "XY XY", CamelString("x_y _x_y"), "CamelString7")
	test.AssertEq(t, "  XY", CamelString("  x_y"), "CamelString8")
}
