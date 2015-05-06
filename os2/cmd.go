package os2

import (
	"io"
	"os"
	"os/exec"
)

// RunCmd run commond, if reader or writer is nil, use os.Stdin/out/err,
// first elelemts of args should be the command to execute
func RunCmd(in io.Reader, out, err io.Writer, args ...string) error {
	if in == nil {
		in = os.Stdin
	}
	if out == nil {
		out = os.Stdout
	}
	if err == nil {
		err = os.Stderr
	}

	var c *exec.Cmd
	if l := len(args); l == 0 {
		panic("no command to run")
	} else if l == 1 {
		c = exec.Command(args[0])
	} else {
		c = exec.Command(args[0], args[1:]...)
	}
	c.Stdin = in
	c.Stdout = out
	c.Stderr = err

	return c.Run()
}
