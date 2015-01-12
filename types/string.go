package types

import (
	"bytes"
	"strings"
	"unicode"
)

// IsLower check letter is lower case or not
func IsLower(b byte) bool {
	return b >= 'a' && b <= 'z'
}

// IsUpper check letter is upper case or not
func IsUpper(b byte) bool {
	return b >= 'A' && b <= 'Z'
}

func IsLetter(b byte) bool {
	return IsLower(b) || IsUpper(b)
}

// IsSpaceQuote return wehter a byte is space or quote characters
func IsSpaceQuote(b byte) bool {
	return IsSpace(b) || b == '"' || b == '\''
}

// IsSpace only call unicode.IsSpace
func IsSpace(b byte) bool {
	return unicode.IsSpace(rune(b))
}

// LowerCase convert a byte to lower case
func LowerCase(b byte) byte {
	if IsUpper(b) {
		b = b - 'A' + 'a'
	}
	return b
}

// UpperCase convert a byte to upper
func UpperCase(b byte) byte {
	if IsLower(b) {
		b = b - 'a' + 'A'
	}
	return b
}

// TrimSpace is only call strings.TrimSpace
func TrimSpace(str string) string {
	return strings.TrimSpace(str)
}

// TrimUpper return the trim and upper format of a string
func TrimUpper(str string) string {
	return strings.ToUpper(strings.TrimSpace(str))
}

// TrimLower return the trim and lower format of a string
func TrimLower(str string) string {
	return strings.ToLower(strings.TrimSpace(str))
}

// BytesTrim2Str trim bytes return as string
func TrimBytes2Str(s []byte) string {
	return string(bytes.TrimSpace(s))
}

// TrimSplit split string and return trim space string
func TrimSplit(s, sep string) []string {
	sp := strings.Split(s, sep)
	for i, n := 0, len(sp); i < n; i++ {
		sp[i] = strings.TrimSpace(sp[i])
	}
	return sp
}

// BytesTrimSplit split string and return trim space string
func BytesTrimSplit(s, sep []byte) [][]byte {
	sp := bytes.Split(s, sep)
	for i, n := 0, len(sp); i < n; i++ {
		sp[i] = bytes.TrimSpace(sp[i])
	}
	return sp
}

// StartWith check whether str is start with another string
// it's a wrapper of strings.HasPrefix
func StartWith(str, start string) bool {
	return strings.HasPrefix(str, start)
}

// EndWith check whether str is end with another string
// it's a wrapper of strings.HasSuffix
func EndWith(str, end string) bool {
	return strings.HasSuffix(str, end)
}

// RepeatJoin repeat s1 count times, then join with s2
func RepeatJoin(s1, s2 string, count int) string {
	if count <= 0 {
		return ""
	}
	str := strings.Repeat(s1+s2, count)
	return str[:len(str)-len(s2)]
}

// StringReader is a wrapper of strings.NewReader
func StringReader(s string) *strings.Reader {
	return strings.NewReader(s)
}

// TrimAfter trim string and remove the section after delimiter and delimiter itself
func TrimAfter(s string, delimiter string) string {
	idx := strings.Index(s, delimiter)
	if idx >= 0 {
		s = s[:idx]
	}
	return strings.TrimSpace(s)
}

// TrimAfter trim bytes and remove the section after delimiter and delimiter itself
func TrimBytesAfter(s []byte, delimiter []byte) []byte {
	idx := bytes.Index(s, delimiter)
	if idx >= 0 {
		s = s[:idx]
	}
	return bytes.TrimSpace(s)
}

// snake string, XxYy to xx_yy, X_Y to x_y
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	num := len(s)
	need := false // need determin if it's necessery to add a '_'
	for i := 0; i < num; i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c = c - 'A' + 'a'
			if need {
				data = append(data, '_')
			}
		} else {
			// if previous is '_' or ' ',
			// there is no need to add extra '_' before
			need = (c != '_' && c != ' ')
		}
		data = append(data, c)
	}
	return string(data)
}

// camel string, xx_yy to XxYy, xx__yy to Xx_Yy
// xx _yy to Xx Yy, the rule is that a lower case letter
// after '_' will combine to a upper case letter
func CamelString(s string) string {
	data := make([]byte, 0, len(s))
	num := len(s)
	need := true
	var prev byte = ' '

	for i := 0; i < num; i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			if need {
				c = c - 'a' + 'A'
				need = false
			}
		} else {
			if prev == '_' {
				data = append(data, '_')
			}
			need = (c == '_' || c == ' ')
			if c == '_' {
				prev = '_'
				continue
			}
		}
		prev = c
		data = append(data, c)
	}
	return string(data)
}

// AbridgeString extract first letter and all upper case letter
// from string as it's abridge case
func AbridgeString(str string) (s string) {
	if l := len(str); l != 0 {
		ret := []byte{str[0]}
		for i := 1; i < l; i++ {
			b := str[i]
			if IsUpper(b) {
				ret = append(ret, b)
			}
		}
		s = string(ret)
	}
	return
}

// AbridgeStringToLower extract first letter and all upper case letter
// from string as it's abridge case, and convert it to lower case
func AbridgeStringToLower(str string) (s string) {
	if l := len(str); l != 0 {
		ret := []byte{LowerCase(str[0])}
		for i := 1; i < l; i++ {
			b := str[i]
			if IsUpper(b) {
				ret = append(ret, LowerCase(b))
			}
		}
		s = string(ret)
	}
	return
}
