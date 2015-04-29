package log

import (
	"github.com/cosiner/gohper/lib/runtime"
	"os"
	"time"

	"github.com/cosiner/gohper/lib/defval"
	"github.com/cosiner/gohper/lib/errors"
)

const (
	ErrInvalidConfig = errors.Err("invalid config parameters or type")
)

type (
	Logger interface {
		AddWriter(Writer, interface{}) error
		Depth(func(depth int) string)
		Level() Level
		Flush()
		Close()

		Debug(...interface{})
		Info(...interface{})
		Warn(...interface{})
		Error(...interface{}) // panic goroutine
		Fatal(...interface{}) // exit process

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

		DebugDepth(int, ...interface{})
		InfoDepth(int, ...interface{})
		WarnDepth(int, ...interface{})
		ErrorDepth(int, ...interface{})
		FatalDepth(int, ...interface{})
	}

	// Writer is actual log writer
	Writer interface {
		// SetLevel will be called before Config
		SetLevel(Level)
		// Config config writer
		Config(interface{}) error
		// Writer output log
		Write(log *Log) error
		// Flush flush output
		Flush()
		// Close close log writer
		Close()
	}

	LoggerOption struct {
		Level
		Flush   int // by seconds
		Backlog int
	}

	// Logger
	logger struct {
		level         Level
		writers       []Writer
		flushInterval time.Duration
		logs          chan *Log
		flush, exit   chan struct{}

		depthFunc 	  func(int)string
	}
)

func (opt *LoggerOption) init() {
	defval.Int(&opt.Flush, 30)
	defval.Int(&opt.Backlog, 100)
}

// NewLogger return a logger, if params is wrong, use default value
func New(opt *LoggerOption) Logger {
	if opt == nil {
		opt = &LoggerOption{}
	}
	opt.init()
	if opt.Level == LEVEL_OFF {
		return &logger{level: LEVEL_OFF}
	}
	l := &logger{
		level:         opt.Level,
		logs:          make(chan *Log, opt.Backlog),
		flush:         make(chan struct{}, 1),
		exit:          make(chan struct{}, 1),
		flushInterval: time.Duration(opt.Flush) * time.Second,
		depthFunc:	runtime.CallerPosition,
	}
	l.start()
	return l
}

func (logger *logger) AddWriter(w Writer, conf interface{}) error {
	if logger.level == LEVEL_OFF {
		return nil
	}

	w.SetLevel(logger.level)
	err := w.Config(conf)
	if err == nil {
		logger.writers = append(logger.writers, w)
	}
	return err
}

func (logger *logger) Depth(d func(int)string) {
	logger.depthFunc = d
}

// Level return logger's level
func (logger *logger) Level() (l Level) {
	return logger.level
}

// start start logger
func (l *logger) start() {
	go func(l *logger) {
		ticker := time.Tick(l.flushInterval)
		for {
			select {
			case log := <-l.logs:
				l.processLogs(log)
			case <-ticker:
				l.processFlush()
			case <-l.flush:
				l.processFlush()
			case <-l.exit:
				l.processClose()
			}
		}
	}(l)
}

func (logger *logger) processFlush() {
	for _, writer := range logger.writers {
		writer.Flush()
	}
}

func (logger *logger) processLogs(log *Log) {
	for _, writer := range logger.writers {
		writer.Write(log)
	}
	if log.Level == LEVEL_FATAL {
		for _, w := range logger.writers {
			w.Close()
		}
		os.Exit(-1)
	}
}

func (logger *logger) processClose() {
	for len(logger.logs) > 0 {
		logger.processLogs(<-logger.logs)
	}
	for _, w := range logger.writers {
		w.Close()
	}
}

// Flush flush logger
func (logger *logger) Flush() {
	logger.flush <- struct{}{}
}

func (logger *logger) Close() {
	logger.level = LEVEL_OFF
	logger.exit <- struct{}{}
}

func (logger *logger) printf(level Level, format string, args ...interface{}) {
	if level >= logger.level {
		logger.logs <- logf(level, format, args...)
	}
}

func (logger *logger) println(level Level, args ...interface{}) {
	if level >= logger.level {
		logger.logs <- logln(level, args...)
	}
}

func (logger *logger) print(level Level, args ...interface{}) {
	if level >= logger.level {
		logger.logs <- log(level, args...)
	}
}

func (logger *logger) printDepth(level Level, depth int, args ...interface{}) {
	if level >= logger.level {
		logger.logs <- logDepth(level, logger.depthFunc(depth + 1), args...)
	}
}

// Debugf log for debug message
func (logger *logger) Debugf(format string, args ...interface{}) {
	logger.printf(LEVEL_DEBUG, format, args...)
}

// Infof log for info message
func (logger *logger) Infof(format string, args ...interface{}) {
	logger.printf(LEVEL_INFO, format, args...)
}

// Warnf log for warning message
func (logger *logger) Warnf(format string, args ...interface{}) {
	logger.printf(LEVEL_WARN, format, args...)
}

// Errorf log for error message
func (logger *logger) Errorf(format string, args ...interface{}) {
	logger.printf(LEVEL_ERROR, format, args)
	panic(log)
}

// Fatalf log for fatal message
func (logger *logger) Fatalf(format string, args ...interface{}) {
	logger.printf(LEVEL_FATAL, format, args)
}

// Debugln log for debug message
func (logger *logger) Debugln(args ...interface{}) {
	logger.println(LEVEL_DEBUG, args...)
}

// Infoln log for info message
func (logger *logger) Infoln(args ...interface{}) {
	logger.println(LEVEL_INFO, args...)
}

// Warnln log for warning message
func (logger *logger) Warnln(args ...interface{}) {
	logger.println(LEVEL_WARN, args...)
}

// Errorln log for error message
func (logger *logger) Errorln(args ...interface{}) {
	logger.println(LEVEL_ERROR, args...)
	panic(log)
}

// Fatalln log for fatal message
func (logger *logger) Fatalln(args ...interface{}) {
	logger.println(LEVEL_FATAL, args...)
}

// Debug log for debug message
func (logger *logger) Debug(args ...interface{}) {
	logger.print(LEVEL_DEBUG, args...)
}

// Info log for info message
func (logger *logger) Info(args ...interface{}) {
	logger.print(LEVEL_INFO, args...)
}

// Warn log for warning message
func (logger *logger) Warn(args ...interface{}) {
	logger.print(LEVEL_WARN, args...)
}

// Error log for error message
func (logger *logger) Error(args ...interface{}) {
	logger.print(LEVEL_ERROR, args...)
	panic(log)
}

// Fatal log for error message
func (logger *logger) Fatal(args ...interface{}) {
	logger.print(LEVEL_FATAL, args...)
}

// Debug log for debug message
func (logger *logger) DebugDepth(depth int, args ...interface{}) {
	logger.printDepth(LEVEL_DEBUG, depth+1, args...)
}

// Info log for info message
func (logger *logger) InfoDepth(depth int, args ...interface{}) {
	logger.printDepth(LEVEL_INFO, depth+1, args...)
}

// Warn log for warning message
func (logger *logger) WarnDepth(depth int, args ...interface{}) {
	logger.printDepth(LEVEL_WARN, depth+1, args...)
}

// Error log for error message
func (logger *logger) ErrorDepth(depth int, args ...interface{}) {
	logger.printDepth(LEVEL_ERROR, depth+1, args...)
	panic(log)
}

// Fatal log for error message
func (logger *logger) FatalDepth(depth int, args ...interface{}) {
	logger.printDepth(LEVEL_FATAL, depth+1, args...)
}
