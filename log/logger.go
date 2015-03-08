package log

import (
	"time"

	"github.com/cosiner/golib/errors"
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
	logger.writers = append(logger.writers, writer)
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

func (logger *logger) logf(level Level, format string, v ...interface{}) {
	if level >= logger.level {
		logger.logs <- NewLogf(level, format, v...)
	}
}

func (logger *logger) logln(level Level, v ...interface{}) {
	if level >= logger.level {
		logger.logs <- NewLogln(level, v...)
	}
}

func (logger *logger) log(level Level, v ...interface{}) {
	if level >= logger.level {
		logger.logs <- NewLog(level, v...)
	}
}

// Debugf log for debug message
func (l *logger) Debugf(format string, v ...interface{}) {
	l.logf(LEVEL_DEBUG, format, v...)
}

// Infof log for info message
func (l *logger) Infof(format string, v ...interface{}) {
	l.logf(LEVEL_INFO, format, v...)
}

// Warnf log for warning message
func (l *logger) Warnf(format string, v ...interface{}) {
	l.logf(LEVEL_WARN, format, v...)
}

// Errorf log for error message
func (l *logger) Errorf(format string, v ...interface{}) {
	l.logf(LEVEL_ERROR, format, v...)
}

// Fatalf log for fatal message
func (l *logger) Fatalf(format string, v ...interface{}) {
	log := NewLogf(LEVEL_FATAL, format, v...)
	if LEVEL_FATAL > l.level {
		l.logs <- log
	}
	panic(log)
}

// Debugln log for debug message
func (l *logger) Debugln(v ...interface{}) {
	l.logln(LEVEL_DEBUG, v...)
}

// Infoln log for info message
func (l *logger) Infoln(v ...interface{}) {
	l.logln(LEVEL_INFO, v...)
}

// Warnln log for warning message
func (l *logger) Warnln(v ...interface{}) {
	l.logln(LEVEL_WARN, v...)
}

// Errorln log for error message
func (l *logger) Errorln(v ...interface{}) {
	l.logln(LEVEL_ERROR, v...)
}

// Fatalln log for fatal message
func (l *logger) Fatalln(v ...interface{}) {
	log := NewLogln(LEVEL_FATAL, v...)
	l.logs <- log
	panic(log)
}

// Debug log for debug message
func (l *logger) Debug(v ...interface{}) {
	l.log(LEVEL_DEBUG, v...)
}

// Info log for info message
func (l *logger) Info(v ...interface{}) {
	l.log(LEVEL_INFO, v...)
}

// Warn log for warning message
func (l *logger) Warn(v ...interface{}) {
	l.log(LEVEL_WARN, v...)
}

// Error log for error message
func (l *logger) Error(v ...interface{}) {
	l.log(LEVEL_ERROR, v...)
}

// Fatal log for error message
func (l *logger) Fatal(v ...interface{}) {
	log := NewLog(LEVEL_FATAL, v...)
	l.logs <- log
	panic(log)
}
