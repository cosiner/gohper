// Package states implements state list, stack, queue based on a uint64
package states

type List struct {
	size   uint
	unit   uint
	states uint64
}

func UnitSize(count uint) uint {
	var n = uint(1)
	for i := uint(1); i <= 64; i++ {
		n <<= 1
		if count <= n {
			return i
		}
	}

	panic("too many states to store")
}

func NewList(unitsize uint) List {
	return List{
		unit: unitsize,
	}
}

func (l *List) MaxSize() uint {
	return 64 / l.unit
}

func (l *List) Size() uint {
	return l.size
}

func (l *List) UnitSize() uint {
	return l.unit
}

func (l *List) PushBack(state uint) *List {
	l.states <<= l.unit // append to lower bits
	l.states |= uint64(state & (1<<l.unit - 1))
	l.size++

	return l
}

func (l *List) PopBack() uint {
	state := l.states & (1<<l.unit - 1) // remove from lower bits
	l.states >>= l.unit
	l.size--

	return uint(state)
}

func (l *List) PushFront(state uint) *List {
	state = state & (1<<l.unit - 1) // prepend to higher bits
	l.states |= uint64(state << (l.size * l.unit))
	l.size++

	return l
}

func (l *List) PopFront() uint {
	l.size--
	state := l.states >> (l.size * l.unit) // remove from higher bits
	l.states &= uint64((1<<(l.size*l.unit) - 1))
	return uint(state)
}

func (l *List) IsEmpty() bool {
	return l.size == 0
}

func (l *List) IsFull() bool {
	return l.size*(l.unit+1) > 64
}
