package time2

import (
	"testing"
	"time"

	"github.com/cosiner/gohper/testing2"
)

func TestParseHuman(t *testing.T) {
	testing2.
		Expect(time.Hour, nil).Arg("1H").
		Expect(time.Minute, nil).Arg("1M").
		Expect(time.Second, nil).Arg("1S").
		Expect(time.Millisecond, nil).Arg("1m").
		Expect(time.Microsecond, nil).Arg("1u").
		Expect(time.Nanosecond, nil).Arg("1n").
		Expect(time.Duration(0), testing2.NonNil).Arg("1z").
		Run(t, ParseHuman)
}

func TestMonthDays(t *testing.T) {
	tt := testing2.Wrap(t)
	type YearMonth struct {
		Year, Month int
		Days        int
		IsLeap      bool
	}

	tests := []YearMonth{
		YearMonth{2000, 3, 31, true},
		YearMonth{2000, 4, 30, true},
		YearMonth{2000, 2, 29, true},
		YearMonth{2001, 2, 28, false},
		YearMonth{2300, 2, 28, false},
		YearMonth{2400, 2, 29, true},
	}

	for _, t := range tests {
		tt.Eq(t.Days, MonthDays(t.Year, t.Month))
		tt.Eq(t.IsLeap, IsLeapYear(t.Year))
	}

	defer tt.Recover()
	MonthDays(2014, 13)
}

func TestTiming(t *testing.T) {
	tt := testing2.Wrap(t)
	timing := Timing()

	tt.Log(timing())
	tt.Log(timing())
	tt.Log(timing())
	tt.Log(timing())
	tt.Log(timing())
}
