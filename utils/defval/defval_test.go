package defval

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestDefVal(t *testing.T) {
	tt := testing2.Wrap(t)

	var val int = 0
	Int(&val, 10)
	tt.True(val == 10)

	var val1 int8 = 0
	Int8(&val1, 10)
	tt.True(val1 == 10)

	var val2 int16 = 0
	Int16(&val2, 10)
	tt.True(val2 == 10)

	var val3 int32 = 0
	Int32(&val3, 10)
	tt.True(val3 == 10)

	var val4 int64 = 0
	Int64(&val4, 10)
	tt.True(val4 == 10)

	var val5 uint = 0
	Uint(&val5, 10)
	tt.True(val5 == 10)

	var val6 uint8 = 0
	Uint8(&val6, 10)
	tt.True(val6 == 10)

	var val7 uint16 = 0
	Uint16(&val7, 10)
	tt.True(val7 == 10)

	var val8 uint32 = 0
	Uint32(&val8, 10)
	tt.True(val8 == 10)

	var val9 uint64 = 0
	Uint64(&val9, 10)
	tt.True(val9 == 10)

	var val10 uint = 0
	Uint(&val10, 10)
	tt.True(val10 == 10)

	var val11 string = ""
	String(&val11, "10")
	tt.Eq("10", val11)

	var f func()
	var v bool
	Nil(&f, func() {
		v = true
	})
	f()
	tt.True(v)

	BoolStr(true, &val11)
	tt.Eq("true", val11)
	BoolStr(false, &val11)
	tt.Eq("false", val11)

	BoolInt(true, &val)
	tt.True(1 == val)
	BoolInt(false, &val)
	tt.True(0 == val)
}
