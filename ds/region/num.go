package region

// Direction means positive or reverse
type Direction bool

const (
	POSITIVE Direction = true
	REVERSE  Direction = false
)

// SeqByDir return the sequence of a and b sorted by direction
// if it's positive, a, b is returned, otherwise, b, a is returned
func SeqByDir(a, b int, dir Direction) (int, int) {
	if dir == POSITIVE {
		return a, b
	}

	return b, a
}

// MinByDir return minimum of a and b sorted by direction
func MinByDir(a, b int, dir Direction) int {
	if dir == POSITIVE {
		return a
	}

	return b
}

// MaxByDir return maxium of a and b sorted by direction
func MaxByDir(a, b int, dir Direction) int {
	if dir == POSITIVE {
		return b
	}

	return a
}

// Seq return sequenced a and b
func Seq(a, b int) (int, int) {
	if a <= b {
		return a, b
	}

	return b, a
}

// Min return smaller in a and b
func Min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

// Max return bigger in a and b
func Max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// Mid return mid of three
func Mid(a, b, c int) int {
	if a >= b {
		if b >= c {
			return b
		} else if a >= c {
			return c
		} else {
			return a
		}
	} else {
		if c <= a {
			return a
		} else if c <= b {
			return c
		} else {
			return b
		}
	}
}

// Abs return absolute a
func Abs(a int) int {
	if a < 0 {
		return -a
	}

	return a
}

// Pow return power of a base number
func Pow(base, power uint) uint64 {
	n := uint64(1)
	ubase := uint64(base)
	for power > 0 {
		n *= ubase
		power--
	}

	return uint64(n)
}
