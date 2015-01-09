package errors

import (
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

// ExitWith print message, then normal exit
func ExitWith(str string) {
	fmt.Println(str)
	os.Exit(NOREXIT_CODE)
}

// FexitWith print message, then normal exit
func FexitWith(w io.Writer, str string) {
	fmt.Fprintln(w, str)
	os.Exit(NOREXIT_CODE)
}

// ExitOnErrPrint exit process, if error, print it
func ExitOnErrPrint(err error) {
	OnErrDo(err, PrintErr)
	os.Exit(ERREXIT_CODE)
}

// ExitErr print error message then exit process
func ExitErr(err error) {
	FexitError(os.Stderr, err)
}

// ExitErrln print error message then exit process
func ExitErrln(err error) {
	FexitErrorln(os.Stderr, err)
}

// FexitErr print error message to writer then exit process
func FexitErr(w io.Writer, err error) {
	FexitError(w, err)
}

// FexitErrln print error message to writer then exit process
func FexitErrln(w io.Writer, err error) {
	FexitErrorln(w, err)
}

// ExitError print error message to stderr then exit process
func ExitError(v ...interface{}) {
	FexitError(os.Stderr, v...)
}

// ExitErrorln print error message to stderr then exit process
func ExitErrorln(v ...interface{}) {
	FexitErrorln(os.Stderr, v...)
}

// ExitErrorf  format and output error message to stderr, then exit process
func ExitErrorf(format string, v ...interface{}) {
	FexitErrorf(os.Stderr, format, v...)
}

// FexitError print error message to writer then exit process
func FexitError(w io.Writer, v ...interface{}) {
	FprintError(w, v...)
	os.Exit(ERREXIT_CODE)
}

// FexitErrorln print error message to writer then exit process
func FexitErrorln(w io.Writer, v ...interface{}) {
	FprintErrorln(w, v...)
	os.Exit(ERREXIT_CODE)
}

// FexitErrorf  format and output error message to writer, then exit process
func FexitErrorf(w io.Writer, format string, v ...interface{}) {
	FprintErrorf(w, format, v...)
	os.Exit(ERREXIT_CODE)
}
