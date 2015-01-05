package sys

import (
	"runtime"
	"testing"
)

func TestHomeDir(t *testing.T) {

}

func TestExpandHome(t *testing.T) {

}

func TestExpandAbs(t *testing.T) {

}

func TestProgramDir(t *testing.T) {

}

func TestLastDir(t *testing.T) {

}

func TestMkdirWithParent(t *testing.T) {
	t.Log(MkdirWithParent(ExpandAbs("~/test/test")))
}

func TestArch(t *testing.T) {
	t.Log(runtime.GOOS)
}
