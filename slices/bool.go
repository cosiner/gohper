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

func (v Bools) IsSame(i, j int) bool {
	return v[i] == v[j]
}

func (v Bools) Merge(dst, src int) {}

func (v Bools) Move(dst, src int) {
	v[dst] = v[src]
}
