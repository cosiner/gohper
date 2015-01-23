package log

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/cosiner/golib/termcolor"

	"github.com/cosiner/gomodule/config"

	"github.com/cosiner/golib/errors"
	t "github.com/cosiner/golib/time"
)

var (
	timenow  = time.Now
	datetime = t.DateTime
)

//==============================================================================
//                         Log
//==============================================================================

// Log represend a log with level and log message
type Log struct {
	Level   Level
	Message string
	Time    string
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
		Time:    datetime(),
	}
}

func NewLog(level Level, v ...interface{}) *Log {
	return &Log{
		Level:   level,
		Message: fmt.Sprint(v...),
		Time:    datetime(),
	}
}

func NewLogln(level Level, v ...interface{}) *Log {
	return &Log{
		Level:   level,
		Message: fmt.Sprintln(v...),
		Time:    datetime(),
	}
}

//==============================================================================
//                              LogWriter
//==============================================================================

// LogWriter is actual log writer
type LogWriter interface {
	// Config config writer
	Config(conf string) error
	// Writer output log
	Write(log *Log) error
	// Resetlevel reset log writer's level
	ResetLevel(level Level) error
	// Flush flush output
	Flush()
	// Close close log writer
	Close()
}

//==============================================================================
//                         Console Log Writer
//==============================================================================

// bgColor create color render use given background color, default highlight
func bgColor(bg string) *termcolor.TermColor {
	return termcolor.NewColor().Highlight().Bg(bg)
}

// defTermColor define default color for each log level
var defTermColor = [5]*termcolor.TermColor{
	bgColor(termcolor.GREEN),  //debug
	bgColor(termcolor.WHITE),  //info
	bgColor(termcolor.YELLOW), //warn
	bgColor(termcolor.BLUE),   //error
	bgColor(termcolor.RED),    //fatal
}

// ConsoleLogWriter output log to console
type ConsoleLogWriter struct {
	termColor [5]*termcolor.TermColor
}

// Config config console log writer
// parameter conf can use to config color for each log level, such as
// warn="black"&info="green"&error="red"...
func (clw *ConsoleLogWriter) Config(conf string) error {
	clw.termColor = defTermColor
	if conf != "" {
		c := config.NewConfig(config.LINE)
		c.ParseString(conf)
		for l := _LEVEL_MIN; l < _LEVEL_MAX; l++ {
			s := strings.ToLower(l.String())
			if color := c.ValDef(s, ""); color != "" {
				clw.termColor[l] = bgColor(color)
			}
		}
	}
	return nil
}

// DisableColor disable color output
func (clw *ConsoleLogWriter) DisableColor() {
	for _, tc := range clw.termColor {
		tc.Disable()
	}
}

// Write write
func (clw *ConsoleLogWriter) Write(log *Log) error {
	_, err := fmt.Print(clw.termColor[log.Level].Render(log.String()))
	return err
}

func (clw *ConsoleLogWriter) ResetLevel(level Level) error { return nil }
func (clw *ConsoleLogWriter) Flush()                       {}
func (clw *ConsoleLogWriter) Close()                       {}

//==============================================================================
//                Logger
//==============================================================================

type signalType uint8

const (
	_SIGNAL_FLUSH signalType = iota // flush all writer
	_SIGNAL_STOP                    // stop logger
	_SIGNAL_EXIT                    // exit process
)

// Logger
type Logger struct {
	*sync.RWMutex
	level         Level
	writers       []LogWriter
	flushInterval time.Duration
	logs          chan *Log
	signal        chan signalType
	running       bool
}

// NewLogger return a logger, panic an error on parameter
func NewLogger(flushInterval int, level Level) *Logger {
	errors.Assert(level >= _LEVEL_MIN && level <= _LEVEL_MAX,
		UnknownLevelErr(level.String()))
	errors.Assert(flushInterval > 0,
		errors.Errorf("Flush interval should not be negative:%d", flushInterval))
	return &Logger{RWMutex: new(sync.RWMutex),
		level:         level,
		logs:          make(chan *Log, DEF_BACKLOG),
		signal:        make(chan signalType),
		flushInterval: time.Duration(flushInterval) * time.Second,
		running:       false,
	}
}

// AddLogWroter add a  log writer, nil writer will be auto-ignored
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

// SetLevel change logger's level, it will apply to all log writers
func (logger *Logger) SetLevel(level Level) (err error) {
	errors.Assert(level >= _LEVEL_MIN && level <= _LEVEL_MAX,
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
				logger.RUnlock()
			case <-ticker:
				logger.Lock()
				for _, writer := range logger.writers {
					writer.Flush()
				}
				logger.Unlock()
			case signal := <-logger.signal:
				logger.Lock()
				switch signal {
				case _SIGNAL_FLUSH:
					for _, writer := range logger.writers {
						writer.Flush()
					}
				case _SIGNAL_STOP:
					logger.running = false
					logger.Unlock()
					return
				case _SIGNAL_EXIT:
					logger.running = false
					// if there remains some logs to output, then continue this loop
					// and set a alarm for later exit process
					if len(logger.logs) > 0 {
						logger.Unlock()
						time.AfterFunc(20*time.Millisecond, logger.Exit)
						continue
					}
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

// Flush flush logger
func (logger *Logger) Flush() {
	logger.signal <- _SIGNAL_FLUSH
}

// Exit exit process
func (logger *Logger) Exit() {
	logger.signal <- _SIGNAL_EXIT
}

// Stop stop logger
func (logger *Logger) Stop() {
	logger.signal <- _SIGNAL_STOP
}

//==============================================================================
//                              Output
//==============================================================================
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
