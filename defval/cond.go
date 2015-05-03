package defval

type Cond bool

func (c Cond) String(a, b string) string {
	if c {
		return a
	}
	return b
}

func (c Cond) Int(a, b int) int {
	if c {
		return a
	}
	return b
}

func (c Cond) Int8(a, b int8) int8 {
	if c {
		return a
	}
	return b
}

func (c Cond) Int16(a, b int16) int16 {
	if c {
		return a
	}
	return b
}

func (c Cond) Int32(a, b int32) int32 {
	if c {
		return a
	}
	return b
}

func (c Cond) Int64(a, b int64) int64 {
	if c {
		return a
	}
	return b
}

func (c Cond) Uint(a, b uint) uint {
	if c {
		return a
	}
	return b
}

func (c Cond) Uint8(a, b uint8) uint8 {
	if c {
		return a
	}
	return b
}

func (c Cond) Uint16(a, b uint16) uint16 {
	if c {
		return a
	}
	return b
}

func (c Cond) Uint32(a, b uint32) uint32 {
	if c {
		return a
	}
	return b
}

func (c Cond) Uint64(a, b uint64) uint64 {
	if c {
		return a
	}
	return b
}
