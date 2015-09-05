package time2

import "time"

type Timer struct {
	C      <-chan time.Time
	timer  *time.Timer
	ticker *time.Ticker
	tick   time.Duration
}

func NewTimer(first time.Time, tick time.Duration) *Timer {
	t := &Timer{
		tick:  tick,
		timer: time.NewTimer(first.Sub(time.Now())),
	}
	t.C = t.timer.C
	return t
}

func (t *Timer) Wait() (time.Time, bool) {
	var now time.Time
	var ok bool

	if t.ticker == nil {
		now, ok = <-t.timer.C
	} else {
		now, ok = <-t.ticker.C
	}
	if ok {
		t.Switch()
	}

	return now, ok
}

func (t *Timer) Switch() {
	if t.ticker != nil {
		return
	}

	t.timer.Stop()
	t.timer = nil
	t.ticker = time.NewTicker(t.tick)
	t.C = t.ticker.C
}

func (t *Timer) Stop() {
	if t.ticker == nil {
		t.timer.Stop()
	} else {
		t.ticker.Stop()
	}
}
