package routinepool

import (
	"sync"

	"github.com/cosiner/gohper/sync2"
)

type Job interface{}

type Pool struct {
	processor func(Job)
	maxIdle   uint64
	maxActive uint64

	lock      sync.RWMutex
	jobs      chan Job
	numIdle   uint64
	numActive uint64

	closeCond *sync2.LockCond
}

// New create a pool with fix number of goroutine, if maxActive is 0, there is no
// limit of goroutine number
func New(processor func(Job), jobBufsize int, maxIdle, maxActive uint64) *Pool {
	p := &Pool{
		processor: processor,
		jobs:      make(chan Job, jobBufsize),
		maxIdle:   maxIdle,
		maxActive: maxActive,
	}
	return p
}

// Info return current infomation about idle and activing goroutine number of the pool
func (p *Pool) Info() (numIdle, numActive uint64) {
	p.lock.RLock()
	numIdle = p.numIdle
	numActive = p.numActive
	p.lock.RUnlock()
	return
}

// Do process a job. If there is no goroutine available and goroutine number already
// reach the limitation, it will blocked untile a goroutine is free. Otherwise
// create a new goroutine. Return false only if pool already closed
func (p *Pool) Do(job Job) bool {
	p.lock.RLock()
	closeCond := p.closeCond
	numIdle := p.numIdle
	numActive := p.numActive
	p.lock.RUnlock()

	if closeCond != nil {
		return false
	}

	if numIdle == 0 && (p.maxActive == 0 || numActive < p.maxActive) {
		p.lock.Lock()
		p.numActive++
		p.numIdle++
		p.lock.Unlock()

		go p.routine()
	}

	p.jobs <- job
	return true
}

func (p *Pool) routine() {
	for job, ok := range p.jobs {
		if !ok {
			return
		}

		p.lock.Lock()
		p.numIdle--
		p.lock.Unlock()

		p.processor(job)

		p.lock.Lock()
		closeCond := p.closeCond
		jobs := len(p.jobs)

		if jobs == 0 && closeCond != nil {
			closeCond.Signal()
			p.numActive--
			p.lock.Unlock()
			return
		}

		if p.numIdle+1 > p.maxIdle {
			p.numActive--
			p.lock.Unlock()
			return
		}

		p.numIdle++
		p.lock.Unlock()
	}
}

// Close stop receive new job, and waiting for all exists jobs to be processed
func (p *Pool) Close() {
	p.lock.Lock()
	if p.closeCond != nil {
		p.lock.Unlock()
		return
	}

	p.closeCond = sync2.NewLockCond(nil)
	if len(p.jobs) != 0 {
		p.lock.Unlock()
		p.closeCond.Wait()
	} else {
		p.lock.Unlock()
	}

	close(p.jobs)
}
