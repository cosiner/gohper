package errors

import (
	"fmt"
	"os"
)

type Err string

func (e Err) Error() string {
	return string(e)
}

func Assert(val bool, err error) {
	if !val {
		panic(err)
	}
}

func Exclude(err error, errors ...error) error {
	for _, e := range errors {
		if err == e {
			return nil
		}
	}
	return err
}

func New(v ...interface{}) error {
	return Err(fmt.Sprint(v...))
}

func Newln(v ...interface{}) error {
	return Err(fmt.Sprintln(v...))
}

func Newf(format string, v ...interface{}) error {
	return Err(fmt.Sprintf(format, v...))
}

func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

func Panicln(err error) {
	if err != nil {
		panic(fmt.Sprintln(err))
	}
}

func Panicf(format string, err error) {
	if err != nil {
		panic(fmt.Sprintf(format, err))
	}
}

func Print(err error) {
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
}

func Println(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func Printf(format string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, format, err)
	}
}

func Exit(err error) {
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(0)
	}
}

func Exitln(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(0)
	}
}

func Exitf(format string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, format, err)
		os.Exit(0)
	}
}

func Fatal(err error) {
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(-1)
	}
}

func Fatalln(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

func Fatalf(format string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, format, err)
		os.Exit(-1)
	}
}

func ExitAny(err error) {
	fmt.Fprint(os.Stderr, err)
	os.Exit(0)
}

func ExitAnyf(format string, err error) {
	fmt.Fprintln(os.Stderr, format, err)
	os.Exit(0)
}

func ExitAnyln(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(0)
}

func FatalAny(err error) {
	fmt.Fprint(os.Stderr, err)
	os.Exit(-1)
}

func FatalAnyf(format string, err error) {
	fmt.Fprint(os.Stderr, format, err)
	os.Exit(-1)
}

func FatalAnyln(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(-1)
}

func Do(err error, fn func(err error)) {
	if err != nil {
		fn(err)
	}
}

func CondDo(val bool, err error, fn func(error)) {
	if val {
		fn(err)
	}
}

func Nil(err error, errs ...error) error {
	for _, e := range errs {
		if err == e {
			return nil
		}
	}
	return err
}
