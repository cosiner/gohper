package bytes2

import (
	"bytes"
)

func NewBuffer(size int) *bytes.Buffer {
	return bytes.NewBuffer(make([]byte, 0, size))
}
