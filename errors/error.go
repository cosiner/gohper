package errors

import (
	"fmt"
	"io"
	"os"
)

//==============================================================================
//                           Assert
//==============================================================================
// Assert val is true, else panic error
func Assert(val bool, err error) {
	if !val {
		panic(err)
	}
}

//==============================================================================
//                         Error Format
//==============================================================================
// ERR is a nil error for those condition only need an error
var ERR = Err("")

type errStr struct {
	s string
}

func (es *errStr) Error() string {
	return es.s
}

// Err wrap a string to error
func Err(str string) error {
	return &errStr{s: str}
}

// Error make an error with given params
func Error(v ...interface{}) error {
	return Err(fmt.Sprint(v...))
}

// Errorln make an error with given params and append an newline character
func Errorln(v ...interface{}) error {
	return Err(fmt.Sprintln(v...))
}

// Errorf make an error with given format and  params
func Errorf(format string, v ...interface{}) error {
	return Err(fmt.Sprintf(format, v...))
}

//==============================================================================
//                         Error Event
//==============================================================================
// OnErrExit exit process and print message on error
func OnErrExit(err error) {
	if err != nil {
		ExitErrln(err)
	}
}

// OnErrExitStr exit process and print error string on error
func OnErrExitStr(err error, errStr string) {
	if err != nil {
		ExitErrorln(errStr)
	}
}

// OnErrPanic panic on error
func OnErrPanic(err error) {
	if err != nil {
		panic(err)
	}
}

// OnErrPanicStr panic error string on error
func OnErrPanicStr(err error, errStr string) {
	if err != nil {
		panic(errStr)
	}
}

// OnErr call param function when err is not null
func OnErrDo(err error, fn func(err error)) {
	if err != nil {
		fn(err)
	}
}

// OnFuncErrDo call second funcion when first function return error
func OnFuncErrDo(f func() error, fn func(err error)) {
	OnErrDo(f(), fn)
}

// OnErrDoChain is a chain to process error
// when error is not completed process, it will finally throw again
// to stop it, only return nil in one of the process chain is needed
func OnErrDoChain(err error, fns ...func(err error) error) error {
	for i, end := 0, len(fns); err != nil && i < end; i++ {
		err = fns[i](err)
	}
	return err
}

//==============================================================================
//                         Error Print
//==============================================================================
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

//==============================================================================
//                         Error Exit
//==============================================================================

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
