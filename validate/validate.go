package validate

type Validator func(string) error
type ValidChain []Validator

func New(validators ...Validator) ValidChain {
	return ValidChain(validators)
}

func Use(vc ...Validator) Validator {
	return New(vc...).Validate
}

func UseMul(vc ...Validator) func(...string) error {
	return New(vc...).ValidateM
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

// ValidateM validate multiple string with validators, first validator process first string,
// second process next string, etc.., return first error or nil
func (vc ValidChain) ValidateM(s ...string) error {
	if i, last := 0, len(s)-1; last > -1 {
		for _, v := range vc {
			if i < last {
				if e := v(s[i]); e != nil {
					return e
				}
				i++
			} else {
				if e := v(s[last]); e != nil {
					return e
				}
			}
		}
	}

	return nil
}
