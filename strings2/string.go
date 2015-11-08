package strings2

import (
	"bytes"
	"strings"

	"github.com/cosiner/gohper/index"
	"github.com/cosiner/gohper/unibyte"
)

// TrimQuote trim quote for string, return error if quote don't match
func TrimQuote(str string) (string, bool) {
	str = strings.TrimSpace(str)
	l := len(str)
	if l == 0 {
		return "", true
	}

	if s, e := str[0], str[l-1]; s == '\'' || s == '"' || s == '`' || e == '\'' || e == '"' || e == '`' {
		if l != 1 && s == e {
			str = str[1 : l-1]
		} else {
			return "", false
		}
	}

	return str, true
}

func TrimWrap(str, left, right string, strict bool) (string, bool) {
	l := len(str)
	lstr := l

	if l == 0 {
		return "", true
	}

	ll := len(left)
	if strings.HasPrefix(str, left) {
		str = str[ll:]
		l -= ll
	}

	lr := len(right)
	if strings.HasSuffix(str, right) {
		str = str[:l-lr]
		l -= lr
	}

	return str, lstr == l || lstr == l+ll+lr || !strict
}

// TrimAndToUpper return the trim and upper format of a string
func TrimAndToUpper(str string) string {
	return strings.ToUpper(strings.TrimSpace(str))
}

// TrimAndToLower return the trim and lower format of a string
func TrimAndToLower(str string) string {
	return strings.ToLower(strings.TrimSpace(str))
}

// TrimSplit split string and return trim space string
func SplitAndTrim(s, sep string) []string {
	sp := strings.Split(s, sep)
	for i, n := 0, len(sp); i < n; i++ {
		sp[i] = strings.TrimSpace(sp[i])
	}

	return sp
}

// TrimAfter trim string and remove the section after delimiter and delimiter itself
func TrimAfter(s, delimiter string) string {
	if idx := strings.Index(s, delimiter); idx >= 0 {
		s = s[:idx]
	}

	return strings.TrimSpace(s)
}

func TrimBefore(s, delimiter string) string {
	if idx := strings.Index(s, delimiter); idx >= 0 {
		s = s[idx+len(delimiter):]
	}

	return strings.TrimSpace(s)
}

// IndexN find index of n-th sep string
func IndexN(str, sep string, n int) (index int) {
	index, idx, sepLen := 0, -1, len(sep)
	for i := 0; i < n; i++ {
		if idx = strings.Index(str, sep); idx == -1 {
			break
		}
		str = str[idx+sepLen:]
		index += idx
	}

	if idx == -1 {
		index = -1
	} else {
		index += (n - 1) * sepLen
	}

	return
}

// LastIndexN find last index of n-th sep string
func LastIndexN(str, sep string, n int) (index int) {
	for i := 0; i < n; i++ {
		if index = strings.LastIndex(str, sep); index == -1 {
			break
		}
		str = str[:index]
	}

	return
}

// Separate string by separator, the separator must in the middle of string,
// not first and last
func Separate(s string, sep byte) (string, string) {
	i := MidIndex(s, sep)
	if i < 0 {
		return "", ""
	}

	return s[:i], s[i+1:]
}

func LastIndexByte(s string, b byte) int {
	for l := len(s) - 1; l >= 0; l-- {
		if s[l] == b {
			return l
		}
	}

	return -1
}

// IsAllCharsIn check whether all chars of string is in encoding string
func IsAllCharsIn(s, encoding string) bool {
	var is = true
	for i := 0; i < len(s) && is; i++ {
		is = index.CharIn(s[i], encoding) >= 0
	}

	return is
}

// MidIndex find middle separator index of string, not first and last
func MidIndex(s string, sep byte) int {
	i := strings.IndexByte(s, sep)
	if i <= 0 || i == len(s)-1 {
		return -1
	}

	return i
}

// RemoveSpace remove all space characters from string by unibyte.IsSpace
func RemoveSpace(s string) string {
	idx, end := 0, len(s)
	bs := make([]byte, end)
	for i := 0; i < end; i++ {
		if !unibyte.IsSpace(s[i]) {
			bs[idx] = s[i]
			idx++
		}
	}

	return string(bs[:idx])
}

// MergeSpace merge multiple space to one, trim determine whether remove space at prefix and suffix
func MergeSpace(s string, trim bool) string {
	space := false

	idx, end := 0, len(s)
	bs := make([]byte, end)
	for i := 0; i < end; i++ {
		if unibyte.IsSpace(s[i]) {
			space = true
		} else {
			if space && (!trim || idx != 0) {
				bs[idx] = ' '
				idx++
			}

			bs[idx] = s[i]
			idx++

			space = false
		}
	}

	if space && !trim {
		bs[idx] = ' '
		idx++
	}

	return string(bs[:idx])
}

// IndexNonSpace find index of first non-space character, if not exist, -1 was returned
func IndexNonSpace(s string) int {
	for i := range s {
		if !unibyte.IsSpace(s[i]) {
			return i
		}
	}

	return -1
}

// LastIndexNonSpace find index of last non-space character, if not exist, -1 was returned
func LastIndexNonSpace(s string) int {
	for i := len(s) - 1; i >= 0; i-- {
		if !unibyte.IsSpace(s[i]) {
			return i
		}
	}

	return -1
}

// WriteStringsToBuffer write strings to buffer, it avoid memory allocation of join
// strings
func WriteStringsToBuffer(buffer *bytes.Buffer, strings []string, sep string) {
	i, last := 0, len(strings)-1
	for ; i < last; i++ {
		buffer.WriteString(strings[i])
		buffer.WriteString(sep)
	}

	if last != -1 {
		buffer.WriteString(strings[last])
	}
}

func MultipleLineOperate(s, delim string, operate func(line, delim string) string) string {
	lines := strings.Split(s, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		lines[i] = operate(lines[i], delim)
	}

	return strings.Join(lines, "\n")
}

func TrimLastN(s, delim string, n int) string {
	s = strings.TrimSpace(s)

	if n <= 0 {
		n = -1
	}
	sl, dl := len(s), len(delim)
	for n != 0 && strings.HasSuffix(s, delim) {
		s = s[:sl-dl]
		sl = len(s)
		n--
	}
	return s
}

func TrimFirstN(s, delim string, n int) string {
	s = strings.TrimSpace(s)

	if n <= 0 {
		n = -1
	}
	dl := len(delim)
	for n != 0 && strings.HasPrefix(s, delim) {
		s = s[dl:]
		n--
	}
	return s
}

func JoinPairs(pairs map[string]string, eq, sep string) string {
	var s string
	for k, v := range pairs {
		if s != "" {
			s += sep
		}
		s += (k + eq + v)
	}

	return s
}
