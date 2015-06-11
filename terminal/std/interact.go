package std

import (
	"fmt"
	"os"
)

const READ_BUFSIZE = 256

type Interactor struct {
	Error error
	Buf   []byte
}

func (i *Interactor) ReadInput(prompt, def string) string {
	if i.Error != nil {
		return ""
	}

	if i.Buf == nil {
		i.Buf = make([]byte, READ_BUFSIZE)
	}

	_, i.Error = fmt.Print(prompt)
	if i.Error != nil {
		return ""
	}

	_ = os.Stdout.Sync()

	var n int
	n, i.Error = os.Stdin.Read(i.Buf)
	if i.Error != nil {
		return ""
	}

	if n <= 1 {
		return def
	}

	return string(i.Buf[:n-1]) // remove '\n'
}
