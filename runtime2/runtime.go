package runtime2

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

type Pos struct {
	File string
	Pc   uintptr
	Line int
}

func pathSepFunc(r rune) bool {
	return r == '/' || r == os.PathSeparator
}

// Caller report caller's position with file:function:line format
// depth means which caller, 0 means yourself, 1 means your caller
func Caller(depth int) string {
	_, file, line, _ := runtime.Caller(depth + 1)
	i := strings.LastIndexFunc(file, pathSepFunc)
	if i >= 0 {
		j := strings.LastIndexFunc(file[:i], pathSepFunc)
		if j >= 0 {
			i = j
		}
		file = file[i+1:]
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func Stack(bufsize int, all bool) []byte {
	buf := make([]byte, bufsize)
	n := runtime.Stack(buf, all)
	return buf[:n]
}

func Recover(bufsize int, logger func(...interface{})) {
	if err := recover(); err != nil {
		logger(err, string(Stack(bufsize, false)))
	}
}

func RecoverRun(bufsize int, fn func(), logger func(...interface{})) {
	defer Recover(bufsize, logger)
	fn()
}
