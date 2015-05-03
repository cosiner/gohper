package reflect2

import (
	"reflect"
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestUnmarshalPrimitive(t *testing.T) {
	tt := testing2.Wrap(t)
	bs := "tzz"
	var b bool
	tt.True(UnmarshalPrimitive(bs, reflect.ValueOf(&b)) == nil)
	tt.True(b)
}

func TestMarshalStruct(t *testing.T) {
	tt := testing2.Wrap(t)

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
