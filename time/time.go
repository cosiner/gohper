package time

import (
	"strings"
	"time"
)

const DATETIME_FMT = "2006/01/02 15:04:05"
const DATE_FMT = "2006/01/02"
const TIME_FMT = "15:04:05"

// FormatLayout convert human readable time format to golang time format
// yyyy : year, yy:year, mm:month, dd:day, HH:hour, MM:minute, SS:second
func FormatLayout(format string) string {
	format = strings.Replace(format, "yyyy", "2006", 1)
	strings.Replace(format, "yy", "06", 1)
	strings.Replace(format, "mm", "01", 1)
	strings.Replace(format, "dd", "02", 1)
	strings.Replace(format, "HH", "15", 1)
	strings.Replace(format, "MM", "04", 1)
	strings.Replace(format, "SS", "05", 1)
	return string(format)
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
