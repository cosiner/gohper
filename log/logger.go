package log

import (
	"fmt"
	"time"

	"github.com/cosiner/gohper/lib/runtime"

	"github.com/cosiner/gohper/lib/errors"
)

type (
	Logger interface {
		AddLogWriter(LogWriter)
		AddLogWriterWithConf(LogWriter, string) error
		Start()
		LogLevel() Level
		SetLevel(Level) error
		Flush()

		Debugf(string, ...interface{})
		Infof(string, ...interface{})
		Warnf(string, ...interface{})
		Errorf(string, ...interface{})
		Fatalf(string, ...interface{})
		Debugln(...interface{})
		Infoln(...interface{})
		Warnln(...interface{})
		Errorln(...interface{})
		Fatalln(...interface{})
		Debug(...interface{})
		Info(...interface{})
		Warn(...interface{})
		Error(...interface{})
		Fatal(...interface{})
	}

	// LogWriter is actual log writer
	LogWriter interface {
		// Config config writer
		Config(conf string) error
		// Writer output log
		Write(log *Log) error
		// Flush flush output
		Flush()
		// Close close log writer
		Close()
	}

	// Logger
	logger struct {
		level         Level
		writers       []LogWriter
		flushInterval time.Duration
		logs          chan *Log
		signal        chan byte
	}
)

const (
	_SIGNAL_FLUSH byte = iota // flush all writer
)

// NewLogger return a logger, if params is wrong, use default value
func NewLogger(flushInterval int, level Level) Logger {
	if level < _LEVEL_MIN || level > LEVEL_OFF {
		level = DEF_LEVEL
	}
	if flushInterval <= 0 {
		flushInterval = DEF_FLUSHINTERVAL
	}
	return &logger{
		level:         level,
		logs:          make(chan *Log, DEF_BACKLOG),
		signal:        make(chan byte, 1),
		flushInterval: time.Duration(flushInterval) * time.Second,
	}
}

// AddLogWroter add a  log writer, nil writer will be auto-ignored
func (logger *logger) AddLogWriter(writer LogWriter) {
	if logger.level < LEVEL_OFF {
		logger.writers = append(logger.writers, writer)
	}
}

func (logger *logger) AddLogWriterWithConf(writer LogWriter, conf string) error {
	err := writer.Config(conf)
	if err == nil {
		logger.AddLogWriter(writer)
	}
	return err
}

// LogLevel return logger's level
func (logger *logger) LogLevel() (l Level) {
	return logger.level
}

// SetLevel change logger's level, it will apply to all log writers
func (logger *logger) SetLevel(level Level) (err error) {
	errors.Assert(level >= _LEVEL_MIN && level <= _LEVEL_MAX,
		UnknownLevelErr(level.String()).Error())
	logger.level = level
	return
}

// Start start logger
func (logger *logger) Start() {
	go func() {
		ticker := time.Tick(logger.flushInterval)
		for {
			select {
			case log := <-logger.logs:
				for _, writer := range logger.writers {
					writer.Write(log)
				}
			case <-ticker:
				for _, writer := range logger.writers {
					writer.Flush()
				}
			case <-logger.signal:
				for _, writer := range logger.writers {
					writer.Flush()
				}
			}
		}
	}()
}

// Flush flush logger
func (logger *logger) Flush() {
	logger.signal <- _SIGNAL_FLUSH
}

func (logger *logger) logf(level Level, format string, v ...interface{}) *Log {
	if level >= logger.level {
		log := NewLogf(level, format, v...)
		logger.logs <- log
		return log
	}
	return nil
}

func (logger *logger) logln(level Level, v ...interface{}) *Log {
	if level >= logger.level {
		log := NewLogln(level, v...)
		logger.logs <- log
		return log
	}
	return nil
}

func (logger *logger) log(level Level, v ...interface{}) *Log {
	if level >= logger.level {
		log := NewLog(level, v...)
		logger.logs <- log
		return log
	}
	return nil
}

// Debugf log for debug message
func (logger *logger) Debugf(format string, v ...interface{}) {
	format = fmt.Sprintf("%s %s", runtime.CallerPosition(1), format)
	logger.logf(LEVEL_DEBUG, format, v...)
}

// Infof log for info message
func (logger *logger) Infof(format string, v ...interface{}) {
	logger.logf(LEVEL_INFO, format, v...)
}

// Warnf log for warning message
func (logger *logger) Warnf(format string, v ...interface{}) {
	logger.logf(LEVEL_WARN, format, v...)
}

// Errorf log for error message
func (logger *logger) Errorf(format string, v ...interface{}) {
	if log := logger.logf(LEVEL_ERROR, format, v...); logger.level == LEVEL_DEBUG {
		panic(log)
	}
}

// Fatalf log for fatal message
func (logger *logger) Fatalf(format string, v ...interface{}) {
	if log := logger.logf(LEVEL_FATAL, format, v...); log != nil {
		panic(log)
	}
}

// Debugln log for debug message
func (logger *logger) Debugln(v ...interface{}) {
	logger.logln(LEVEL_DEBUG, append([]interface{}{runtime.CallerPosition(1)}, v...)...)
}

// Infoln log for info message
func (logger *logger) Infoln(v ...interface{}) {
	logger.logln(LEVEL_INFO, v...)
}

// Warnln log for warning message
func (logger *logger) Warnln(v ...interface{}) {
	logger.logln(LEVEL_WARN, v...)
}

// Errorln log for error message
func (logger *logger) Errorln(v ...interface{}) {
	if log := logger.logln(LEVEL_ERROR, v...); logger.level == LEVEL_DEBUG {
		panic(log)
	}
}

// Fatalln log for fatal message
func (logger *logger) Fatalln(v ...interface{}) {
	if log := logger.logln(LEVEL_FATAL, v...); log != nil {
		panic(log)
	}
}

// Debug log for debug message
func (logger *logger) Debug(v ...interface{}) {
	logger.log(LEVEL_DEBUG, append([]interface{}{runtime.CallerPosition(1)}, v...)...)
}

// Info log for info message
func (logger *logger) Info(v ...interface{}) {
	logger.log(LEVEL_INFO, v...)
}

// Warn log for warning message
func (logger *logger) Warn(v ...interface{}) {
	logger.log(LEVEL_WARN, v...)
}

// Error log for error message
func (logger *logger) Error(v ...interface{}) {
	if log := logger.log(LEVEL_ERROR, v...); logger.level == LEVEL_DEBUG {
		panic(log)
	}
}

// Fatal log for error message
func (logger *logger) Fatal(v ...interface{}) {
	if log := logger.log(LEVEL_FATAL, v...); log != nil {
		panic(log)
	}
}
