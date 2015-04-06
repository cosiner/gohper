package log

import (
	e "github.com/cosiner/gohper/lib/errors"

	"testing"
)

func TestConsoleLog(t *testing.T) {
	logger := New(DEF_FLUSHINTERVAL, LEVEL_INFO)
	logger.AddConfWriter(new(ConsoleLogWriter), "info=green")
	logger.Start()
	logger.Errorln("aaa1")
	logger.Debugln("aaa2")
	logger.Errorln("aaa3")
	logger.Infoln("aaa4")
	logger.Warnln("aaa4")
	logger.Errorln("aaa4")
	// Fatalln("dddddddddd")
	// Fatalln("ddddddddddddddddddaaaaaaaaaaaa")
}

func TestFileLog(t *testing.T) {
	logger := New(DEF_FLUSHINTERVAL, LEVEL_INFO)
	e.OnErrExit(logger.AddConfWriter(new(FileLogWriter), "bufsize=10240&maxsize=10240&logdir=/tmp/logs&level=info"))
	logger.Start()
	logger.Warnln("DDDDDDDDDDDDDDDD")
	logger.Warnln("DDDDDDDDDDDDDDDD")
	logger.Warnln("DDDDDDDDDDDDDDDD")
	logger.Warnln("DDDDDDDDDDDDDDDD")
	logger.Warnln("DDDDDDDDDDDDDDDD")
	logger.Warnln("DDDDDDDDDDDDDDDD")
	logger.Warnln("DDDDDDDDDDDDDDDD")
	logger.Warnln("DDDDDDDDDDDDDDDD")
	logger.Warnln("DDDDDDDDDDDDDDDD")
	logger.Warnln("DDDDDDDDDDDDDDDD")
	logger.Warnln("DDDDDDDDDDDDDDDD")
	logger.Warnln("DDDDDDDDDDDDDDDD")
	logger.Warnln("DDDDDDDDDDDDDDDD")
	logger.Flush()
}
