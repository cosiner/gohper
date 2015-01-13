package log

import (
	"testing"

	e "github.com/cosiner/golib/errors"
)

func testConsoleLog(t *testing.T) {
	Init(DEF_FLUSHINTERVAL, LEVEL_WARN)
	AddConsoleWriter()
	Start()
	Errorln("aaa1")
	Debugln("aaa2")
	Errorln("aaa3")
	Infoln("aaa4")
	Warnln("aaa4")
	Errorln("aaa4")
	Flush()
	Exit()
	Errorln("aaadddd")
}

func TestFileLog(t *testing.T) {
	Init(DEF_FLUSHINTERVAL, LEVEL_WARN)
	fwr := new(FileLogWriter)
	e.OnErrExit(fwr.Config("bufsize=10240&maxsize=10240&logdir=/tmp/logs"))
	e.OnErrExit(AddLogWriter(fwr))
	AddConsoleWriter()
	Start()
	Warnln("DDDDDDDDDDDDDDDD")
	Flush()
	Exit()
}
