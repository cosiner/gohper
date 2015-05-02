package reflect

import (
	"reflect"
	"testing"

	"github.com/cosiner/gohper/lib/test"
)

func TestUnmarshalPrimitive(t *testing.T) {
	tt := test.Wrap(t)
	bs := "tzz"
	var b bool
	tt.True(UnmarshalPrimitive(bs, reflect.ValueOf(&b)) == nil)
	tt.True(b)
}

func TestMarshalStruct(t *testing.T) {
	tt := test.Wrap(t)

	st := struct {
		Name string `fd:"nm"`
		Age  int    `fd:"Age"`
	}{
		"aaa",
		123,
	}
	mp := make(map[string]string)
	MarshalStruct(&st, mp, "fd")

	tt.DeepEq(map[string]string{
		"nm":  "aaa",
		"Age": "123",
	}, mp)

	mp["nm"] = "bbb"
	mp["age"] = "234"
	UnmarshalStruct(&st, mp, "fd")

	tt.Eq("bbb", st.Name)
	tt.Eq(123, st.Age)
}
