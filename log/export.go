package log

//==================format======================================================
var (
	logger       *Logger
	nilFunc      = func() {}
	Start        = nilFunc
	Stop         = nilFunc
	Flush        = nilFunc
	Exit         = nilFunc
	LogLevel     = func() Level { return unknownLevel }
	SetLevel     = func(_ Level) error { return nil }
	AddLogWriter = func(_ LogWriter) error { return nil }
)

// Init init global logger
func Init(flushInterval int, level Level) {
	logger = NewLogger(flushInterval, level)
	Start = logger.Start
	Stop = logger.Stop
	Flush = logger.Flush
	Exit = logger.Exit
	LogLevel = logger.LogLevel
	SetLevel = logger.SetLevel
	AddLogWriter = logger.AddLogWriter
}

// GlobalLogger return current global logger
func GlobalLogger() *Logger {
	return logger
}

// AddConsoleWriter add an console log writer
func AddConsoleWriter() {
	logger.AddLogWriter(new(ConsoleLogWriter))
}

// Debugf log for debug message
func Debugf(format string, v ...interface{}) {
	logger.logf(LEVEL_DEBUG, format, v...)
}

// Infof log for info message
func Infof(format string, v ...interface{}) {
	logger.logf(LEVEL_INFO, format, v...)
}

// Warnf log for warning message
func Warnf(format string, v ...interface{}) {
	logger.logf(LEVEL_WARN, format, v...)
}

// Errorf log for error message
func Errorf(format string, v ...interface{}) {
	logger.logf(LEVEL_ERROR, format, v...)
}

// Fatalf log for fatal message
func Fatalf(format string, v ...interface{}) {
	logger.logf(LEVEL_FATAL, format, v...)
	logger.Exit()
}

// Debugln log for debug message
func Debugln(v ...interface{}) {
	logger.logln(LEVEL_DEBUG, v...)
}

// Infoln log for info message
func Infoln(v ...interface{}) {
	logger.logln(LEVEL_INFO, v...)
}

// Warnln log for warning message
func Warnln(v ...interface{}) {
	logger.logln(LEVEL_WARN, v...)
}

// Errorln log for error message
func Errorln(v ...interface{}) {
	logger.logln(LEVEL_ERROR, v...)
}

// Fatalln log for fatal message
func Fatalln(v ...interface{}) {
	logger.logln(LEVEL_FATAL, v...)
	logger.Exit()
}

// Debug log for debug message
func Debug(v ...interface{}) {
	logger.log(LEVEL_DEBUG, v...)
}

// Info log for info message
func Info(v ...interface{}) {
	logger.log(LEVEL_INFO, v...)
}

// Warn log for warning message
func Warn(v ...interface{}) {
	logger.log(LEVEL_WARN, v...)
}

// Error log for error message
func Error(v ...interface{}) {
	logger.log(LEVEL_ERROR, v...)
}

// Fatalflog for fatal message
func Fatal(v ...interface{}) {
	logger.log(LEVEL_FATAL, v...)
	logger.Exit()
}
