package time2

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cosiner/gohper/errors"
)

const DATETIME_FMT = "2006/01/02 15:04:05"
const DATE_FMT = "2006/01/02"
const TIME_FMT = "15:04:05"

func UnixNanoSinceNow(sec int) int64 {
	return time.Now().Add(time.Duration(sec) * time.Second).UnixNano()
}

// NowTimeUnix is a wrapper of time.Now().Unix()
func NowTimeUnix() uint64 {
	return uint64(time.Now().Unix())
}

// NowTimeUnixNano is a wrapper of time.Now().UnixNano()
func NowTimeUnixNano() uint64 {
	return uint64(time.Now().UnixNano())
}

func DateAndTime() (string, string) {
	now := time.Now()
	return now.Format(DATE_FMT), now.Format(TIME_FMT)
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

const (
	ErrWrongHumanTimeFormat = errors.Err("wrong human time format")
)

// ParseHuman convert human time format to duration.
// support:
// 	'H': hour,
// 	'M': minute,
// 	'S': second,
// 	'm': millsecond,
// 	'u': microsecond,
// 	'n': nanosecond
func ParseHuman(timestr string) (time.Duration, error) {
	var t, counter time.Duration
	for i, l := 0, len(timestr); i < l; i++ {
		c := timestr[i]
		if c >= '0' && c <= '9' {
			counter = counter*10 + time.Duration(c-'0')
			continue
		}

		switch c {
		case 'H':
			t += counter * time.Hour
		case 'M':
			t += counter * time.Minute
		case 'S':
			t += counter * time.Second
		case 'm':
			t += counter * time.Millisecond
		case 'u':
			t += counter * time.Microsecond
		case 'n':
			t += counter * time.Nanosecond
		default:
			return 0, ErrWrongHumanTimeFormat
		}
		counter = 0
	}

	return t, nil
}

// MonthDays return days of the month/year
func MonthDays(year, month int) int {
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	case 2:
		if IsLeapYear(year) {
			return 29
		}
		return 28
	default:
		panic(fmt.Sprintf("Illegal month:%d", month))
	}
}

func YearDays(year int) int {
	if IsLeapYear(year) {
		return 366
	}
	return 365
}

// IsLeapYear check whether a year is leay
func IsLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}
