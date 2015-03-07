package log

import (
	"testing"

	e "github.com/cosiner/golib/errors"
)

func TestConsoleLog(t *testing.T) {
	Init(DEF_FLUSHINTERVAL, LEVEL_DEBUG)
	AddConsoleWriter("info=green")
	Start()
	Errorln("aaa1")
	Debugln("aaa2")
	Errorln("aaa3")
	Infoln("aaa4")
	Warnln("aaa4")
	Errorln("aaa4")
	// Fatalln("dddddddddd")
	// Fatalln("ddddddddddddddddddaaaaaaaaaaaa")
}

func TestFileLog(t *testing.T) {
	Init(DEF_FLUSHINTERVAL, LEVEL_WARN)
	e.OnErrExit(AddFileWriter("bufsize=10240&maxsize=10240&logdir=/tmp/logs"))
	Start()
	Warnln("DDDDDDDDDDDDDDDD")
	Flush()
	Exit()
}
