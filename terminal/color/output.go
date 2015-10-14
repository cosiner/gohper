package color

import (
	"io"
	"os"
	"runtime"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

var (
	Stdout io.Writer
	Stderr io.Writer
	IsTTY  bool
)

func init() {
	IsTTY = isatty.IsTerminal(os.Stdout.Fd())

	if runtime.GOOS == "windows" {
		Stdout = colorable.NewColorableStdout()
		Stderr = colorable.NewColorableStderr()
	} else {
		Stdout = os.Stdout
		Stderr = os.Stderr
	}
}
