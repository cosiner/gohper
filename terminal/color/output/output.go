package output

import (
	"io"
	"os"

	"github.com/cosiner/gohper/os2"
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

	if os2.IsWindows() {
		Stdout = colorable.NewColorableStdout()
		Stderr = colorable.NewColorableStderr()
	} else {
		Stdout = os.Stdout
		Stderr = os.Stderr
	}
}
