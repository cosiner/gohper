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
	tt.True(UnmarshalPrimitive(bs, reflect.ValueOf(&b)) == nil)
	tt.True(b)
}
