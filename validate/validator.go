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
	if l := len(s); l >= length.Min && l <= length.Max {
		return nil
	}
	return length.Err
}

type Chars struct {
	// Chars should be sorted in ascending
	Chars string
	Err   error
}

func (c Chars) Validate(s string) error {
	if strings2.IsAllCharsIn(s, c.Chars) {
		return nil
	}
	return c.Err
}

type Regexp struct {
	Regexp *regexp.Regexp
	Err    error
}

func (r Regexp) Validate(s string) error {
	if r.Regexp.MatchString(s) {
		return nil
	}
	return r.Err
}

// SimpleEmail only check '@'' and '.' character
type SimpleEmail struct {
	Err error
}

func (e SimpleEmail) Validate(s string) error {
	at := strings.IndexByte(s, '@')
	if at > 0 && at < len(s)-1 {
		s = s[at+1:]
		dot := strings.IndexByte(s, '.')
		if dot > 0 && dot < len(s)-1 {
			return nil
		}
	}
	return e.Err
}
