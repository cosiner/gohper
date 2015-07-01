package validate

import (
	"regexp"
	"strings"

	"github.com/cosiner/gohper/strings2"
)

type Length struct {
	Min, Max int
	Err      error
}

func (length Length) Validate(s string) error {
	if l := len(s); (length.Min > 0 && l < length.Min) || (length.Max > 0 && l > length.Max) {
		return length.Err
	}

	return nil
}

func ValidLength(min, max int, err error) Validator {
	return Length{
		Min: min,
		Max: max,
		Err: err,
	}.Validate
}

type Chars struct {
	// Chars should be sorted in ascending
	Chars string
	Err   error
}

func (c Chars) Validate(s string) error {
	if !strings2.IsAllCharsIn(s, c.Chars) {
		return c.Err
	}

	return nil
}

func ValidChars(chars string, err error) Validator {
	return Chars{
		Chars: chars,
		Err:   err,
	}.Validate
}

type Regexp struct {
	Regexp *regexp.Regexp
	Err    error
}

func (r Regexp) Validate(s string) error {
	if !r.Regexp.MatchString(s) {
		return r.Err
	}

	return nil
}

func ValidRegexp(reg *regexp.Regexp, err error) Validator {
	return Regexp{
		Regexp: reg,
		Err:    err,
	}.Validate
}

// SimpleEmail only check '@'' and '.' character
type SimpleEmail struct {
	Err error
}

func (e SimpleEmail) Validate(s string) error {
	at := strings.IndexByte(s, '@')
	if at <= 0 || at == len(s)-1 {
		return e.Err
	}

	s = s[at+1:]
	dot := strings.IndexByte(s, '.')
	if dot <= 0 || dot == len(s)-1 {
		return e.Err
	}

	return nil
}

func ValidSimpleEmail(err error) Validator {
	return SimpleEmail{
		Err: err,
	}.Validate
}
