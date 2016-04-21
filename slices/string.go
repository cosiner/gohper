package slices

import (
	"math/rand"
	"sort"
	"strings"
)

func EqualString(i string) func(string) bool {
	return func(u string) bool {
		return u == i
	}
}

type Strings []string

func NewStrings(s ...string) Strings {
	return Strings(s)
}

func MakeStrings(i string, n int) Strings {
	v := make(Strings, n)
	for n = n - 1; n >= 0; n-- {
		v[n] = i
	}

	return v
}

func (v Strings) Strings() []string {
	return []string(v)
}

func (v Strings) Len() int {
	return len(v)
}

func (v Strings) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v Strings) Less(i, j int) bool {
	return v[i] < v[j]
}

func (v Strings) SafeGet(index int) string {
	if isSafeIndex(index, len(v)) {
		return v[index]
	}
	return ""
}

func (v Strings) SafeSet(index int, val string) bool {
	if isSafeIndex(index, len(v)) {
		v[index] = val
		return true
	}
	return false
}

func (v Strings) FitCapToLen() Strings {
	if l := len(v); l != cap(v) {
		new := make(Strings, l)
		copy(new, v)
		return new
	}

	return v
}

func (v Strings) IncrAppend(i string) Strings {
	l := len(v)
	if l < cap(v) {
		return append(v, i)
	}

	new := make(Strings, l+1)
	copy(new, v)
	new[l] = i
	return new
}

func (v Strings) Map(mapper func(string) string) Strings {
	for i, l := 0, len(v); i < l; i++ {
		v[i] = mapper(v[i])
	}
	return v
}

func (v Strings) Filter(filter func(string) bool) Strings {
	var new Strings
	for i, l := 0, len(v); i < l; i++ {
		if e := v[i]; filter(e) {
			new = append(new, e)
		}
	}

	return new
}

func (v Strings) Find(filter func(string) bool) int {
	for i, l := 0, len(v); i < l; i++ {
		if filter(v[i]) {
			return i
		}
	}
	return -1
}

func (v Strings) NumMatched(filter func(string) bool) int {
	var n int
	for i, l := 0, len(v); i < l; i++ {
		if filter(v[i]) {
			n++
		}
	}
	return n
}

func (v Strings) FilterInplace(filter func(string) bool) Strings {
	var prev = -1

	for i, l := 0, len(v); i < l; i++ {
		if e := v[i]; filter(e) {
			prev++
			v[prev] = e
		}
	}

	return v[:prev+1]
}

func (v Strings) RmDups() Strings {
	l := len(v)
	if l == 0 {
		return v
	}

	sort.Sort(v)
	prev := 0
	for i := 1; i < l; i++ {
		if s := v[i]; s != v[prev] {
			prev++
			v[prev] = s
		}
	}

	return v[:prev+1]
}

func (v Strings) Remove(index int) Strings {
	l := len(v)
	if index < 0 || index >= l {
		return v
	}

	if l >= 2*(index+1) {
		for ; index > 0; index-- {
			v[index] = v[index-1]
		}
		v = v[1:]
	} else {
		for ; index < l-1; index++ {
			v[index] = v[index+1]
		}
		v = v[:l-1]
	}

	return v
}

func (v Strings) Append(i string) Strings {
	return append(v, i)
}

func (v Strings) Rand() string {
	l := len(v)
	if l == 0 {
		return ""
	}

	return v[rand.Intn(l)]
}

func (v Strings) Clear(i string) Strings {
	return v.FilterInplace(func(vi string) bool {
		return vi != i
	})
}

func (v Strings) Replace(old, new string) Strings {
	return v.Map(func(vi string) string {
		if vi == old {
			return new
		}
		return old
	})
}

func (v Strings) Join(suffix, sep string) string {
	if suffix == "" {
		return strings.Join(v, sep)
	}

	l := len(v)
	if l == 0 {
		return ""
	}

	buf := make([]byte, 0, l*2+len(suffix)*l+len(sep)*(l-1))
	for i := 0; i < l; i++ {
		buf = append(buf, v[i]...)
		buf = append(buf, suffix...)
		if i != l-1 {
			buf = append(buf, sep...)
		}
	}
	return string(buf)
}

func (v Strings) Contains(s string) bool {
	return v.Find(EqualString(s)) >= 0
}

func (v Strings) ToInterfaces() []interface{} {
	ifs := make([]interface{}, len(v))
	for i, s := range v {
		ifs[i] = s
	}
	return ifs
}

func (v Strings) IsSame(i, j int) bool {
	return v[i] == v[j]
}

func (v Strings) Merge(dst, src int) {}

func (v Strings) Move(dst, src int) {
	v[dst] = v[src]
}
