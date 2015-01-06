package types

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

// MidIn return mid of three, min must smaller than max
func MidIn(min, max, val int) int {
	if val >= max {
		return max
	} else if val >= min {
		return val
	} else {
		return min
	}
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
