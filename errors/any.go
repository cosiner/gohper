package errors

import (
	"fmt"
	"io"
	"os"
)

// ExitAny print error message to Stdout and exit code is 0
func ExitAny(err error) {
	Fexit(os.Stdout, err)
}

// ExitAnyln print error message to Stdout and exit code is 0
func ExitAnyln(err error) {
	Fexitln(os.Stdout, err)
}

// FexitAny print error message to Writer and exit code is 0
func FexitAny(w io.Writer, err error) {
	fmt.Fprint(w, err)
	os.Exit(0)
}

// FexitAnyln print error message to Writer and exit code is 0
func FexitAnyln(w io.Writer, err error) {
	fmt.Fprintln(w, err)
	os.Exit(0)
}

// FatalAny print message, then Fatal with -1 as error code
func FatalAny(err error) {
	Ffatal(os.Stderr, err)
}

// FatalAnyln print message, then Fatal with  -1 as error code
func FatalAnyln(err error) {
	Ffatalln(os.Stderr, err)
}

// FfatalAny print message, then exit with  -1 as error code when err is not nil
func FfatalAny(w io.Writer, err error) {
	fmt.Fprint(w, err)
	os.Exit(-1)
}

// FfatalAnyln print message, then exit with  -1 as error code when err is not nil
func FfatalAnyln(w io.Writer, err error) {
	fmt.Fprintln(w, err)
	os.Exit(-1)
}
