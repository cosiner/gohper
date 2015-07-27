package strings2

import (
	"bytes"
	"strconv"
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

func TrimBefore(s, delimeter string) string {
	if idx := strings.Index(s, delimeter); idx >= 0 {
		s = s[idx+len(delimeter):]
	}

	return strings.TrimSpace(s)
}

// IndexN find index of n-th sep string
func IndexN(str, sep string, n int) (index int) {
	index, idx, seplen := 0, -1, len(sep)
	for i := 0; i < n; i++ {
		if idx = strings.Index(str, sep); idx == -1 {
			break
		}
		str = str[idx+seplen:]
		index += idx
	}

	if idx == -1 {
		index = -1
	} else {
		index += (n - 1) * seplen
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

// Seperate string by seperator, the seperator must in the middle of string,
// not first and last
func Seperate(s string, sep byte) (string, string) {
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

// MidIndex find middle seperator index of string, not first and last
func MidIndex(s string, sep byte) int {
	i := strings.IndexByte(s, sep)
	if i <= 0 || i == len(s)-1 {
		return -1
	}

	return i
}

// RepeatJoin repeat s count times as a string slice, then join with sep
func RepeatJoin(s, sep string, count int) string {
	switch {
	case count <= 0:
		return ""
	case count == 1:
		return s
	case count == 2:
		return s + sep + s
	default:
		bs := make([]byte, 0, (len(s)+len(sep))*count-len(sep))
		buf := bytes.NewBuffer(bs)
		buf.WriteString(s)

		for i := 1; i < count; i++ {
			buf.WriteString(sep)
			buf.WriteString(s)
		}

		return buf.String()
	}
}

// SuffixJoin join string slice with suffix
func SuffixJoin(s []string, suffix, sep string) string {
	if len(s) == 0 {
		return ""
	}

	if len(s) == 1 {
		return s[0] + suffix
	}

	n := len(sep) * (len(s) - 1)
	for i, sl := 0, len(suffix); i < len(s); i++ {
		n += len(s[i]) + sl
	}

	b := make([]byte, n)
	bp := copy(b, s[0])
	bp += copy(b[bp:], suffix)
	for _, s := range s[1:] {
		bp += copy(b[bp:], sep)
		bp += copy(b[bp:], s)
		bp += copy(b[bp:], suffix)
	}

	return string(b)
}

// JoinInt join int slice as string
func JoinInt(v []int, sep string) string {
	if len(v) == 0 {
		return ""
	}

	buf := bytes.NewBuffer([]byte{})
	buf.WriteString(strconv.Itoa(v[0]))
	for _, s := range v[1:] {
		buf.WriteString(sep)
		buf.WriteString(strconv.Itoa(s))
	}

	return buf.String()
}

// JoinInt join int slice as string
func JoinUint(v []uint, sep string) string {
	if len(v) == 0 {
		return ""
	}

	buf := bytes.NewBuffer([]byte{})
	buf.WriteString(strconv.FormatUint(uint64(v[0]), 10))
	for _, s := range v[1:] {
		buf.WriteString(sep)
		buf.WriteString(strconv.FormatUint(uint64(s), 10))
	}

	return buf.String()
}

// Compare compare two string, if equal, 0 was returned, if s1 > s2, 1 was returned,
// otherwise -1 was returned
func Compare(s1, s2 string) int {
	l1, l2 := len(s1), len(s2)
	for i := 0; i < l1 && i < l2; i++ {
		if s1[i] < s2[i] {
			return -1
		} else if s1[i] > s2[i] {
			return 1
		}
	}

	switch {
	case l1 < l2:
		return -1
	case l1 == l2:
		return 0
	default:
		return 1
	}
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

// MergeSpace merge mutiple space to one, trim determin whether remove space at prefix and suffix
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

// WriteStringsToBuffer write strings to buffer, it avoid memory allocaton of join
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

func IsEmpty(s string) bool {
	return s == ""
}

func IsNotEmpty(s string) bool {
	return s != ""
}

// NumMatched return number of strings matched by matcher
func NumMatched(matcher func(string) bool, strings ...string) int {
	var m int
	for i := 0; i < len(strings); i++ {
		if matcher(strings[i]) {
			m++
		}
	}

	return m
}

// Filter strings matched by matcher
func Filter(matcher func(string) bool, strings ...string) []string {
	filtered := make([]string, 0, len(strings))

	for i, l := 0, len(strings); i < l; i++ {
		if s := strings[i]; matcher(s) {
			filtered = append(filtered, s)
		}
	}

	return filtered
}

// Map convert string using the mapper
func Map(mapper func(string) string, strings ...string) []string {
	for i, l := 0, len(strings); i < l; i++ {
		strings[i] = mapper(strings[i])
	}

	return strings
}

func FilterInPlace(matcher func(string) bool, strings ...string) []string {
	pos := 0
	for i, l := 0, len(strings); i < l; i++ {
		if s := strings[i]; matcher(s) {
			strings[pos] = s
			pos++
		}
	}

	return strings[:pos]
}

func MultipleLineOperate(s, delim string, operate func(line, delim string) string) string {
	lines := strings.Split(s, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		lines[i] = operate(lines[i], delim)
	}

	return strings.Join(lines, "\n")
}

// MakeSlice create a string slice with given size
func MakeSlice(element string, size int) []string {
	slice := make([]string, size)
	for i := 0; i < size; i++ {
		slice[i] = element
	}
	return slice
}
