package log

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/cosiner/golib/errors"
	t "github.com/cosiner/golib/time"
)

var (
	timeNow  = time.Now
	dateTime = t.DateTime
)

// Log represend a log with level and log message
type Log struct {
	Level   Level
	Message string
	Time    string
}

// String return a log as string with format "[level] time message"
func (l *Log) String() string {
	return fmt.Sprintf("[%s] %s %s", l.Level.String(), l.Time, l.Message)
}

// buildLog format log
func NewLogf(level Level, format string, v ...interface{}) *Log {
	return &Log{
		Level:   level,
		Message: fmt.Sprintf(format, v...),
		Time:    dateTime(),
	}
}

func NewLog(level Level, v ...interface{}) *Log {
	return &Log{
		Level:   level,
		Message: fmt.Sprint(v...),
		Time:    dateTime(),
	}
}

func NewLogln(level Level, v ...interface{}) *Log {
	return &Log{
		Level:   level,
		Message: fmt.Sprintln(v...),
		Time:    dateTime(),
	}
}

// LogWriter is actual log writer
type LogWriter interface {
	Config(conf string) error
	Write(log *Log) error
	ResetLevel(level Level) error
	Flush()
	Close()
}

type ConsoleLogWriter struct {
}

func (clw *ConsoleLogWriter) Config(conf string) error {
	return nil
}

func (clw *ConsoleLogWriter) Write(log *Log) error {
	_, err := fmt.Printf("[%s] %s %s", log.Level.String(), dateTime(), log.Message)
	return err
}

func (clw *ConsoleLogWriter) ResetLevel(level Level) {}
func (clw *ConsoleLogWriter) Flush()                 {}
func (clw *ConsoleLogWriter) Close()                 {}

type LoggerSignal uint8

const (
	SIGNAL_FLUSH LoggerSignal = iota // flush all writer
	SIGNAL_STOP                      // stop logger
	SIGNAL_PAUSE                     // pause logger
	SIGNAL_EXIT                      // exit process
)

// Logger
type Logger struct {
	*sync.RWMutex
	level         Level
	writers       []LogWriter
	flushInterval time.Duration
	logs          chan *Log
	signal        chan LoggerSignal
	running       bool
}

// NewLogger return a logger, if parameter is wrong, it will be automicly set to default value
// default use file logger, and default started
func NewLogger(flushInterval int, level Level) *Logger {
	errors.Assert(level >= LEVEL_MIN && level <= LEVEL_MAX,
		UnknownLevelErr(level.String()))
	errors.Assert(flushInterval > 0,
		errors.Errorf("Flush interval should not be negative:%d", flushInterval))
	return &Logger{RWMutex: new(sync.RWMutex),
		level:         level,
		logs:          make(chan *Log, DEF_BACKLOG),
		signal:        make(chan LoggerSignal),
		flushInterval: time.Duration(flushInterval) * time.Second,
		running:       false,
	}
}

// SetLogWriter set log writer, writer must not be nil
func (logger *Logger) AddLogWriter(writer LogWriter) (err error) {
	if writer != nil {
		logger.Lock()
		logger.writers = append(logger.writers, writer)
		err = writer.ResetLevel(logger.level)
		logger.Unlock()
	}
	return
}

// LogLevel return logger's level
func (logger *Logger) LogLevel() (l Level) {
	logger.RLock()
	l = logger.level
	logger.RUnlock()
	return
}

// SetLevel change logger's level
func (logger *Logger) SetLevel(level Level) (err error) {
	errors.Assert(level >= LEVEL_MIN && level <= LEVEL_MAX,
		UnknownLevelErr(level.String()))
	logger.Lock()
	logger.level = level
	for _, writer := range logger.writers {
		err = writer.ResetLevel(level)
		if err != nil {
			break
		}
	}
	logger.Unlock()
	return
}

func (logger *Logger) Signal(signal LoggerSignal) {
	logger.signal <- signal
}

// Start start logger
func (logger *Logger) Start() {
	logger.Lock()
	logger.running = true
	logger.Unlock()
	go func() {
		ticker := time.Tick(logger.flushInterval)
		for {
			select {
			case log := <-logger.logs:
				logger.RLock()
				for _, writer := range logger.writers {
					writer.Write(log)
				}
				logger.Unlock()
			case <-ticker:
				logger.Lock()
				for _, writer := range logger.writers {
					writer.Flush()
				}
				logger.Unlock()
			case signal := <-logger.signal:
				logger.Lock()
				switch signal {
				case SIGNAL_FLUSH:
					for _, writer := range logger.writers {
						writer.Flush()
					}
				case SIGNAL_PAUSE:
					logger.running = false
				case SIGNAL_STOP:
					logger.running = false
					logger.Unlock()
					return
				case SIGNAL_EXIT:
					logger.running = false
					for _, writer := range logger.writers {
						writer.Close()
					}
					logger.Unlock()
					os.Exit(1)
				}
				logger.Unlock()
			}
		}
	}()
}

// logf format the log, send it to log write's goroutine
func (logger *Logger) logf(level Level, format string, v ...interface{}) {
	logger.RLock()
	r, l := logger.running, logger.level
	logger.RUnlock()
	if r && level >= l {
		logger.logs <- NewLogf(level, format, v...)
	}
}

// logln  send log to log write's goroutine with an new line
func (logger *Logger) logln(level Level, v ...interface{}) {
	logger.RLock()
	r, l := logger.running, logger.level
	logger.RUnlock()
	if r && level >= l {
		logger.logs <- NewLogln(level, v...)
	}
}

// log send log to log write's goroutine
func (logger *Logger) log(level Level, v ...interface{}) {
	logger.RLock()
	r, l := logger.running, logger.level
	logger.RUnlock()
	if r && level >= l {
		logger.logs <- NewLog(level, v...)
	}
}
