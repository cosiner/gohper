package strings2

import "github.com/cosiner/gohper/unibyte"

// ToSnake string, XxYy to xx_yy, X_Y to x_y
func ToSnake(s string) string {
	num := len(s)
	need := false // need determin if it's necessery to add a '_'

	snake := make([]byte, 0, len(s)*2)
	for i := 0; i < num; i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c = c - 'A' + 'a'
			if need {
				snake = append(snake, '_')
			}
		} else {
			// if previous is '_' or ' ',
			// there is no need to add extra '_' before
			need = (c != '_' && c != ' ')
		}

		snake = append(snake, c)
	}

	return string(snake)
}

// ToCamel string, xx_yy to XxYy, xx__yy to Xx_Yy
// xx _yy to Xx Yy, the rule is that a lower case letter
// after '_' will combine to a upper case letter
func ToCamel(s string) string {
	num := len(s)
	need := true

	var prev byte = ' '
	camel := make([]byte, 0, len(s))
	for i := 0; i < num; i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			if need {
				c = c - 'a' + 'A'
				need = false
			}
		} else {
			if prev == '_' {
				camel = append(camel, '_')
			}
			need = (c == '_' || c == ' ')
			if c == '_' {
				prev = '_'
				continue
			}
		}

		prev = c
		camel = append(camel, c)
	}

	return string(camel)
}

// ToAbridge extract first letter and all upper case letter
// from string as it's abridge case
func ToAbridge(str string) string {
	l := len(str)
	if l == 0 {
		return ""
	}

	arbi := []byte{str[0]}
	for i := 1; i < l; i++ {
		b := str[i]
		if unibyte.IsUpper(b) {
			arbi = append(arbi, b)
		}
	}

	return string(arbi)
}

// ToLowerAbridge extract first letter and all upper case letter
// from string as it's abridge case, and convert it to lower case
func ToLowerAbridge(str string) (s string) {
	l := len(str)
	if l == 0 {
		return ""
	}

	arbi := []byte{unibyte.ToLower(str[0])}
	for i := 1; i < l; i++ {
		b := str[i]
		if unibyte.IsUpper(b) {
			arbi = append(arbi, unibyte.ToLower(b))
		}
	}

	return string(arbi)
}
