package types

import (
	"reflect"
	"testing"

	"github.com/cosiner/golib/test"
)

func TestUnmarshalPrimitive(t *testing.T) {
	tt := test.WrapTest(t)
	bs := []byte("true")
	var b bool
	tt.AssertNil("1", UnmarshalPrimitive(bs, reflect.ValueOf(&b)))
	tt.AssertTrue("2", b)
}
