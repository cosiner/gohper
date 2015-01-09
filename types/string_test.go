package types

import (
	"testing"

	"github.com/cosiner/golib/test"
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

func TestSnakeString(t *testing.T) {
	test.AssertEq(t, "SnakeString1", "_xy_xy", SnakeString("_xy_xy"))
	test.AssertEq(t, "SnakeString2", "_xy_xy", SnakeString("_xyXy"))
	test.AssertEq(t, "SnakeString3", "_xy xy", SnakeString("_Xy Xy"))
	test.AssertEq(t, "SnakeString4", "_xy_xy", SnakeString("_Xy_Xy"))
}

func TestCamelString(t *testing.T) {
	test.AssertEq(t, "CamelString1", "XyXy", CamelString("xy_xy"))
	test.AssertEq(t, "CamelString2", "Xy__Xy", CamelString("xy__Xy"))
	test.AssertEq(t, "CamelString3", "Xy Xy", CamelString("xy Xy"))
	test.AssertEq(t, "CamelString4", "XY Xy", CamelString("x_y Xy"))
	test.AssertEq(t, "CamelString5", "X_Y XY", CamelString("x__y XY"))
	test.AssertEq(t, "CamelString6", "XY XY", CamelString("x_y xY"))
	test.AssertEq(t, "CamelString7", "XY XY", CamelString("x_y _x_y"))
	test.AssertEq(t, "CamelString8", "  XY", CamelString("  x_y"))
}
