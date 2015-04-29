package log

import (
	"fmt"
	"io"

	"github.com/cosiner/gohper/lib/errors"
	"github.com/cosiner/gohper/lib/time"
	"github.com/cosiner/gohper/lib/types"
)

const (
	LEVEL_TRACE Level = iota
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_PANIC
	LEVEL_FATAL
	LEVEL_OFF

	LEVEL_ALL  = LEVEL_TRACE
	_LEVEL_MIN = LEVEL_TRACE
	_LEVEL_MAX = LEVEL_FATAL

	ErrUnknownLevel = errors.Err("Unknown log level")
)

var (
	// levelName specified the all log level name
	levelName  = [...]string{"TRACE", "INFO", "WARN", "PANIC", "FATAL", "OFF"}
	timeFormat = time.FormatLayout("yyyymmdd-HHMMSS")
)

type (
	// Level is log level,
	// TRACE, INFO, WARN, PANIC, FATAL,
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

// String return level name, if level is no more than level_off, return actual name
// else panic
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
