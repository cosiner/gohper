package trace

import (
	"fmt"

	"github.com/cosiner/gohper/runtime2"
)

var TraceEnabled = true

type traceError struct {
	pos string
	err error
}

func (e traceError) Pos() string {
	return e.pos
}
func (e traceError) Error() string {
	return e.pos + ":" + e.err.Error()
}

func (e traceError) Unwrap() error {
	return e.err
}

func TraceDepth(err error, depth int) error {
	if err == nil {
		return nil
	}
	if !TraceEnabled {
		fmt.Println(err)
		return err
	}

	if _, is := err.(traceError); is {
		return err
	}
	return traceError{
		pos: runtime2.Caller(depth + 1),
		err: err,
	}
}

func Trace(err error) error {
	return TraceDepth(err, 1)
}
