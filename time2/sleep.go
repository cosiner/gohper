package time2

import "time"

func LimitSleep(currMs, minMs, maxMs int) int {
	if currMs == 0 {
		currMs = minMs
	} else {
		currMs = 2 * currMs
	}
	if currMs > maxMs {
		currMs = maxMs
	}
	time.Sleep(time.Millisecond * time.Duration(currMs))
	return currMs
}
