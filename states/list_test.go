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

	list.PushBack(SECOND)
	list.PushBack(THIRD)
	list.PushBack(FOURTH)
	list.PushBack(FIRST)

	tt.
		Eq(FIRST, list.PopBack()).
		Eq(FOURTH, list.PopBack()).
		Eq(THIRD, list.PopBack()).
		Eq(SECOND, list.PopBack())

	// stack: front in, front out
	list.PushFront(SECOND)
	list.PushFront(THIRD)
	list.PushFront(FOURTH)
	list.PushFront(FIRST)

	tt.
		Eq(FIRST, list.PopFront()).
		Eq(FOURTH, list.PopFront()).
		Eq(THIRD, list.PopFront()).
		Eq(SECOND, list.PopFront())

	// queue: front in, back out
	list.PushFront(SECOND)
	list.PushFront(THIRD)
	list.PushFront(FOURTH)
	list.PushFront(FIRST)

	tt.
		Eq(SECOND, list.PopBack()).
		Eq(THIRD, list.PopBack()).
		Eq(FOURTH, list.PopBack()).
		Eq(FIRST, list.PopBack())

	// queue: back in, front out
	list.PushBack(SECOND)
	list.PushBack(THIRD)
	list.PushBack(FOURTH)
	list.PushBack(FIRST)

	tt.
		Eq(SECOND, list.PopFront()).
		Eq(THIRD, list.PopFront()).
		Eq(FOURTH, list.PopFront()).
		Eq(FIRST, list.PopFront())
}
