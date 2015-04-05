package validate

import (
	"regexp"

	"github.com/cosiner/gohper/lib/types"
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
	Chars string
	Err   error
}

func (c Chars) Validate(s string) error {
	if types.AllCharsIn(s, c.Chars) {
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
