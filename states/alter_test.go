package states

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestStack(t *testing.T) {
	tt := testing2.Wrap(t)
	stack := NewStack(3)

	stack.Push(SECOND)
	stack.Push(THIRD)
	stack.Push(FOURTH)
	stack.Push(FIRST)

	tt.
		Eq(FIRST, stack.Pop()).
		Eq(FOURTH, stack.Pop()).
		Eq(THIRD, stack.Pop()).
		Eq(SECOND, stack.Pop())
}

func TestQueue(t *testing.T) {
	tt := testing2.Wrap(t)
	queue := NewQueue(3)

	queue.Push(SECOND)
	queue.Push(THIRD)
	queue.Push(FOURTH)
	queue.Push(FIRST)

	tt.
		Eq(SECOND, queue.Pop()).
		Eq(THIRD, queue.Pop()).
		Eq(FOURTH, queue.Pop()).
		Eq(FIRST, queue.Pop())
}
