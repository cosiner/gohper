package time2

import (
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	now := time.Now()
	timer := NewTimer(now.Add(time.Millisecond), time.Millisecond)
	go func() {
		// for {
		// 	select {
		// 	case tnow, ok := <-timer.C:
		// 		if !ok {
		// 			return
		// 		}
		// 		timer.Switch()
		// 		tt.Eq(time.Millisecond, tnow.Sub(now))
		// 	}
		// }
		for {
			timer.Wait()
		}
	}()
	time.Sleep(time.Millisecond * 100)
}
