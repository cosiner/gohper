package log

import (
	"fmt"
	"io"

	"github.com/cosiner/gohper/lib/errors"
	"github.com/cosiner/gohper/lib/time"
	"github.com/cosiner/gohper/lib/types"
)

type (
	// Level is log level,
	// DEBUG, INFO, WARN, ERROR, FATAL,
	Level uint8
	// Log represend a log with level and log message
	Log struct {
		Level Level
		Time  string

		format  string
		newline bool
		depth   string
		args    []interface{}
	}
)

var (
	// levelName specified the all log level name
	levelName  = [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL", "OFF"}
	timeFormat = time.FormatLayout("yyyymmdd-HHMMSS")
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

	ErrUnknownLevel = errors.Err("Unknown log level")
)

// String return level name, if level is no more than level_off, return actual name
// else return UNKNOWN
func (l Level) String() string {
	if l >= _LEVEL_MIN && l <= LEVEL_OFF {
		return levelName[l]
	}
	panic(ErrUnknownLevel)
}

// ParseLevel parse level from string regardless of string case
func ParseLevel(str string) Level {
	s := types.TrimUpper(str)
	for l := _LEVEL_MIN; l <= LEVEL_OFF; l++ {
		if s == levelName[l] {
			return l
		}
	}
	panic(ErrUnknownLevel)
}

// String return a log as string with format "[level] time message"
func (l *Log) WriteTo(w io.Writer) error {
	_, err := fmt.Fprintf(w, "[%5s] %s ", l.Level, l.Time)
	if l.depth != "" {
		_, err = w.Write(types.UnsafeBytes(l.depth))
		_, err = w.Write(types.UnsafeBytes(":"))
	}
	if l.format != "" {
		_, err = fmt.Fprintf(w, l.format, l.args...)
	} else if l.newline {
		_, err = fmt.Fprintln(w, l.args...)
	} else {
		_, err = fmt.Fprint(w, l.args...)
	}
	return err
}

func logf(level Level, format string, args ...interface{}) *Log {
	return &Log{
		Level:  level,
		Time:   time.DateTime(),
		format: format,
		args:   args,
	}
}

func log(level Level, args ...interface{}) *Log {
	return &Log{
		Level: level,
		Time:  time.DateTime(),
		args:  args,
	}
}

func logln(level Level, args ...interface{}) *Log {
	return &Log{
		Level:   level,
		Time:    time.DateTime(),
		newline: true,
		args:    args,
	}
}

func logDepth(level Level, depth string, args ...interface{}) *Log {
	return &Log{
		Level:   level,
		Time:    time.DateTime(),
		newline: true,
		depth:   depth,
		args:    args,
	}
}
