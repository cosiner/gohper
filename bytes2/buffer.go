package bytes2

import (
	"encoding/binary"
	"errors"
	"io"
	"unicode/utf8"

	"github.com/cosiner/gohper/unsafe2"
)

type Buffer struct {
	Buf     []byte
	readPos int
}

func NewBuffer(data []byte) *Buffer {
	return &Buffer{Buf: data}
}

func MakeBuffer(size, capacity int) *Buffer {
	return NewBuffer(make([]byte, size, capacity))
}

func (b *Buffer) Grows(n int) (i int) {
	i = len(b.Buf)

	newLen := len(b.Buf) + n
	if cap(b.Buf) >= newLen {
		b.Buf = b.Buf[:newLen]
		return
	}

	data := make([]byte, newLen, cap(b.Buf)/4+newLen)
	copy(data, b.Buf)
	b.Buf = data
	return
}

func (b *Buffer) Truncate(size int) {
	if len(b.Buf) > size {
		b.Buf = b.Buf[:size]
		if b.readPos > size {
			b.readPos = size
		}
	}
}

func (b *Buffer) ResetUndelay(data []byte) {
	b.Buf = data
	b.readPos = 0
}

func (b *Buffer) Reset() {
	b.Truncate(0)
}

func (b *Buffer) Bytes() []byte {
	return b.Buf[b.readPos:]
}

func (b *Buffer) Len() int {
	return len(b.Buf) - b.readPos
}

func (b *Buffer) Cap() int {
	return cap(b.Buf)
}

func (b *Buffer) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	if b.readPos >= len(b.Buf) {
		return 0, io.EOF
	}
	n := copy(p, b.Buf[b.readPos:])
	b.readPos += n
	return n, nil
}

func (b *Buffer) ReadByte() (byte, error) {
	if b.readPos >= len(b.Buf) {
		return 0, io.EOF
	}
	r := b.Buf[b.readPos]
	b.readPos += 1
	return r, nil
}

func (b *Buffer) ReadUint16(order binary.ByteOrder) (uint16, error) {
	if b.readPos >= len(b.Buf)-2 {
		return 0, io.EOF
	}
	u := order.Uint16(b.Buf[b.readPos:])
	b.readPos += 2
	return u, nil
}

func (b *Buffer) ReadUint32(order binary.ByteOrder) (uint32, error) {
	if b.readPos >= len(b.Buf)-4 {
		return 0, io.EOF
	}
	u := order.Uint32(b.Buf[b.readPos:])
	b.readPos += 4
	return u, nil
}

func (b *Buffer) ReadUint64(order binary.ByteOrder) (uint64, error) {
	if b.readPos >= len(b.Buf)-8 {
		return 0, io.EOF
	}
	u := order.Uint64(b.Buf[b.readPos:])
	b.readPos += 8
	return u, nil
}

func (b *Buffer) Skip(i int) int {
	pos := b.readPos + i
	if pos >= 0 && pos < len(b.Buf) {
		b.readPos = pos
		return pos
	}

	return -1
}

func (b *Buffer) ReadAt(p []byte, off int64) (int, error) {
	if off < 0 {
		return 0, errors.New("binary.Buffer.ReadAt: negative offset")
	}

	if int(off) >= len(b.Buf) {
		return 0, io.EOF
	}
	n := len(p)
	if n+int(off) > len(b.Buf) {
		n = len(b.Buf) - int(off)
	}
	copy(p, b.Buf[off:])
	return n, nil
}

func (b *Buffer) ReadRune() (rune, int, error) {
	if b.readPos == len(b.Buf) {
		return 0, 0, io.EOF
	}
	if c := b.Buf[b.readPos]; c < utf8.RuneSelf {
		b.readPos += 1
		return rune(c), 1, nil
	}
	r, n := utf8.DecodeRune(b.Buf[b.readPos:])
	b.readPos += n
	return r, n, nil
}

func (b *Buffer) ReadBytes(delim byte) ([]byte, error) {
	if b.readPos >= len(b.Buf) {
		return nil, io.EOF
	}
	s := b.readPos
	for i := b.readPos; i < len(b.Buf); i++ {
		if b.Buf[i] == delim {
			b.readPos = i + 1
			return b.Buf[s:b.readPos], nil
		}
	}
	return nil, io.EOF
}

func (b *Buffer) Write(p []byte) (int, error) {
	i := b.Grows(len(p))
	copy(b.Buf[i:], p)
	return len(p), nil
}

func (b *Buffer) WriteString(s string) (int, error) {
	return b.Write(unsafe2.Bytes(s))
}

func (b *Buffer) WriteByte(c byte) error {
	i := b.Grows(1)
	b.Buf[i] = c
	return nil
}

func (b *Buffer) WriteUint16(u uint16, order binary.ByteOrder) error {
	i := b.Grows(2)
	order.PutUint16(b.Buf[i:], u)
	return nil
}

func (b *Buffer) WriteUint32(u uint32, order binary.ByteOrder) error {
	i := b.Grows(4)
	order.PutUint32(b.Buf[i:], u)
	return nil
}

func (b *Buffer) WriteUint64(u uint64, order binary.ByteOrder) error {
	i := b.Grows(8)
	order.PutUint64(b.Buf[i:], u)
	return nil
}

func (b *Buffer) WriteRune(r rune) (int, error) {
	i := b.Grows(utf8.UTFMax)
	s := utf8.EncodeRune(b.Buf[i:], r)
	n := utf8.UTFMax - s
	b.Buf = b.Buf[:len(b.Buf)-n]
	return s, nil
}

func (b *Buffer) String() string {
	return string(b.Bytes())
}
