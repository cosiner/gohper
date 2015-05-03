package errors

import (
	"fmt"
	"io"
	"os"
)

const (
	// ExitCode is the default exit code on error
	ExitCode = -1
)

type Err string

func (e Err) Error() string {
	return string(e)
}

// Assert assert val is true, else panic error
func Assert(val bool, err interface{}) {
	if !val {
		panic(err)
	}
}

// New make an error with given params
func New(v ...interface{}) error {
	return Err(fmt.Sprint(v...))
}

// Newln make an error with given params and append an newline character
func Newln(v ...interface{}) error {
	return Err(fmt.Sprintln(v...))
}

// Newf make an error with given format and  params
func Newf(format string, v ...interface{}) error {
	return Err(fmt.Sprintf(format, v...))
}

// Panic panic on error
func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

// Print error message to stderr when err is not nil
func Print(err error) {
	Fprint(os.Stderr, err)
}

// Fprint error message to given writer when err is not nil
func Fprint(w io.Writer, err error) {
	if err != nil {
		fmt.Fprint(w, err)
	}
}

// Fatal print message, then Fatal with error code when err is not nil
func Fatal(err error) {
	Ffatal(os.Stderr, err)
}

// Ffatal print message, then exit with error code when err is not nil
func Ffatal(w io.Writer, err error) {
	if err != nil {
		Fprint(w, err)
		os.Exit(ExitCode)
	}
}

// Do call param function when err is not null when err is not nil
func Do(err error, fn func(err error)) {
	if err != nil {
		fn(err)
	}
}
