package states

type Stack struct {
	List
}

func NewStack(unitsize uint) Stack {
	return Stack{
		List: NewList(unitsize),
	}
}

func (s *Stack) Push(state uint) *Stack {
	s.PushBack(state)
	return s
}

func (s *Stack) Pop() uint {
	return s.PopBack()
}

type Queue struct {
	List
}

func NewQueue(unitsize uint) Queue {
	return Queue{
		List: NewList(unitsize),
	}
}

func (q *Queue) Push(state uint) *Queue {
	q.PushFront(state)

	return q
}

func (q *Queue) Pop() uint {
	return q.PopBack()
}
