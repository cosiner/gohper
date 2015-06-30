package validate

import (
	"github.com/cosiner/gohper/errors"
)

const (
	ErrParamsCountNotMatch = errors.Err("parameters count not matched with validators")
)

type Validator func(string) error
type ValidChain []Validator

func (v Validator) ValidateM(s ...string) error {
	var err error
	for i, l := 0, len(s); err == nil && i < l; i++ {
		err = v(s[i])
	}

	return err
}

func New(validators ...Validator) ValidChain {
	return ValidChain(validators)
}

func Use(vc ...Validator) Validator {
	return New(vc...).Validate
}

func UseMul(vc ...Validator) func(...string) error {
	return New(vc...).ValidateM
}

func UseStrictMul(vc ...Validator) func(...string) error {
	return New(vc...).StrictValidateM
}

// Validate string with validators, return first error or nil
func (vc ValidChain) Validate(s string) error {
	for _, v := range vc {
		if e := v(s); e != nil {
			return e
		}
	}

	return nil
}

func (vc ValidChain) StrictValidateM(s ...string) error {
	l1, l2 := len(vc)-1, len(s)-1
	errors.Assert(l1 == l2, ErrParamsCountNotMatch)

	var err error
	for i := 0; err == nil && i <= l1; i++ {
		err = vc[i](s[i])
	}

	return err
}

func (vc ValidChain) ValidateM(s ...string) error {
	l1, l2 := len(vc)-1, len(s)-1
	var (
		i   int
		err error
	)
	if l1 < 0 || l2 < 0 {
	} else if l1 <= l2 {
		for i = 0; err == nil && i <= l1; i++ {
			err = vc[i](s[i])
		}

		for i = l1 + 1; err == nil && i <= l2; i++ {
			err = vc[l1](s[i])
		} // last validator process all remains string
	} else {
		for i = 0; err == nil && i <= l2; i++ {
			err = vc[i](s[i])
		}

		for i = l2 + 1; err == nil && i <= l1; i++ {
			err = vc[i](s[l2])
		} // remains validators process last string
	}
	return err
}
