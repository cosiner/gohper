package pager

import (
	"strconv"
	"sync"
)

// Pager is a tool to help paging
type Pager struct {
	BeginPage  int
	BeginIndex int
	PageSize   int

	MaxPage int
}

func (p *Pager) Begin(page int) int {
	if page <= p.BeginPage {
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

func (p *Pager) IsOverRange(start, count int) bool {
	if p.MaxPage > 0 {
		return (start + count) > p.PageSize*p.MaxPage
	}
	return false
}

func (p *Pager) IsReachBottom(start, count, maxPage int) bool {
	if maxPage <= 0 {
		maxPage = p.MaxPage
	} else if p.MaxPage > 0 && p.MaxPage < maxPage {
		maxPage = p.MaxPage
	}
	if maxPage > 0 {
		return (start + count) >= p.PageSize*maxPage
	}
	return false
}

func (p *Pager) EndByString(page string) int {
	return p.BeginByString(page) + p.PageSize
}

type PagerGroup struct {
	pagers []Pager
	lock   sync.Mutex
}

func (pg *PagerGroup) Add(beginPage, beginIndex, pageSize, maxPage int) *Pager {
	if beginPage < 0 {
		beginPage = 1
	}

	if beginIndex < 0 {
		beginIndex = 0
	}

	pg.lock.Lock()
	l := len(pg.pagers)
	pg.pagers = append(pg.pagers, Pager{
		BeginPage:  beginPage,
		BeginIndex: beginIndex,
		PageSize:   pageSize,
		MaxPage:    maxPage,
	})
	p := &pg.pagers[l]
	pg.lock.Unlock()
	return p
}
