package states

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

const (
	FIRST uint = iota
	SECOND
	THIRD
	FOURTH
)

func TestUnitSize(t *testing.T) {
	testing2.
		Expect(uint(1)).Arg(uint(0)).
		Expect(uint(1)).Arg(uint(1)).
		Expect(uint(1)).Arg(uint(2)).
		Expect(uint(2)).Arg(uint(3)).
		Expect(uint(2)).Arg(uint(4)).
		Expect(uint(3)).Arg(uint(5)).
		Expect(uint(3)).Arg(uint(6)).
		Expect(uint(3)).Arg(uint(7)).
		Expect(uint(3)).Arg(uint(8)).
		Expect(uint(4)).Arg(uint(9)).
		Run(t, UnitSize)
}

func TestList(t *testing.T) {
	tt := testing2.Wrap(t)

	list := NewList(3)

	tt.True(21 == list.MaxSize())
	for i := 0; i < 21; i++ {
		list.PushBack(FIRST)
	}
	tt.True(list.IsFull())

	for i := 0; i < 21; i++ {
		tt.Eq(FIRST, list.PopBack())
	}
	tt.True(list.IsEmpty())

	// stack: back in, back out

	list.
		PushBack(SECOND).
		PushBack(THIRD).
		PushBack(FOURTH).
		PushBack(FIRST)

	tt.
		Eq(FIRST, list.PopBack()).
		Eq(FOURTH, list.PopBack()).
		Eq(THIRD, list.PopBack()).
		Eq(SECOND, list.PopBack())

	// stack: front in, front out
	list.
		PushFront(SECOND).
		PushFront(THIRD).
		PushFront(FOURTH).
		PushFront(FIRST)

	tt.
		Eq(FIRST, list.PopFront()).
		Eq(FOURTH, list.PopFront()).
		Eq(THIRD, list.PopFront()).
		Eq(SECOND, list.PopFront())

	// queue: front in, back out
	list.
		PushFront(SECOND).
		PushFront(THIRD).
		PushFront(FOURTH).
		PushFront(FIRST)

	tt.
		Eq(SECOND, list.PopBack()).
		Eq(THIRD, list.PopBack()).
		Eq(FOURTH, list.PopBack()).
		Eq(FIRST, list.PopBack())

	// queue: back in, front out
	list.
		PushBack(SECOND).
		PushBack(THIRD).
		PushBack(FOURTH).
		PushBack(FIRST)

	tt.
		Eq(SECOND, list.PopFront()).
		Eq(THIRD, list.PopFront()).
		Eq(FOURTH, list.PopFront()).
		Eq(FIRST, list.PopFront())
}
