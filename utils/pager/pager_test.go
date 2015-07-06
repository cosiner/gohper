package pager

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestPager(t *testing.T) {
	p := &Pager{
		BeginPage:  1,
		BeginIndex: 0,
		PageSize:   10,
	}
	testPager(p, t)
}

func TestPagerGroup(t *testing.T) {
	pg := PagerGroup{}
	testPager(pg.Add(1, 0, 10), t)
	testPager(pg.Add(-1, -1, 10), t)
	testPager(pg.Add(-1, -2, 10), t)
	testPager(pg.Add(-1, -2, 10), t)
	testPager(pg.Add(-1, -2, 10), t)
}

func testPager(p *Pager, t *testing.T) {
	testing2.
		Expect(p.BeginIndex).Arg("abcde").
		Expect(p.BeginIndex).Arg("").
		Expect(p.BeginIndex).Arg("-1").
		Expect(p.BeginIndex).Arg("0").
		Expect(p.BeginIndex).Arg("1").
		Expect(p.PageSize).Arg("2").
		Expect(p.PageSize*2).Arg("3").
		Run(t, p.BeginByString)

	testing2.
		Expect(p.PageSize).Arg("abcde").
		Expect(p.PageSize).Arg("").
		Expect(p.PageSize).Arg("-1").
		Expect(p.PageSize).Arg("0").
		Expect(p.PageSize).Arg("1").
		Expect(p.PageSize*2).Arg("2").
		Expect(p.PageSize*3).Arg("3").
		Run(t, p.EndByString)

	testing2.
		Expect(p.PageSize).Arg(-1).
		Expect(p.PageSize).Arg(0).
		Expect(p.PageSize).Arg(1).
		Expect(p.PageSize*2).Arg(2).
		Expect(p.PageSize*3).Arg(3).
		Run(t, p.End)
}
