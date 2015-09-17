package time2

import "time"

var TimingEnable = true

func Timing() func(...func()) time.Duration {
	if !TimingEnable {
		return func(...func()) time.Duration {
			return 0
		}
	}

	var last time.Time
	f := func(fn ...func()) time.Duration {
		for _, f := range fn {
			if f != nil {
				f()
			}
		}

		now := Now()
		d := now.Sub(last)
		last = now
		return d
	}
	last = Now()
	return f
}
