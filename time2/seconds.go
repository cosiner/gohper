package time2

import "time"

var TimeSince *time.Time

func Seconds() int {
	now := Now()
	if TimeSince == nil {
		return int(now.Unix())
	}
	return int(now.Sub(*TimeSince).Seconds())
}
