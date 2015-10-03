package time2

import (
	"testing"
	"time"
	"fmt"
	"sync"
)

func TestTimer(t *testing.T) {
	ticker := NewTimeTicker(DateDefNow(-1, -1, -1, 15, 41, 0, 0), time.Second)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for i := 0; i < 5; i++ {
			tm := <-ticker.C()
			fmt.Println(DateTime(tm))
		}
		wg.Done()
	}()
	wg.Wait()
}
