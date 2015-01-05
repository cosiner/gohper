package test

import (
	"mlib/log"
	"mlib/log/test/test"
	"testing"
)

func TestLog(t *testing.T) {
	if err := log.Init(30, log.LEVEL_INFO, 1024*10, 1024*1024*100, "/tmp/logs"); err != nil {
		panic(err)
	}
	log.SetLevel(log.LEVEL_WARN)
	test.WarnLog()
	log.Error("aaa")
	log.Error("aaa")
	log.Error("aaa")
	log.Error("aaa")
	log.Error("aaa")
	log.Error("aaa")
	log.Error("aaa")
	log.Error("aaa")
	log.Error("aaa")
	log.Error("aaa")
	log.Error("aaa")
	log.Error("aaa")
	log.Error("aaa")
	log.FlushAndExit()
}
