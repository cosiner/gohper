package time2

import "time"

type Ticker struct {
	C      <-chan time.Time
	timer  *time.Timer
	ticker *time.Ticker
	tick   time.Duration
}

func NewTicker(first time.Time, tick time.Duration) *Ticker {
	now := time.Now()
	sub := now.Sub(first)
	sub2 := sub / tick * tick
	if sub2 < sub {
		sub2 += tick
	}

	first = first.Add(sub2)
	t := &Ticker{
		tick:  tick,
		timer: time.NewTimer(sub2 - sub),
	}
	t.C = t.timer.C
	return t
}

func (t *Ticker) Wait() (time.Time, bool) {
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

func (t *Ticker) Switch() {
	if t.ticker != nil {
		return
	}

	t.timer.Stop()
	t.timer = nil
	t.ticker = time.NewTicker(t.tick)
	t.C = t.ticker.C
}

func (t *Ticker) Stop() {
	if t.ticker == nil {
		t.timer.Stop()
	} else {
		t.ticker.Stop()
	}
}
