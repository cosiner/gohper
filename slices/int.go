package slices

import (
	"math/rand"
	"sort"
	"strconv"
)

func EqualInt(i int) func(int) bool {
	return func(u int) bool {
		return u == i
	}
}

type Ints []int

func NewInts(i ...int) Ints {
	return Ints(i)
}

func MakeInts(i int, n int) Ints {
	v := make(Ints, n)
	for n = n - 1; n >= 0; n-- {
		v[n] = i
	}

	return v
}

func (v Ints) Ints() []int {
	return []int(v)
}

func (v Ints) Len() int {
	return len(v)
}

func (v Ints) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v Ints) Less(i, j int) bool {
	return v[i] < v[j]
}

func (v Ints) SafeGet(index int) int {
	if isSafeIndex(index, len(v)) {
		return v[index]
	}
	return 0
}

func (v Ints) SafeSet(index int, val int) bool {
	if isSafeIndex(index, len(v)) {
		v[index] = val
		return true
	}
	return false
}

func (v Ints) FitCapToLen() Ints {
	if l := len(v); l != cap(v) {
		new := make(Ints, l)
		copy(new, v)
		return new
	}

	return v
}

func (v Ints) IncrAppend(i int) Ints {
	l := len(v)
	if l < cap(v) {
		return append(v, i)
	}

	new := make(Ints, l+1)
	copy(new, v)
	new[l] = i
	return new
}

func (v Ints) Map(mapper func(int) int) Ints {
	for i, l := 0, len(v); i < l; i++ {
		v[i] = mapper(v[i])
	}
	return v
}

func (v Ints) Filter(filter func(int) bool) Ints {
	var new Ints
	for i, l := 0, len(v); i < l; i++ {
		if e := v[i]; filter(e) {
			new = append(new, e)
		}
	}

	return new
}

func (v Ints) Find(filter func(int) bool) int {
	for i, l := 0, len(v); i < l; i++ {
		if filter(v[i]) {
			return i
		}
	}
	return -1
}

func (v Ints) FilterInplace(filter func(int) bool) Ints {
	var prev = -1

	for i, l := 0, len(v); i < l; i++ {
		if e := v[i]; filter(e) {
			prev++
			v[prev] = e
		}
	}

	return v[:prev+1]
}

func (v Ints) RmDups() Ints {
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

func (v Ints) Remove(index int) Ints {
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

func (v Ints) Append(i int) Ints {
	return append(v, i)
}

func (v Ints) NumMatched(filter func(int) bool) int {
	var n int
	for i, l := 0, len(v); i < l; i++ {
		if filter(v[i]) {
			n++
		}
	}
	return n
}

func (v Ints) Rand() int {
	l := len(v)
	if l == 0 {
		return 0
	}

	return v[rand.Intn(l)]
}

func (v Ints) Clear(i int) Ints {
	return v.FilterInplace(func(vi int) bool { return vi != i })
}

func (v Ints) Replace(old, new int) Ints {
	return v.Map(func(vi int) int {
		if vi == old {
			return new
		}
		return old
	})
}

func (v Ints) Join(suffix, sep string) string {
	l := len(v)
	if l == 0 {
		return ""
	}

	buf := make([]byte, 0, l*2+len(suffix)*l+len(sep)*(l-1))
	for i := 0; i < l; i++ {
		buf = strconv.AppendInt(buf, int64(v[i]), 10)
		buf = append(buf, suffix...)
		if i != l-1 {
			buf = append(buf, sep...)
		}
	}
	return string(buf)
}

func (v Ints) Contains(i int) bool {
	return v.Find(EqualInt(i)) > 0
}

func (v Ints) ToInterfaces() []interface{} {
	ifs := make([]interface{}, len(v))
	for i, s := range v {
		ifs[i] = s
	}
	return ifs
}
