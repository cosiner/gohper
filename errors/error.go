package errors

import (
	"fmt"
	"io"
	"os"
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

// Panic on error
func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

// Print error message to stderr when err is not nil
func Print(err error) {
	Fprint(os.Stderr, err)
}

// Println error message to stderr when err is not nil
func Println(err error) {
	Fprintln(os.Stderr, err)
}

// Fprint error message to given writer when err is not nil
func Fprint(w io.Writer, err error) {
	if err != nil {
		fmt.Fprint(w, err)
	}
}

// Fprintln error message to given writer when err is not nil
func Fprintln(w io.Writer, err error) {
	if err != nil {
		fmt.Fprint(w, err)
	}
}

// Exit if error is not nil, print error message to Stdout and exit code is 0
func Exit(err error) {
	Fexit(os.Stdout, err)
}

// Exitln if error is not nil, print error message to Stdout and exit code is 0
func Exitln(err error) {
	Fexitln(os.Stdout, err)
}

// Fexit if error is not nil, print error message to Writer and exit code is 0
func Fexit(w io.Writer, err error) {
	if err != nil {
		fmt.Fprint(w, err)
		os.Exit(0)
	}
}

// Fexitln if error is not nil, print error message to Writer and exit code is 0
func Fexitln(w io.Writer, err error) {
	if err != nil {
		fmt.Fprintln(w, err)
		os.Exit(0)
	}
}

// Fatal print message, then Fatal with -1 as error code when err is not nil
func Fatal(err error) {
	Ffatal(os.Stderr, err)
}

// Fatalln print message, then Fatal with  -1 as error code when err is not nil
func Fatalln(err error) {
	Ffatalln(os.Stderr, err)
}

// Ffatal print message, then exit with  -1 as error code when err is not nil
func Ffatal(w io.Writer, err error) {
	if err != nil {
		fmt.Fprint(w, err)
		os.Exit(-1)
	}
}

// Ffatalln print message, then exit with  -1 as error code when err is not nil
func Ffatalln(w io.Writer, err error) {
	if err != nil {
		fmt.Fprintln(w, err)
		os.Exit(-1)
	}
}

// Do call param function when err is not null when err is not nil
func Do(err error, fn func(err error)) {
	if err != nil {
		fn(err)
	}
}

// CondDo call function if val is true
func CondDo(val bool, err error, fn func(error)) {
	if val {
		fn(err)
	}
}
