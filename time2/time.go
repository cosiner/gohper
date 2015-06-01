package time2

import (
	"strconv"
	"strings"
	"time"
)

const DATETIME_FMT = "2006/01/02 15:04:05"
const DATE_FMT = "2006/01/02"
const TIME_FMT = "15:04:05"
const FMT_USAGE = `yyyy - 2006
yy   - 06
mm   - 01
dd   - 02
HH   - 15
MM   - 04
SS   - 05`

var formatReplacer = strings.NewReplacer([]string{
	"yyyy", "2006",
	"yy", "06",
	"mm", "01",
	"dd", "02",
	"HH", "15",
	"MM", "04",
	"SS", "05",
}...)

func UnixNanoSinceNow(sec int) int64 {
	return time.Now().Add(time.Duration(sec) * time.Second).UnixNano()
}

// FormatLayout convert human readable time format to golang time format
// yyyy : year, yy:year, mm:month, dd:day, HH:hour, MM:minute, SS:second
func FormatLayout(format string) string {
	return formatReplacer.Replace(format)
}

// NowTimeUnix is a wrapper of time.Now().Unix()
func NowTimeUnix() uint64 {
	return uint64(time.Now().Unix())
}

// NowTimeUnixNano is a wrapper of time.Now().UnixNano()
func NowTimeUnixNano() uint64 {
	return uint64(time.Now().UnixNano())
}

// NowTime return time string in gived format
func NowTime(format string) string {
	return time.Now().Format(FormatLayout(format))
}

// DateTime return curremt datetime in format yyyy/mm/dd HH:MM:SS
func DateTime() string {
	return time.Now().Format(DATETIME_FMT)
}

// Time return current time in format HH:MM:SS
func Time() string {
	return time.Now().Format(TIME_FMT)
}

// Date return current date in format: yyyy/mm/dd
func Date() string {
	return time.Now().Format(DATE_FMT)
}

// ParseTime parse gived time string in format: yyyy/mm/dd HH:MM:SS
func ParseTime(t string) (time.Time, error) {
	return time.Parse(DATETIME_FMT, t)
}

// Timing the cost of function call, unix nano was returned
func Timing(f func()) int64 {
	now := time.Now().UnixNano()
	f()

	return time.Now().UnixNano() - now
}

// ToHuman convert nano to human time size, insufficient portion will be discarded
// performs rounding.
//
// support 0-999ns, 0-999us, 0-999ms, 0-Maxs,
func ToHuman(nano int64) string {
	var base int64 = 1
	if nano < 1000*base {
		return strconv.Itoa(int(nano/base)) + "ns"
	}

	base *= 1000
	if nano < 1000*base {
		var us = int(nano / base)
		if nano%base >= base/2 {
			us++
		}

		return strconv.Itoa(us) + "us"
	}

	base *= 1000
	if nano < 1000*base {
		var ms = int(nano / base)
		if nano%base >= base/2 {
			ms++
		}
		return strconv.Itoa(ms) + "ms"
	}

	base *= 1000
	var s = int(nano / base)
	if nano%base >= base/2 {
		s++
	}
	return strconv.Itoa(s) + "s"
}
