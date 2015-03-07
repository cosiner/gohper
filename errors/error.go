package errors

import (
	"fmt"
	"io"
	"os"
)

//==============================================================================
//                           Assert
//==============================================================================

// Assert assert val is true, else panic error
func Assert(val bool, err string) {
	if !val {
		panic(err)
	}
}

// AssertNoErr assert error is nil, else panic
func AssertNoErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Assertf assert val is true, else panic error
func Assertf(val bool, errformat string, v ...interface{}) {
	if !val {
		panic(fmt.Sprintf(errformat, v...))
	}
}

//==============================================================================
//                         Error Format
//==============================================================================

// ERR is a nil error for those condition only need an error
var ERR = Err("")

// errStr is a error implementation
type errStr struct {
	s string
}

// Error implements builtin error
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

// OnErrDo call param function when err is not null
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

// PrintErr output error message to stderr
func PrintErr(err error) {
	fmt.Fprint(os.Stderr, err)
}

// PrintErrln output error message to stderr
func PrintErrln(err error) {
	fmt.Fprintln(os.Stderr, err)
}

// FprintErr output error message to given writer
func FprintErr(w io.Writer, err error) {
	fmt.Fprint(w, err)
}

// FprintErrln output error message to given writer
func FprintErrln(w io.Writer, err error) {
	fmt.Fprintln(w, err)
}

// PrintError output error message to stderr
func PrintError(v ...interface{}) {
	fmt.Fprint(os.Stderr, v...)
}

// PrintErrorln output error message to stderr
func PrintErrorln(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
}

// PrintErrorf format and output error message to stderr
func PrintErrorf(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format, v...)
}

// FprintError output error message to given writer
func FprintError(w io.Writer, v ...interface{}) {
	fmt.Fprint(w, v...)
}

// FprintErrorln output error message to given writer
func FprintErrorln(w io.Writer, v ...interface{}) {
	fmt.Fprintln(w, v...)
}

// FprintErrorf output error message to given writer
func FprintErrorf(w io.Writer, format string, v ...interface{}) {
	fmt.Fprintf(w, format, v...)
}

//==============================================================================
//                         Error Exit
//==============================================================================

const (
	// ErrorExitCode is the default exit code on error
	ErrorExitCode = 1
	// NormalExitCode is the default exit code on normal
	NormalExitCode = 0
)

// ExitWith print message, then normal exit
func ExitWith(str string) {
	fmt.Println(str)
	os.Exit(NormalExitCode)
}

// FexitWith print message, then normal exit
func FexitWith(w io.Writer, str string) {
	fmt.Fprintln(w, str)
	os.Exit(NormalExitCode)
}

// ExitOnErrPrint exit process, if error, print it
func ExitOnErrPrint(err error) {
	OnErrDo(err, PrintErr)
	os.Exit(ErrorExitCode)
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
	os.Exit(ErrorExitCode)
}

// FexitErrorln print error message to writer then exit process
func FexitErrorln(w io.Writer, v ...interface{}) {
	FprintErrorln(w, v...)
	os.Exit(ErrorExitCode)
}

// FexitErrorf  format and output error message to writer, then exit process
func FexitErrorf(w io.Writer, format string, v ...interface{}) {
	FprintErrorf(w, format, v...)
	os.Exit(ErrorExitCode)
}
