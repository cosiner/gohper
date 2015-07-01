package attrs

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestAttrs(t *testing.T) {
	testAttrs(t, New())
	testAttrs(t, NewLocked())
}

func testAttrs(t *testing.T, attrs Attrs) {
	tt := testing2.Wrap(t)

	kvs := map[string]interface{}{
		"A": 1,
		"B": 2,
		"C": 3,
		"D": 4,
		"E": 5,
		"F": 6,
	}

	for k, v := range kvs {
		attrs.SetAttr(k, v)
	}
	attrs.SetAttr("A", 7)

	tt.Nil(attrs.Attr("G"))

	for k, v := range kvs {
		val := attrs.Attr(k).(int)
		if k == "A" {
			v = 7
		}

		tt.Eq(v.(int), val)
	}

	tt.Eq(7, attrs.GetSetAttr("A", nil).(int))
	tt.False(attrs.IsAttrExist("A"))

	delete(kvs, "A")
	tt.DeepEq(Values(kvs), attrs.AllAttrs())

	attrs.Clear()
	tt.DeepEq(Values(map[string]interface{}{}), attrs.AllAttrs())
}
