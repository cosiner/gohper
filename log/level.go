package log

import (
	"github.com/cosiner/golib/types"
)

// Level is log level,
// DEBUG, INFO, WARN, ERROR, FATAL, OFF
type Level uint8

// levelName specified the all log level name
var levelName = [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL", "OFF"}

const (
	LEVEL_DEBUG Level = iota
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
	LEVEL_FATAL
	LEVEL_OFF
	_LEVEL_NUM = LEVEL_OFF // level count, don't include LEVEL_OFF
	LEVEL_ALL  = LEVEL_DEBUG
	LEVEL_MIN  = LEVEL_DEBUG
	// LEVEL_MAX is mean LEVEL_OFF which is not an actual log level, only for check
	LEVEL_MAX         = _LEVEL_NUM
	DEF_FLUSHINTERVAL = 30               // flush interval for a flush timer
	DEF_BUFSIZE       = 1024 * 10        // bufsize for log buffer
	DEF_BACKLOG       = 10               // channel's back log count
	DEF_FILESIZE      = 1024 * 1024 * 10 // max log file size
	DEF_LEVEL         = LEVEL_INFO       // default log level
)

// String return level name, if level is no more than level_off, return actual name
// else return UNKNOWN
func (l Level) Name() string {
	if l <= LEVEL_OFF {
		return levelName[l]
	} else {
		return "UNKNOWN"
	}
}

// ParseLevel parse level from string regardless of string case
// on error, return LEVEL_OFF
func ParseLevel(str string) (level Level) {
	levelStr := types.TrimUpper(str)
	level = LEVEL_OFF
	for l := LEVEL_MIN; l <= LEVEL_MAX; l++ {
		if levelStr == levelName[l] {
			level = l
		}
	}
	return
}
