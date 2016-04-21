package slices

import (
	"math/rand"
	"sort"
)

func EqualInterface(i interface{}) func(interface{}) bool {
	return func(ii interface{}) bool {
		return i == ii
	}
}

type Interfaces []interface{}

func NewInterfaces(i ...interface{}) Interfaces {
	return Interfaces(i)
}

func MakeInterfaces(i interface{}, n int) Interfaces {
	v := make(Interfaces, n)
	for n = n - 1; n >= 0; n-- {
		v[n] = i
	}

	return v
}

func (v Interfaces) Interfaces() []interface{} {
	return []interface{}(v)
}

func (v Interfaces) Len() int {
	return len(v)
}

func (v Interfaces) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v Interfaces) Less(i, j int) bool {
	return true
}

func (v Interfaces) SafeGet(index int) interface{} {
	if isSafeIndex(index, len(v)) {
		return v[index]
	}
	return nil
}

func (v Interfaces) SafeSet(index int, val interface{}) bool {
	if isSafeIndex(index, len(v)) {
		v[index] = val
		return true
	}
	return false
}

func (v Interfaces) FitCapToLen() Interfaces {
	if l := len(v); l != cap(v) {
		new := make(Interfaces, l)
		copy(new, v)
		return new
	}

	return v
}

func (v Interfaces) IncrAppend(i interface{}) Interfaces {
	l := len(v)
	if l < cap(v) {
		return append(v, i)
	}

	new := make(Interfaces, l+1)
	copy(new, v)
	new[l] = i
	return new
}

func (v Interfaces) Map(mapper func(interface{}) interface{}) Interfaces {
	for i, l := 0, len(v); i < l; i++ {
		v[i] = mapper(v[i])
	}
	return v
}

func (v Interfaces) Filter(filter func(interface{}) bool) Interfaces {
	var new Interfaces
	for i, l := 0, len(v); i < l; i++ {
		if e := v[i]; filter(e) {
			new = append(new, e)
		}
	}

	return new
}

func (v Interfaces) Find(filter func(interface{}) bool) int {
	for i, l := 0, len(v); i < l; i++ {
		if filter(v[i]) {
			return i
		}
	}
	return -1
}

func (v Interfaces) NumMatched(filter func(interface{}) bool) int {
	var n int
	for i, l := 0, len(v); i < l; i++ {
		if filter(v[i]) {
			n++
		}
	}
	return n
}

func (v Interfaces) FilterInplace(filter func(interface{}) bool) Interfaces {
	var prev = -1

	for i, l := 0, len(v); i < l; i++ {
		if e := v[i]; filter(e) {
			prev++
			v[prev] = e
		}
	}

	return v[:prev+1]
}

func (v Interfaces) RmDups() Interfaces {
	sort.Sort(v)
	prev := -1
	for i, l := 1, len(v); i < l; i++ {
		if v[i] != v[prev+1] {
			prev++
			v[prev] = v[i]
		}
	}
	return v[:prev+1]
}

func (v Interfaces) Remove(index int) Interfaces {
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

func (v Interfaces) Append(i ...interface{}) Interfaces {
	return append(v, i...)
}

func (v Interfaces) AppendStrings(ss ...string) Interfaces {
	for _, s := range ss {
		v = append(v, s)
	}
	return v
}

func (v Interfaces) Rand() interface{} {
	l := len(v)
	if l == 0 {
		return nil
	}

	return v[rand.Intn(l)]
}

func (v Interfaces) Clear(i interface{}) Interfaces {
	return v.FilterInplace(func(vi interface{}) bool { return vi != i })
}

func (v Interfaces) Replace(old, new interface{}) Interfaces {
	return v.Map(func(vi interface{}) interface{} {
		if vi == old {
			return new
		}
		return old
	})
}

func (v Interfaces) Contains(i interface{}) bool {
	return v.Find(EqualInterface(i)) >= 0
}
