package reflect

import (
	"reflect"
	"testing"

	"github.com/cosiner/golib/test"
)

func TestUnmarshalPrimitive(t *testing.T) {
	tt := test.WrapTest(t)
	bs := "true"
	var b bool
	tt.AssertNil(UnmarshalPrimitive(bs, reflect.ValueOf(&b)))
	tt.AssertTrue(b)
}
