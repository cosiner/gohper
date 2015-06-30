package utils

import (
	"strconv"
	"sync"
)

// Pager is a tool to help paging
type Pager struct {
	BeginPage  int
	BeginIndex int
	PageSize   int
}

func (p *Pager) Begin(page int) int {
	if page <= 0 {
		return p.BeginIndex
	}

	return (page-p.BeginPage)*p.PageSize + p.BeginIndex
}

func (p *Pager) End(page int) int {
	return p.Begin(page) + p.PageSize
}

func (p *Pager) BeginByString(page string) int {
	if page == "" {
		return p.BeginIndex
	}

	val, err := strconv.Atoi(page)
	if err != nil {
		return p.BeginIndex
	}

	return p.Begin(val)
}

func (p *Pager) EndByString(page string) int {
	return p.BeginByString(page) + p.PageSize
}

type PagerGroup struct {
	pagers []Pager
	lock   sync.Mutex
}

func (pg *PagerGroup) Add(beginPage, beginIndex, pageSize int) (p *Pager) {
	if beginPage < 0 {
		beginPage = 1
	}

	if beginIndex < 0 {
		beginIndex = 0
	}

	pg.lock.Lock()
	if l := len(pg.pagers); l < cap(pg.pagers) {
		pg.pagers = pg.pagers[:l+1]

		p = &pg.pagers[l]
		p.BeginPage = beginPage
		p.BeginIndex = beginIndex
		p.PageSize = pageSize
	} else {
		pg.pagers = append(pg.pagers, Pager{
			BeginPage:  beginPage,
			BeginIndex: beginIndex,
			PageSize:   pageSize,
		})
		p = &pg.pagers[l]
	}

	pg.lock.Unlock()
	return
}
