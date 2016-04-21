package slices

import (
	"math/rand"
	"sort"
	"strconv"
)

func EqualUint(i uint) func(uint) bool {
	return func(u uint) bool {
		return u == i
	}
}

type Uints []uint

func NewUints(u ...uint) Uints {
	return Uints(u)
}

func MakeUints(i uint, n int) Uints {
	v := make(Uints, n)
	for n = n - 1; n >= 0; n-- {
		v[n] = i
	}

	return v
}

func (v Uints) Uints() []uint {
	return []uint(v)
}

func (v Uints) Len() int {
	return len(v)
}

func (v Uints) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v Uints) Less(i, j int) bool {
	return v[i] < v[j]
}

func (v Uints) SafeGet(index int) uint {
	if isSafeIndex(index, len(v)) {
		return v[index]
	}
	return 0
}

func (v Uints) SafeSet(index int, val uint) bool {
	if isSafeIndex(index, len(v)) {
		v[index] = val
		return true
	}
	return false
}

func (v Uints) FitCapToLen() Uints {
	if l := len(v); l != cap(v) {
		new := make(Uints, l)
		copy(new, v)
		return new
	}

	return v
}

func (v Uints) IncrAppend(i uint) Uints {
	l := len(v)
	if l < cap(v) {
		return append(v, i)
	}

	new := make(Uints, l+1)
	copy(new, v)
	new[l] = i
	return new
}

func (v Uints) Map(mapper func(uint) uint) Uints {
	for i, l := 0, len(v); i < l; i++ {
		v[i] = mapper(v[i])
	}
	return v
}

func (v Uints) Filter(filter func(uint) bool) Uints {
	var new Uints
	for i, l := 0, len(v); i < l; i++ {
		if e := v[i]; filter(e) {
			new = append(new, e)
		}
	}

	return new
}

func (v Uints) Find(filter func(uint) bool) int {
	for i, l := 0, len(v); i < l; i++ {
		if filter(v[i]) {
			return i
		}
	}
	return -1
}

func (v Uints) NumMatched(filter func(uint) bool) int {
	var n int
	for i, l := 0, len(v); i < l; i++ {
		if filter(v[i]) {
			n++
		}
	}
	return n
}

func (v Uints) FilterInplace(filter func(uint) bool) Uints {
	var prev = -1

	for i, l := 0, len(v); i < l; i++ {
		if e := v[i]; filter(e) {
			prev++
			v[prev] = e
		}
	}

	return v[:prev+1]
}

func (v Uints) RmDups() Uints {
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

func (v Uints) Remove(index int) Uints {
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

func (v Uints) Append(i uint) Uints {
	return append(v, i)
}

func (v Uints) Rand() uint {
	l := len(v)
	if l == 0 {
		return 0
	}

	return v[rand.Intn(l)]
}

func (v Uints) Clear(i uint) Uints {
	return v.FilterInplace(func(vi uint) bool { return vi != i })
}

func (v Uints) Replace(old, new uint) Uints {
	return v.Map(func(vi uint) uint {
		if vi == old {
			return new
		}
		return old
	})
}

func (v Uints) Join(suffix, sep string) string {
	l := len(v)
	if l == 0 {
		return ""
	}

	buf := make([]byte, 0, l*2+len(suffix)*l+len(sep)*(l-1))
	for i := 0; i < l; i++ {
		buf = strconv.AppendUint(buf, uint64(v[i]), 10)
		buf = append(buf, suffix...)
		if i != l-1 {
			buf = append(buf, sep...)
		}
	}
	return string(buf)
}

func (v Uints) Contains(u uint) bool {
	return v.Find(EqualUint(u)) >= 0
}

func (v Uints) ToInterfaces() []interface{} {
	ifs := make([]interface{}, len(v))
	for i, s := range v {
		ifs[i] = s
	}
	return ifs
}

func (v Uints) IsSame(i, j int) bool {
	return v[i] == v[j]
}

func (v Uints) Merge(dst, src int) {}

func (v Uints) Move(dst, src int) {
	v[dst] = v[src]
}
