package validate

import (
	"regexp"
	"testing"

	"github.com/cosiner/gohper/errors"
	"github.com/cosiner/gohper/testing2"
)

func TestFunc(t *testing.T) {
	tt := testing2.Wrap(t)

	errReg := errors.Err("regexp not match")
	r := Regexp{
		Regexp: regexp.MustCompile("^\\d+$"),
		Err:    errReg,
	}
	tt.Nil(r.Validate("123"))
	tt.Eq(errReg, r.Validate("123dqe"))

	errLength := errors.Err("Wrong length")
	errChars := errors.Err("Wrong chars")
	length := Length{
		Min: 3, Max: 10, Err: errLength,
	}.Validate
	chars := Chars{
		Chars: "0123456789", Err: errChars,
	}.Validate

	single := Use(length, chars)
	tt.Nil(single("023456"))

	vc := UseMul(length, chars)
	tt.Eq(vc("0"), errLength)
	tt.Eq(vc("0000000000000000000000"), errLength)
	tt.Nil(vc("000"))

	tt.Eq(vc("abcde"), errChars)

	// length process first, chars process remains
	tt.Eq(vc("01", "abc"), errLength)
	tt.Eq(vc("012", "abc"), errChars)
	tt.Eq(vc("012", "a"), errChars)
	tt.Nil(vc("012", "0"))
	tt.Nil(vc("012", "0", "1111111111111111"))

	// length process first, chars process remains
	vc = UseMul(length, Use(length, chars))
	tt.Eq(vc("012", "a"), errLength)
	tt.Eq(vc("012", "0"), errLength)
	tt.Eq(vc("012", "0", "1111111111111111"), errLength)

	vc = UseStrictMul(length, chars)
	tt.Nil(vc("abcd", "1234"))

	defer tt.Recover()
	vc("012", "0", "1111111111111111")
}

func TestSimpleEmail(t *testing.T) {
	tt := testing2.Wrap(t)
	err := errors.Err("Wrong email")
	se := &SimpleEmail{Err: err}

	tt.Nil(se.Validate("11@1.a"))
	tt.Eq(err, se.Validate("11@1."))
	tt.Eq(err, se.Validate("@1.a"))
	tt.Eq(err, se.Validate("11@.11"))
	tt.Eq(err, se.Validate("11@"))
}

func TestNewValidator(t *testing.T) {
	_ = ValidLength(1, 2, nil)
	_ = ValidChars("123", nil)
	_ = ValidRegexp(nil, nil)
	_ = ValidSimpleEmail(nil)
}
