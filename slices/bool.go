package slices

type Bools []bool

func MakeBools(i bool, n int) Bools {
	v := make(Bools, n)
	if !i {
		return v
	}

	for n := n - 1; n >= 0; n-- {
		v[n] = i
	}
	return v
}

func (v Bools) Bools() []bool {
	return []bool(v)
}
