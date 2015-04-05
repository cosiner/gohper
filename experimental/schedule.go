package experiament

import (
	"sync"

	"github.com/cosiner/gohper/lib/signal"
)

// WARNNING: UNTESTED

// Task Schedule By scale, each user has his scale
type Task interface {
	Execute()
}

type TaskFunc func()

func (f TaskFunc) Execute() {
	f()
}

type Queue struct {
	sync.Mutex
	tasks []Task
}

func (q *Queue) Enqueue(t Task) bool {
	q.Lock()
	defer q.Unlock()
	if len(q.tasks) == cap(q.tasks) {
		return false
	}
	q.tasks = append(q.tasks, t)
	return true
}

func (q *Queue) Dequeue() Task {
	q.Lock()
	defer q.Unlock()
	if l := len(q.tasks); l == 0 {
		return nil
	} else {
		t := q.tasks[l-1]
		q.tasks = q.tasks[:l-1]
		return t
	}
}

type Scheduler struct {
	cond   *signal.Signal
	queues map[uint32]*Queue
}

func New() *Scheduler {
	return &Scheduler{
		cond:   signal.New(),
		queues: make(map[uint32]*Queue),
	}
}

func (s *Scheduler) AddTask(id uint32, t Task) (success bool) {
	if q := s.queues[id]; q != nil {
		if success = q.Enqueue(t); !success {
			for _, q = range s.queues {
				if q.Enqueue(t) {
					success = true
					s.cond.Notify()
					break
				}
			}
		} else {
			s.cond.Notify()
		}
	}
	return
}

func (s *Scheduler) GetTask() Task {
	for _, q := range s.queues {
		if t := q.Dequeue(); t != nil {
			return t
		}
	}
	return nil
}

func (s *Scheduler) AddQueue(id uint32, capacity int) {
	s.queues[id] = &Queue{tasks: make([]Task, 0, capacity)}
}

func (s *Scheduler) Run() {
	go func() {
		for {
			if t := s.GetTask(); t != nil {
				t.Execute()
			} else {
				s.cond.Wait()
			}
		}
	}()
}
