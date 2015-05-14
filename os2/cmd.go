package os2

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// StdRunCmd run commond, if reader or writer is nil, use os.Stdin/out/err,
func StdRunCmd(in io.Reader, out, err io.Writer, args ...string) error {
	if in == nil {
		in = os.Stdin
	}
	if out == nil {
		out = os.Stdout
	}
	if err == nil {
		err = os.Stderr
	}

	return RunCmd(in, out, err, args...)
}

// RunCmd is a quick way to exec.Command().Run,
// first elelemts of args should be the command to execute
func RunCmd(in io.Reader, out, err io.Writer, args ...string) error {
	if len(args) == 0 {
		panic("no command to run")
	}

	name := args[0]
	if filepath.Base(name) == name {
		if lp, err := exec.LookPath(name); err != nil {
			return err
		} else {
			name = lp
		}
	}

	return (&exec.Cmd{
		Path:   name,
		Args:   args,
		Stdin:  in,
		Stdout: out,
		Stderr: err,
	}).Run()
}
