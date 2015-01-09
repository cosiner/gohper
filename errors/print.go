package errors

import (
	"fmt"
	"io"
	"os"
)

// ErrPrint output error message to stderr
func PrintErr(err error) {
	fmt.Fprint(os.Stderr, err)
}

// ErrPrintln output error message to stderr
func PrintErrln(err error) {
	fmt.Fprintln(os.Stderr, err)
}

// ErrFprint output error message to stderr
func FprintErr(w io.Writer, err error) {
	fmt.Fprint(w, err)
}

// ErrFprintln output error message to stderr
func FprintErrln(w io.Writer, err error) {
	fmt.Fprintln(w, err)
}

// ErrorPrint output error message to stderr
func PrintError(v ...interface{}) {
	fmt.Fprint(os.Stderr, v...)
}

// ErrorPrintln output error message to stderr
func PrintErrorln(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
}

// ErrorPrintf format and output error message to stderr
func PrintErrorf(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format, v...)
}

// ErrorFprint output error message to stderr
func FprintError(w io.Writer, v ...interface{}) {
	fmt.Fprint(w, v...)
}

// ErrorFprintln output error message to stderr
func FprintErrorln(w io.Writer, v ...interface{}) {
	fmt.Fprintln(w, v...)
}

// ErrorFprintf output error message to stderr
func FprintErrorf(w io.Writer, format string, v ...interface{}) {
	fmt.Fprintf(w, format, v...)
}
