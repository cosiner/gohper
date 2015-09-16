package time2

import (
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	now := time.Now()
	ticker := NewTicker(now.Add(time.Millisecond), time.Millisecond)
	go func() {
		// for {
		// 	select {
		// 	case tnow, ok := <-ticker.C:
		// 		if !ok {
		// 			return
		// 		}
		// 		ticker.Switch()
		// 		tt.Eq(time.Millisecond, tnow.Sub(now))
		// 	}
		// }
		for {
			ticker.Wait()
		}
	}()
	time.Sleep(time.Millisecond * 100)
}
