package validate

import (
	"testing"

	. "github.com/cosiner/gohper/lib/errors"
	"github.com/cosiner/gohper/lib/test"
)

func TestSimpleEmail(t *testing.T) {
	tt := test.Wrap(t)
	se := &SimpleEmail{Err: Err("Wrong email")}
	tt.Log(se.Validate("11@1.a"))
}

func BenchmarkSimpleEmail(b *testing.B) {
	se := Use(
		Length{
			Min: 3,
			Max: 10,
			Err: Err("aa"),
		}.Validate,
		SimpleEmail{
			Err: Err("Wrong email"),
		}.Validate,
		Chars{
			Chars: "111@1.a",
			Err:   Err("aa"),
		}.Validate)
	for i := 0; i < b.N; i++ {
		_ = se("11@1.a")
	}
}
