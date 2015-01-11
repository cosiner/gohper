package log

import (
	"fmt"
	"os"
	"sync"
	"time"

	e "github.com/cosiner/golib/errors"

	t "github.com/cosiner/golib/time"
	"github.com/cosiner/golib/types"
)

// LogWriter is actual log writer
type LogWriter interface {
	Write(log *Log) error
	ResetLevel(level Level) error
	Flush()
	Close()
}

// Log represend a log with level and log message
type Log struct {
	level Level
	msg   string
}

// Logger
type Logger struct {
	mu            *sync.Mutex
	level         Level
	lw            LogWriter
	flushInterval time.Duration
	logs          chan *Log
	flushexit     chan bool
}

var (
	DEF_LOGDIR = os.TempDir() + "/logs"
	timeNow    = time.Now
	dateTime   = t.DateTime
)

// NewLogger return a logger, if parameter is wrong, it will be automicly set to default value
// default use file logger, and default started
func NewLogger(flushInterval int, level Level, bufSize, maxSize uint64, logDir string) (*Logger, error) {
	if level == LEVEL_MAX || level < LEVEL_MIN {
		return nil, e.Err("Invalid log level:" + level.Name())
	}

	if bufSize == 0 {
		bufSize = DEF_BUFSIZE
	}
	if maxSize == 0 {
		maxSize = DEF_FILESIZE
	}
	if types.TrimSpace(logDir) == "" {
		logDir = DEF_LOGDIR
	}

	var logger *Logger
	lw, err := newLogWriter(level, bufSize, maxSize, logDir)
	logger = NewEmptyLogger(flushInterval, level)
	if err == nil {
		logger.lw = lw
		logger.Start()
	}
	return logger, err
}

// NewEmptyLogger new an empty logger, can't be used directly
// to use it, must SetLogWriter, then Start it, and commonly,
// assign it to global variable log.L
func NewEmptyLogger(flushInterval int, level Level) *Logger {
	if flushInterval <= 0 {
		flushInterval = DEF_FLUSHINTERVAL
	}
	if level < LEVEL_MIN || level > LEVEL_OFF {
		level = DEF_LEVEL
	}

	logger := &Logger{mu: new(sync.Mutex),
		level:         level,
		lw:            nil,
		logs:          make(chan *Log, DEF_BACKLOG),
		flushexit:     make(chan bool),
		flushInterval: time.Duration(flushInterval) * time.Second}
	return logger
}

// SetLogWriter set log writer, writer must not be nil
func (logger *Logger) SetLogWriter(lw LogWriter) error {
	if lw == nil {
		return e.Err("LogWrite can't be nil")
	}
	if logger.lw != nil {
		logger.mu.Lock()
		logger.lw.Close()
		logger.mu.Unlock()
	}
	logger.lw = lw
	lw.ResetLevel(logger.level)
	return nil
}

// LogLevel return logger's level
func (logger *Logger) LogLevel() Level {
	return logger.level
}

// SetLevel change logger's level
func (logger *Logger) SetLevel(level Level) (err error) {
	logger.level = level
	if logger.lw != nil {
		logger.mu.Lock()
		err = logger.lw.ResetLevel(level)
		logger.mu.Unlock()
	}
	return
}

// lockFor is a convenience for function need to lock
func (logger *Logger) lockFor(fn func()) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	fn()
}

// lockWrite lock then call logwriter to write log
func (logger *Logger) lockWrite(log *Log) (err error) {
	logger.lockFor(func() {
		err = logger.lw.Write(log)
	})
	return
}

// LockFlush flush logwriter
func (logger *Logger) LockFlush() {
	logger.lockFor(func() {
		logger.lw.Flush()
	})
}

// FlushAndExit flush logger then exit
func (logger *Logger) FlushAndExit() {
	logger.flushexit <- true
}

// Start start logger
func (logger *Logger) Start() {
	go func() {
		ticker := time.Tick(logger.flushInterval)
		for {
			select {
			case log := <-logger.logs:
				logger.lockWrite(log)
			case <-ticker:
				logger.LockFlush()
			case fe := <-logger.flushexit:
				if fe {
					// if len(logger.logs) > 0 {
					// 	logger.FlushAndExit()
					// 	continue
					// }
					close(logger.flushexit)
					logger.lockFor(func() {
						close(logger.logs)
						logger.lw.Close()
					})
					os.Exit(0)
				}
			}
		}
	}()
}

// logf format the log, send it to log write's goroutine
func (logger *Logger) logf(level Level, format string, v ...interface{}) {
	if log := logger.buildLog(level, format, v...); log != nil {
		logger.logs <- log
	}
}

// buildLog format log
func (logger *Logger) buildLog(level Level, format string, v ...interface{}) (log *Log) {
	if logger.level != LEVEL_OFF && level >= logger.level {
		format = fmt.Sprintf("[%s]%s %s", level.Name(), dateTime(), format)
		msg := fmt.Sprintf(format, v...)
		log = &Log{level, msg}
	}
	return
}

// Debugf log for debug message
func (logger *Logger) Debugf(format string, v ...interface{}) {
	logger.logf(LEVEL_DEBUG, format, v...)
}

// Infof log for info message
func (logger *Logger) Infof(format string, v ...interface{}) {
	logger.logf(LEVEL_INFO, format, v...)
}

// Warnf log for warning message
func (logger *Logger) Warnf(format string, v ...interface{}) {
	logger.logf(LEVEL_WARN, format, v...)
}

// Errorf log for error message
func (logger *Logger) Errorf(format string, v ...interface{}) {
	logger.logf(LEVEL_ERROR, format, v...)
}

// Fatalf log for fatal message
func (logger *Logger) Fatalf(format string, v ...interface{}) {
	logger.lockWrite(logger.buildLog(LEVEL_FATAL, format, v...))
	logger.FlushAndExit()
}
