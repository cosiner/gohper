package types

import (
	"testing"

	"github.com/cosiner/golib/test"
)

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

func TestAbridgeString(t *testing.T) {
	tt := test.WrapTest(t)

	tt.AssertEq("Abr1", "ABC", AbridgeString("AaaBbbCcc"))
	tt.AssertEq("Abr2", "ABC", AbridgeString("AaaBbbCcc"))
}
