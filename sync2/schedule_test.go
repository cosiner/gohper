package sync2

import (
	"time"

	"testing"
)

func TestSchedule(t *testing.T) {
	s := New()
	s.AddQueue(1, 10)
	s.AddQueue(2, 20)
	s.AddQueue(3, 90)
	s.Run()
	var u1, u2, u3 int
	var t3 TaskFunc = func() {
	}
	for i := 0; i < 10000; i++ {
		// if s.AddTask(1, t1) {
		// 	u1++
		// }
		// if s.AddTask(2, t2) {
		// 	u2++
		// }
		if s.AddTask(3, t3) {
			u3++
		}
	}
	time.Sleep(2)
	t.Log(u1, u2, u3)
}
