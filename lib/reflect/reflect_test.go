package reflect

import (
	"reflect"
	"testing"

	"github.com/cosiner/gohper/lib/test"
)

func TestUnmarshalPrimitive(t *testing.T) {
	tt := test.Wrap(t)
	bs := "true"
	var b bool
	tt.Nil(UnmarshalPrimitive(bs, reflect.ValueOf(&b)))
	tt.True(b)
}
