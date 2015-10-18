package time2

import (
	"fmt"
	"time"

	"github.com/cosiner/gohper/errors"
)

var (
	Location     = time.Local
	DATETIME_FMT = "2006/01/02 15:04:05"
	DATE_FMT     = "2006/01/02"
	TIME_FMT     = "15:04:05"
)

func Now() time.Time {
	now := time.Now()
	if Location == time.Local {
		return now
	}

	return now.In(Location)
}

func CurrDate() string {
	return Date(Now())
}

func Date(t time.Time) string {
	return t.Format(DATE_FMT)
}

func CurrTime() string {
	return Time(Now())
}

func Time(t time.Time) string {
	return t.Format(TIME_FMT)
}

func CurrDateTime() string {
	return DateTime(Now())
}

func DateTime(t time.Time) string {
	return t.Format(DATETIME_FMT)
}

func CurrFormat(layout string) string {
	return Now().Format(layout)
}

func Format(t time.Time, layout string) string {
	return t.Format(layout)
}

func CurrDateAndTime() (string, string) {
	now := Now()
	return Date(now), Time(now)
}

func Parse(layout, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, Location)
}

func ParseDate(value string) (time.Time, error) {
	return Parse(DATE_FMT, value)
}

func ParseTime(value string) (time.Time, error) {
	return Parse(TIME_FMT, value)
}
func ParseDateTime(value string) (time.Time, error) {
	return Parse(DATETIME_FMT, value)
}

func UnixNanoSinceNow(sec int) int64 {
	return Now().Add(time.Duration(sec) * time.Second).UnixNano()
}

// Unix is a wrapper of Now().Unix()
func Unix() uint64 {
	return uint64(Now().Unix())
}

// UnixNano is a wrapper of Now().UnixNano()
func UnixNano() uint64 {
	return uint64(Now().UnixNano())
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

// DateDefNow create a timestamp with given field, default use value of now if a
// field is less than 0
func DateDefNow(year, month, day, hour, minute, sec, nsec int) time.Time {
	now := Now()
	nyear, nmonth, nday := now.Date()
	nhour, nminute, nsec := now.Clock()
	if year < 0 {
		year = nyear
	}
	if month < 0 {
		month = int(nmonth)
	}
	if day < 0 {
		day = nday
	}
	if hour < 0 {
		hour = nhour
	}
	if minute < 0 {
		minute = nminute
	}
	if sec < 0 {
		sec = nsec
	}

	return time.Date(year, time.Month(month), day, hour, minute, sec, nsec, Location)
}
