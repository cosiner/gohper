package net2

import (
	"net"
	"sync/atomic"

	"github.com/cosiner/gohper/time2"
)

type Sleeper struct {
	maxSleepMs int
	minSleepMs int
	curr       int
}

func NewSleeper(minMs, maxMs int) Sleeper {
	if minMs <= 0 {
		minMs = 5
	}
	if maxMs <= 0 {
		maxMs = 1000
	}
	if minMs > maxMs {
		minMs = maxMs
	}
	return Sleeper{
		maxSleepMs: maxMs,
		minSleepMs: minMs,
	}
}

func (s *Sleeper) Sleep() {
	s.curr = time2.LimitSleep(s.curr, s.minSleepMs, s.maxSleepMs)
}

func (s *Sleeper) Reset() {
	s.curr = 0
}

type retryListener struct {
	sleeper Sleeper
	net.Listener
}

func NewRetryListener(l net.Listener, minSleepMs, maxSleepMs int) net.Listener {
	return &retryListener{
		sleeper:  NewSleeper(minSleepMs, maxSleepMs),
		Listener: l,
	}
}

func RetryListen(netname, addr string, minSleepMs, maxSleepMs int) (net.Listener, error) {
	l, err := net.Listen(netname, addr)
	if err != nil {
		return nil, err
	}
	return NewRetryListener(l, minSleepMs, maxSleepMs), nil
}

func (r *retryListener) Accept() (net.Conn, error) {
	for {
		c, e := r.Listener.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				r.sleeper.Sleep()
				continue
			}
		}
		r.sleeper.Reset()
		return c, e
	}
}

type limitConn struct {
	l *limitListener
	net.Conn
}

func (c limitConn) Close() error {
	c.l.incrCurr(-1)
	return c.Conn.Close()
}

type limitListener struct {
	sleeper Sleeper

	max  int64
	curr int64
	net.Listener
}

func NewLimitListener(l net.Listener, maxConn int64, minSleepMs, maxSleepMs int) net.Listener {
	if maxConn < 0 {
		maxConn = 10240
	}
	return &limitListener{
		sleeper:  NewSleeper(minSleepMs, maxSleepMs),
		max:      maxConn,
		Listener: l,
	}
}

func LimitListen(netname, addr string, maxConn int64, minSleepMs, maxSleepMs int) (net.Listener, error) {
	l, err := net.Listen(netname, addr)
	if err != nil {
		return nil, err
	}
	return NewLimitListener(l, maxConn, minSleepMs, maxSleepMs), nil
}

func (l *limitListener) incrCurr(n int64) {
	atomic.AddInt64(&l.curr, n)
}

func (l *limitListener) Accept() (net.Conn, error) {
	for atomic.LoadInt64(&l.curr) >= l.max {
		l.sleeper.Sleep()
	}
	l.sleeper.Reset()
	c, e := l.Listener.Accept()
	if e != nil {
		return nil, e
	}
	l.incrCurr(1)
	return limitConn{
		Conn: c,
		l:    l,
	}, nil
}
