package log

import (
	"os"
	"time"

	"github.com/cosiner/gohper/lib/defval"
	"github.com/cosiner/gohper/lib/errors"
	"github.com/cosiner/gohper/lib/runtime"
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

		Trace(...interface{})
		Info(...interface{})
		Warn(...interface{})
		Panic(...interface{}) // panic goroutine
		Fatal(...interface{}) // exit process

		Tracef(string, ...interface{})
		Infof(string, ...interface{})
		Warnf(string, ...interface{})
		Panicf(string, ...interface{})
		Fatalf(string, ...interface{})

		Traceln(...interface{})
		Infoln(...interface{})
		Warnln(...interface{})
		Panicln(...interface{})
		Fatalln(...interface{})

		TraceDepth(int, ...interface{})
		InfoDepth(int, ...interface{})
		WarnDepth(int, ...interface{})
		PanicDepth(int, ...interface{})
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

		depthFunc func(int) string
	}
)

func (opt *LoggerOption) init() {
	defval.Int(&opt.Flush, 30)
	defval.Int(&opt.Backlog, 100)
}

// Default create a logger with console writer, use debug level
func Default() Logger {
	l := New(&LoggerOption{
		Level: _LEVEL_MIN,
	})
	l.AddWriter(new(ConsoleWriter), nil)
	return l
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
		depthFunc:     runtime.CallerPosition,
	}
	return l
}

func (l *logger) AddWriter(w Writer, conf interface{}) error {
	if l.level == LEVEL_OFF {
		return nil
	}

	w.SetLevel(l.level)
	err := w.Config(conf)
	if err == nil {
		l.writers = append(l.writers, w)
		if len(l.writers) == 1 {
			l.start()
		}
	}
	return err
}

func (l *logger) Depth(d func(int) string) {
	l.depthFunc = d
}

// Level return logger's level
func (l *logger) Level() Level {
	return l.level
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

func (l *logger) processFlush() {
	for _, writer := range l.writers {
		writer.Flush()
	}
}

func (l *logger) processLogs(log *Log) {
	for _, writer := range l.writers {
		writer.Write(log)
	}
	if log.Level == LEVEL_FATAL {
		for _, w := range l.writers {
			w.Close()
		}
		os.Exit(-1)
	}
}

func (l *logger) processClose() {
	for len(l.logs) > 0 {
		l.processLogs(<-l.logs)
	}
	for _, w := range l.writers {
		w.Close()
	}
}

// Flush flush logger
func (l *logger) Flush() {
	l.flush <- struct{}{}
}

func (l *logger) Close() {
	l.level = LEVEL_OFF
	l.exit <- struct{}{}
}

func (l *logger) printf(level Level, format string, args ...interface{}) {
	if level >= l.level {
		l.logs <- logf(level, format, args...)
	}
}

func (l *logger) println(level Level, args ...interface{}) {
	if level >= l.level {
		l.logs <- logln(level, args...)
	}
}

func (l *logger) print(level Level, args ...interface{}) {
	if level >= l.level {
		l.logs <- log(level, args...)
	}
}

func (l *logger) printDepth(level Level, depth int, args ...interface{}) {
	if level >= l.level {
		l.logs <- logDepth(level, l.depthFunc(depth+1), args...)
	}
}

// Tracef log for debug message
func (l *logger) Tracef(format string, args ...interface{}) {
	l.printf(LEVEL_TRACE, format, args...)
}

// Infof log for info message
func (l *logger) Infof(format string, args ...interface{}) {
	l.printf(LEVEL_INFO, format, args...)
}

// Warnf log for warning message
func (l *logger) Warnf(format string, args ...interface{}) {
	l.printf(LEVEL_WARN, format, args...)
}

// Panicf log for error message
func (l *logger) Panicf(format string, args ...interface{}) {
	l.printf(LEVEL_PANIC, format, args)
	panic(l.depthFunc(1))
}

// Fatalf log for fatal message
func (l *logger) Fatalf(format string, args ...interface{}) {
	l.printf(LEVEL_FATAL, format, args)
}

// Traceln log for debug message
func (l *logger) Traceln(args ...interface{}) {
	l.println(LEVEL_TRACE, args...)
}

// Infoln log for info message
func (l *logger) Infoln(args ...interface{}) {
	l.println(LEVEL_INFO, args...)
}

// Warnln log for warning message
func (l *logger) Warnln(args ...interface{}) {
	l.println(LEVEL_WARN, args...)
}

// Panicln log for error message
func (l *logger) Panicln(args ...interface{}) {
	l.println(LEVEL_PANIC, args...)
	panic(l.depthFunc(1))
}

// Fatalln log for fatal message
func (l *logger) Fatalln(args ...interface{}) {
	l.println(LEVEL_FATAL, args...)
}

// Trace log for debug message
func (l *logger) Trace(args ...interface{}) {
	l.print(LEVEL_TRACE, args...)
}

// Info log for info message
func (l *logger) Info(args ...interface{}) {
	l.print(LEVEL_INFO, args...)
}

// Warn log for warning message
func (l *logger) Warn(args ...interface{}) {
	l.print(LEVEL_WARN, args...)
}

// Panic log for error message
func (l *logger) Panic(args ...interface{}) {
	l.print(LEVEL_PANIC, args...)
	panic(l.depthFunc(1))
}

// Fatal log for error message
func (l *logger) Fatal(args ...interface{}) {
	l.print(LEVEL_FATAL, args...)
}

// Trace log for debug message
func (l *logger) TraceDepth(depth int, args ...interface{}) {
	l.printDepth(LEVEL_TRACE, depth+1, args...)
}

// Info log for info message
func (l *logger) InfoDepth(depth int, args ...interface{}) {
	l.printDepth(LEVEL_INFO, depth+1, args...)
}

// Warn log for warning message
func (l *logger) WarnDepth(depth int, args ...interface{}) {
	l.printDepth(LEVEL_WARN, depth+1, args...)
}

// Panic log for error message
func (l *logger) PanicDepth(depth int, args ...interface{}) {
	l.printDepth(LEVEL_PANIC, depth+1, args...)
	panic(l.depthFunc(depth + 1))
}

// Fatal log for error message
func (l *logger) FatalDepth(depth int, args ...interface{}) {
	l.printDepth(LEVEL_FATAL, depth+1, args...)
}
