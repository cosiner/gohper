// Package sys implement some utilities for common operations
package sys

import (
	"io"
	"os"
	"os/exec"
)

// RunCmdN run commond with N io parameter specified stdin, stdout, stderr

func RunCmd0(cmd string, arg ...string) error {
	return RunCmd3(os.Stdin, os.Stdout, os.Stderr, cmd, arg...)
}

func RunCmd1(in io.Reader, cmd string, arg ...string) error {
	return RunCmd3(in, os.Stdout, os.Stderr, cmd, arg...)
}

func RunCmd2(in io.Reader, out io.Writer, cmd string, arg ...string) error {
	return RunCmd3(in, out, os.Stderr, cmd, arg...)
}

func RunCmd3(in io.Reader, out io.Writer, err io.Writer, cmd string, arg ...string) error {
	c := exec.Command(cmd, arg...)
	c.Stdin = in
	c.Stdout = out
	c.Stderr = err
	return c.Run()
}
