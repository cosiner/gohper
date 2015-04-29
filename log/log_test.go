package log

import (
	"testing"
	"time"

	"github.com/cosiner/gohper/lib/errors"
)

func TestConsoleLog(t *testing.T) {
	logger := New(nil)
	errors.OnErrExit(logger.AddWriter(new(ConsoleWriter), &ConsoleWriterOption{
		Colors: map[string]string{},
	}))
	logger.Infoln("aaa1")
	logger.Traceln("aaa2")
	logger.Infoln("aaa3")
	logger.Infoln("aaa4")
	logger.Warnln("aaa4")
	logger.Infoln("aaa4")
}

func TestFileLog(t *testing.T) {
	logger := New(&LoggerOption{
		Level: LEVEL_TRACE,
	})
	logger.AddWriter(new(ConsoleWriter), nil)
	errors.OnErrExit(logger.AddWriter(new(FileWriter), &FileWriterOption{
		Bufsize: "10K",
		Maxsize: "10M",
		Logdir:  "logss",
		Daily:   true,
	}))
	logger.Warnln("DDDDDDDDDDDDDDDD")
	logger.Infoln("DDDDDDDDDDDDDDDD")
	logger.Traceln("DDDDDDDDDDDDDDDD")
	logger.Close()
	time.Sleep(100 * time.Millisecond)
}
