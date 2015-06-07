package validate

import (
	"testing"

	"github.com/cosiner/gohper/errors"
	"github.com/cosiner/gohper/testing2"
)

func TestFunc(t *testing.T) {
	tt := testing2.Wrap(t)
	errLength := errors.Err("Wrong length")
	errChars := errors.Err("Wrong chars")
	length := Length{
		Min: 3, Max: 10, Err: errLength,
	}.Validate
	chars := Chars{
		Chars: "0123456789", Err: errChars,
	}.Validate

	vc := UseMul(length, chars)
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
	vc = UseMul(length, Use(length, chars))
	tt.Eq(vc("012", "a"), errLength)
	tt.Eq(vc("012", "0"), errLength)
	tt.Eq(vc("012", "0", "1111111111111111"), errLength)

	vc = UseStrictMul(length, chars)
	defer tt.Recover()
	vc("012", "0", "1111111111111111")
}
