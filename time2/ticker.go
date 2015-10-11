package time2

import "time"

type TimeTicker struct {
	timer  *time.Timer
	ticker *time.Ticker

	tm time.Time
	tk time.Duration
}

func NewTimeTicker(first time.Time, tick time.Duration) *TimeTicker {
	now := Now()
	sub := first.Sub(now)
	for sub < 0 {
		sub += tick
	}

	t := &TimeTicker{
		timer: time.NewTimer(sub),

		tm: now.Add(sub),
		tk: tick,
	}
	return t
}

func (t *TimeTicker) C() <-chan time.Time {
	if Now().Before(t.tm) || len(t.timer.C) > 0 {
		return t.timer.C
	}
	if t.ticker == nil {
		t.ticker = time.NewTicker(t.tk)
	}
	return t.ticker.C
}

func (t *TimeTicker) Stop() {
	if t.ticker == nil {
		t.timer.Stop()
	} else {
		t.ticker.Stop()
	}
}
