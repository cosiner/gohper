package errors

import "fmt"

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
