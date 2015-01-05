// Package log implements a level logger
package log

import (
	"mlib/util/sys"
)

// null function set
func nullFunc()        {}
func nullLevel() Level { return LEVEL_OFF }
func nullSetLevel(level Level) error {
	return nil
}
func nullLog(format string)                    {}
func nullLogf(format string, v ...interface{}) {}

var (
	L                                    *Logger
	onceInit                             = new(sys.Once)
	LogLevel                             = nullLevel
	SetLevel                             = nullSetLevel
	LockFlush, FlushAndExit              = nullFunc, nullFunc
	Debugf, Infof, Warnf, Errorf, Fatalf = nullLogf, nullLogf, nullLogf, nullLogf, nullLogf
	Debug, Info, Warn, Error, Fatal      = nullLog, nullLog, nullLog, nullLog, nullLog
)

// Init init the logger, and log functions, it will init only once unless on error
func Init(flushInterval int, level Level, bufSize, maxSize uint64, logDir string) error {
	if level == LEVEL_OFF {
		return nil
	}
	return onceInit.DoCheckError(func() (err error) {
		L, err = NewLogger(flushInterval, level, bufSize, maxSize, logDir)
		if err != nil {
			return
		}
		LogLevel = L.LogLevel
		SetLevel = L.SetLevel
		LockFlush = L.LockFlush
		FlushAndExit = L.FlushAndExit

		Debug = func(format string) {
			L.Debugf(format)
		}
		Debugf = L.Debugf

		Info = func(format string) {
			L.Infof(format)
		}
		Infof = L.Infof

		Warn = func(format string) {
			L.Warnf(format)
		}
		Warnf = L.Warnf

		Error = func(format string) {
			L.Errorf(format)
		}
		Errorf = L.Errorf

		Fatal = func(format string) {
			L.Fatalf(format)
		}
		Fatalf = L.Fatalf
		return
	})
}
