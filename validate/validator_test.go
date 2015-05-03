package validate

import (
	"testing"

	"github.com/cosiner/gohper/errors"
	"github.com/cosiner/gohper/testing2"
)

func TestSimpleEmail(t *testing.T) {
	tt := testing2.Wrap(t)
	se := &SimpleEmail{Err: errors.Err("Wrong email")}
	tt.Log(se.Validate("11@1.a"))
}

func BenchmarkSimpleEmail(b *testing.B) {
	se := Use(
		Length{
			Min: 3,
			Max: 10,
			Err: errors.Err("aa"),
		}.Validate,
		SimpleEmail{
			Err: errors.Err("Wrong email"),
		}.Validate,
		Chars{
			Chars: "111@1.a",
			Err:   errors.Err("aa"),
		}.Validate)
	for i := 0; i < b.N; i++ {
		_ = se("11@1.a")
	}
}
