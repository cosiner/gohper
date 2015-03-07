package log

var (
	l logger = emptyLogger{}
)

// Init init global l
func Init(flushInterval int, level Level) {
	l = NewLogger(flushInterval, level)
}

// AddConsoleWriter add an console log writer
func AddConsoleWriter(conf string) error {
	clw := new(ConsoleLogWriter)
	if err := clw.Config(conf); err != nil {
		return err
	}
	return l.AddLogWriter(clw)
}

// AddFileWriter add file log writer to l
func AddFileWriter(conf string) error {
	flw := new(FileLogWriter)
	if err := flw.Config(conf); err != nil {
		return err
	}
	return l.AddLogWriter(flw)
}

func AddLogWriter(writer LogWriter) error { return l.AddLogWriter(writer) }
func LogLevel() Level                     { return l.LogLevel() }
func SetLevel(level Level) error          { return l.SetLevel(level) }
func Start()                              { l.Start() }
func Flush()                              { l.Flush() }
func Exit()                               { l.Exit() }
func Stop()                               { l.Stop() }

// Debugf log for debug message
func Debugf(format string, v ...interface{}) {
	l.logf(LEVEL_DEBUG, format, v...)
}

// Infof log for info message
func Infof(format string, v ...interface{}) {
	l.logf(LEVEL_INFO, format, v...)
}

// Warnf log for warning message
func Warnf(format string, v ...interface{}) {
	l.logf(LEVEL_WARN, format, v...)
}

// Errorf log for error message
func Errorf(format string, v ...interface{}) {
	l.logf(LEVEL_ERROR, format, v...)
}

// Fatalf log for fatal message
func Fatalf(format string, v ...interface{}) {
	l.logf(LEVEL_FATAL, format, v...)
	l.Exit()
}

// Debugln log for debug message
func Debugln(v ...interface{}) {
	l.logln(LEVEL_DEBUG, v...)
}

// Infoln log for info message
func Infoln(v ...interface{}) {
	l.logln(LEVEL_INFO, v...)
}

// Warnln log for warning message
func Warnln(v ...interface{}) {
	l.logln(LEVEL_WARN, v...)
}

// Errorln log for error message
func Errorln(v ...interface{}) {
	l.logln(LEVEL_ERROR, v...)
}

// Fatalln log for fatal message
func Fatalln(v ...interface{}) {
	l.logln(LEVEL_FATAL, v...)
	l.Exit()
}

// Debug log for debug message
func Debug(v ...interface{}) {
	l.log(LEVEL_DEBUG, v...)
}

// Info log for info message
func Info(v ...interface{}) {
	l.log(LEVEL_INFO, v...)
}

// Warn log for warning message
func Warn(v ...interface{}) {
	l.log(LEVEL_WARN, v...)
}

// Error log for error message
func Error(v ...interface{}) {
	l.log(LEVEL_ERROR, v...)
}

// Fatalflog for fatal message
func Fatal(v ...interface{}) {
	l.log(LEVEL_FATAL, v...)
	l.Exit()
}
