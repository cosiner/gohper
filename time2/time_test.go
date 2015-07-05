package time2

import (
	"testing"
	"time"

	"github.com/cosiner/gohper/testing2"
)

func TestToHuman(t *testing.T) {
	tt := testing2.Wrap(t)
	type Test struct {
		Time  int64
		Human string
	}

	tests := []Test{
		Test{0, "0ns"},
		Test{999, "999ns"},

		Test{1000, "1us"},
		Test{1499, "1us"},
		Test{1500, "2us"},

		Test{1000 * 1000, "1ms"},
		Test{1000 * 1499, "1ms"},
		Test{1000 * 1500, "2ms"},

		Test{1000 * 1000 * 1000, "1s"},
		Test{1000 * 1000 * 1499, "1s"},
		Test{1000 * 1000 * 1500, "2s"},

		Test{1000 * 1000 * 1000 * 10000, "10000s"},
	}

	for _, test := range tests {
		tt.Eq(test.Human, ToHuman(test.Time))
	}
}

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
