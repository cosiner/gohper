package log

import (
	"fmt"

	"github.com/cosiner/gohper/lib/time"

	"github.com/cosiner/gohper/lib/errors"
	"github.com/cosiner/gohper/lib/types"
)

type (
	// Level is log level,
	// DEBUG, INFO, WARN, ERROR, FATAL,
	Level uint8
	// Log represend a log with level and log message
	Log struct {
		Level   Level
		Message string
		Time    string
	}
)

var (
	// levelName specified the all log level name
	levelName = [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL", "OFF"}
)

const (
	LEVEL_DEBUG Level = iota
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
	LEVEL_FATAL
	LEVEL_OFF

	LEVEL_ALL  = LEVEL_DEBUG
	_LEVEL_MIN = LEVEL_DEBUG
	_LEVEL_MAX = LEVEL_FATAL

	DEF_FLUSHINTERVAL = 30               // flush interval for a flush timer
	DEF_BUFSIZE       = 1024 * 10        // bufsize for log buffer
	DEF_BACKLOG       = 100              // channel's back log count
	DEF_FILESIZE      = 1024 * 1024 * 10 // max log file size
	DEF_LEVEL         = LEVEL_INFO       // default log level

	ErrUnknownLevel = errors.Err("Unknown log level")
)

// String return level name, if level is no more than level_off, return actual name
// else return UNKNOWN
func (l Level) String() string {
	if l <= _LEVEL_MAX {
		return levelName[l]
	}
	panic(ErrUnknownLevel)
}

// ParseLevel parse level from string regardless of string case
func ParseLevel(str string) Level {
	levelStr := types.TrimUpper(str)
	for l := _LEVEL_MIN; l <= LEVEL_OFF; l++ {
		if levelStr == levelName[l] {
			return l
		}
	}
	panic(ErrUnknownLevel)
}

// String return a log as string with format "[level] time message"
func (l *Log) String() string {
	return fmt.Sprintf("[%5s] %s %s", l.Level.String(), l.Time, l.Message)
}

// buildLog format log
func NewLogf(level Level, format string, v ...interface{}) *Log {
	return &Log{
		Level:   level,
		Message: fmt.Sprintf(format, v...),
		Time:    time.DateTime(),
	}
}

func NewLog(level Level, v ...interface{}) *Log {
	return &Log{
		Level:   level,
		Message: fmt.Sprint(v...),
		Time:    time.DateTime(),
	}
}

func NewLogln(level Level, v ...interface{}) *Log {
	return &Log{
		Level:   level,
		Message: fmt.Sprintln(v...),
		Time:    time.DateTime(),
	}
}
