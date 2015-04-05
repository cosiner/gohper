package validate

import (
	"testing"

	. "github.com/cosiner/gohper/lib/errors"
	"github.com/cosiner/gohper/lib/test"
)

func TestFunc(t *testing.T) {
	tt := test.Wrap(t)
	errLength := Err("Wrong length")
	errChars := Err("Wrong chars")
	length := Length{
		Min: 3, Max: 10, Err: errLength,
	}
	chars := Chars{
		Chars: "0123456789", Err: errChars,
	}

	vc := UseMul(length.Validate, chars.Validate)
	tt.Eq(vc("0"), errLength)
	tt.Eq(vc("0000000000000000000000"), errLength)
	tt.Eq(vc("000"), nil)

	tt.Eq(vc("abcde"), errChars)

	// length process first, chars process remains
	tt.Eq(vc("01", "abc"), errLength)
	tt.Eq(vc("012", "abc"), errChars)
	tt.Eq(vc("012", "a"), errChars)
	tt.Eq(vc("012", "0"), nil)
	tt.Eq(vc("012", "0", "1111111111111111"), nil)

	// length process first, chars process remains
	vc = UseMul(length.Validate, Use(length.Validate, chars.Validate))
	tt.Eq(vc("012", "a"), errLength)
	tt.Eq(vc("012", "0"), errLength)
	tt.Eq(vc("012", "0", "1111111111111111"), errLength)
}
