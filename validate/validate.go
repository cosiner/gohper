package validate

// ValidateV means one validator process all string, next process all string
// ValidateS means all validator process one string, all process next string
// ValidateM means one validator process one string, next process next string,
// remains validator process last string
// *2Last means validate until last string, use last error

type Preprocessor func(string) string
type PreChain []Preprocessor

type Validator func(string) error
type ValidChain []Validator

func (pc PreChain) Process(s string) string {
	for i := 0; i < len(pc); i++ {
		s = pc[i](s)
	}
	return s
}

func New(vc ...Validator) ValidChain {
	return ValidChain(vc)
}

func NewPre(pc ...Preprocessor) PreChain {
	return PreChain(pc)
}

func Pre(pc ...Preprocessor) Preprocessor {
	return PreChain(pc).Process
}

func PreUse(p Preprocessor, c Validator) Validator {
	return func(s string) error {
		return c(p(s))
	}
}

func Use(vc ...Validator) Validator {
	return New(vc...).Validate
}

func UseMul(vc ...Validator) func(...string) error {
	return New(vc...).ValidateM
}

// func (vc ValidChain) Add(vs ...Validator) ValidChain {
// 	if len(vc) == 0 {
// 		return ValidChain(vs)
// 	}
// 	return append(vc, vs...)
// }

// Validate validate string with validators, return first error or nil
func (vc ValidChain) Validate(s string) error {
	for _, v := range vc {
		if e := v(s); e != nil {
			return e
		}
	}
	return nil
}

// Validate2Last validate string with validators, return last error or nil
func (vc ValidChain) Validate2Last(s string) error {
	var err error
	for _, v := range vc {
		if e := v(s); e != nil {
			err = e
		}
	}
	return err
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

// ValidateM2Last validate multiple string with validators, first validator process first string,
// second process next string, etc.., return last error or nil
func (vc ValidChain) ValidateM2Last(s ...string) error {
	var err error
	if i, last := 0, len(s)-1; last > -1 {
		for _, v := range vc {
			if i < last {
				if e := v(s[i]); e != nil {
					err = e
				}
				i++
			} else {
				if e := v(s[last]); e != nil {
					err = e
				}
			}
		}
	}
	return err
}

func (v Validator) ValidateV(s ...string) error {
	for i := 0; i < len(s); i++ {
		if e := v(s[i]); e != nil {
			return e
		}
	}
	return nil
}

func (v Validator) ValidateV2Last(s ...string) error {
	var err error
	for i := 0; i < len(s); i++ {
		if e := v(s[i]); e != nil {
			err = e
		}
	}
	return err
}

// ValidateV validate multiple string with validators, first validator process all string,
// then next validator, etc.., return first error or nil
func (vc ValidChain) ValidateV(s ...string) error {
	for _, v := range vc {
		if e := v.ValidateV(s...); e != nil {
			return e
		}
	}
	return nil
}

// ValidateV2Last validate multiple string with validators, first validator process all string,
// then next validator, etc.., return last error or nil
func (vc ValidChain) ValidateV2Last(s ...string) error {
	var err error
	for _, v := range vc {
		if e := v.ValidateV2Last(s...); e != nil {
			err = e
		}
	}
	return err
}

// ValidateS validate multiple string with validators, all validator process first string,
// then next string, etc.., return first error or nil
func (vc ValidChain) ValidateS(s ...string) error {
	return Validator(vc.Validate).ValidateV(s...)
}

// ValidateS2Last validate multiple string with validators, all validator process first string,
// then next string, etc.., return last error or nil
func (vc ValidChain) ValidateS2Last(s ...string) error {
	return Validator(vc.Validate2Last).ValidateV2Last(s...)
}
