package time

import (
	"strings"
	"time"
)

const DATETIME_FMT = "2006/01/02 15:04:05"
const DATE_FMT = "2006/01/02"
const TIME_FMT = "15:04:05"

var formatReplacer = strings.NewReplacer([]string{
	"yyyy", "2006",
	"yy", "06",
	"mm", "01",
	"dd", "02",
	"HH", "15",
	"MM", "04",
	"SS", "05",
}...)

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
