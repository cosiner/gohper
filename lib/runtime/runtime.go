package runtime

import (
	"path/filepath"
	"runtime"
	"strconv"
)

// CallerPosition report caller's position with file:function:line format
// depth means which caller, 0 means yourself, 1 means your caller
func CallerPosition(depth int) string {
	pc, file, line, _ := runtime.Caller(depth + 1)
	return filepath.Base(file) + ":" + filepath.Base(runtime.FuncForPC(pc).Name()) + ":" + strconv.Itoa(line)
}
