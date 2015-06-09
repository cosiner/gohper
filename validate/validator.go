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
