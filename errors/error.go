package errors

import (
	"errors"
	"fmt"
	"io"
	"os"
)

const (
	// ERREXIT_CODE is the default exit code on error
	ERREXIT_CODE = 1
	// NOREXIT_CODE is the default exit code on normal
	NOREXIT_CODE = 0
)

// Errorf make an error with given format and  params
func Errorf(format string, v ...interface{}) error {
	return errors.New(fmt.Sprintf(format, v...))
}

// Errorln make an error with given params and append an newline character
func Errorln(v ...interface{}) error {
	return errors.New(fmt.Sprintln(v...))
}

// Error make an error with given params
func Error(v ...interface{}) error {
	return errors.New(fmt.Sprint(v...))
}

// Err is only a wrapper of errors.New
func Err(str string) error {
	return errors.New(str)
}

// ErrorPanic call f and panic on error
func ErrorPanic(f func() error) {
	OnError(f(), func(err error) {
		panic(err)
	})
}

// ErrorStrPanic call f and panic on error print given error message
func ErrorStrPanic(f func() error, errStr string) {
	OnError(f(), func(err error) {
		panic(errStr)
	})
}

// OnError call param function when err is not null
func OnError(err error, fn func(err error)) {
	if err != nil {
		fn(err)
	}
}

// ConsoleError output error message to stderr
func ConsoleError(v ...interface{}) {
	fmt.Fprint(os.Stderr, v...)
}

// ConsoleErrorf format and output error message to stderr
func ConsoleErrorf(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format, v...)
}

// ErrorExit print error message to stderr then exit process
func ErrorExit(v ...interface{}) {
	FerrorExit(os.Stderr, v...)
}

// ErrorfExit  format and output error message to stderr, then exit process
func ErrorfExit(format string, v ...interface{}) {
	FerrorfExit(os.Stderr, format, v...)
}

// FerrorExit print error message to writer then exit process
func FerrorExit(w io.Writer, v ...interface{}) {
	fmt.Fprint(w, v...)
	os.Exit(ERREXIT_CODE)
}

// FerrorfExit  format and output error message to writer, then exit process
func FerrorfExit(w io.Writer, format string, v ...interface{}) {
	fmt.Fprintf(w, format, v...)
	os.Exit(ERREXIT_CODE)
}
