package defval

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestCond(t *testing.T) {
	tt := testing2.Wrap(t)
	tt.True(Cond(true).String("1", "2") == "1")
	tt.True(Cond(true).Int(1, 2) == 1)
	tt.True(Cond(true).Int8(1, 2) == 1)
	tt.True(Cond(true).Int16(1, 2) == 1)
	tt.True(Cond(true).Int32(1, 2) == 1)
	tt.True(Cond(true).Int64(1, 2) == 1)
	tt.True(Cond(true).Uint(1, 2) == 1)
	tt.True(Cond(true).Uint8(1, 2) == 1)
	tt.True(Cond(true).Uint16(1, 2) == 1)
	tt.True(Cond(true).Uint32(1, 2) == 1)
	tt.True(Cond(true).Uint64(1, 2) == 1)

	tt.True(Cond(false).String("1", "2") == "2")
	tt.True(Cond(false).Int(1, 2) == 2)
	tt.True(Cond(false).Int8(1, 2) == 2)
	tt.True(Cond(false).Int16(1, 2) == 2)
	tt.True(Cond(false).Int32(1, 2) == 2)
	tt.True(Cond(false).Int64(1, 2) == 2)
	tt.True(Cond(false).Uint(1, 2) == 2)
	tt.True(Cond(false).Uint8(1, 2) == 2)
	tt.True(Cond(false).Uint16(1, 2) == 2)
	tt.True(Cond(false).Uint32(1, 2) == 2)
	tt.True(Cond(false).Uint64(1, 2) == 2)
}
