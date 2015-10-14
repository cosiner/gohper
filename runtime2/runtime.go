package runtime2

import (
	"path/filepath"
	"runtime"
	"strconv"
)

type Pos struct {
	File string
	Pc   uintptr
	Line int
}

func (c Pos) String() string {
	return filepath.Base(c.File) + ":" + strconv.Itoa(c.Line)
}

func CallerPos(depth int) Pos {
	pc, file, line, _ := runtime.Caller(depth + 1)
	return Pos{
		File: file,
		Pc:   pc,
		Line: line,
	}
}

// Caller report caller's position with file:function:line format
// depth means which caller, 0 means yourself, 1 means your caller
func Caller(depth int) string {
	return CallerPos(depth + 1).String()
}

func Stack(bufsize int, all bool) []byte {
	buf := make([]byte, 0, bufsize)
	n := runtime.Stack(buf, all)
	return buf[:n]
}
