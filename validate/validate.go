package validate

type Validator func(string) error

type Chain []Validator

func New(c ...Validator) Chain {
	return Chain(c)
}

func Use(c ...Validator) Validator {
	return New(c...).Validate
}

func UseMul(c ...Validator) func(...string) error {
	return New(c...).ValidateMul
}

func (c Chain) Add(n ...Validator) Chain {
	if len(c) == 0 {
		return Chain(n)
	}
	return append(c, n...)
}

func (v Validator) Validate(s ...string) error {
	for i := 0; i < len(s); i++ {
		if e := v(s[i]); e != nil {
			return e
		}
	}
	return nil
}

func (v Validator) Validate2Last(s ...string) error {
	var err error
	for i := 0; i < len(s); i++ {
		if e := v(s[i]); e != nil {
			err = e
		}
	}
	return err
}

// Validate validate string with validators, return first error or nil
func (c Chain) Validate(s string) error {
	for _, v := range c {
		if e := v(s); e != nil {
			return e
		}
	}
	return nil
}

// Validate validate string with validators, return last error or nil
func (c Chain) Validate2Last(s string) error {
	var err error
	for _, v := range c {
		if e := v(s); e != nil {
			err = e
		}
	}
	return err
}

// ValidateMul validate multiple string with validators, first validator process first string,
// second process next string, etc.., return first error or nil
func (c Chain) ValidateMul(s ...string) error {
	if i, last := 0, len(s)-1; last > -1 {
		for _, v := range c {
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

// ValidateMul2Last validate multiple string with validators, first validator process first string,
// second process next string, etc.., return last error or nil
func (c Chain) ValidateMul2Last(s ...string) error {
	var err error
	if i, last := 0, len(s)-1; last > -1 {
		for _, v := range c {
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

// ValidateMulV validate multiple string with validators, first validator process all string,
// then next validator, etc.., return first error or nil
func (c Chain) ValidateMulV(s ...string) error {
	for _, v := range c {
		if e := v.Validate(s...); e != nil {
			return e
		}
	}
	return nil
}

// ValidateMulV validate multiple string with validators, first validator process all string,
// then next validator, etc.., return last error or nil
func (c Chain) ValidateMul2LastV(s ...string) error {
	var err error
	for _, v := range c {
		if e := v.Validate2Last(s...); e != nil {
			err = e
		}
	}
	return err
}
