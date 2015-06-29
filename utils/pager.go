package utils

import "strconv"

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
