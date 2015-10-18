package color

import (
	"fmt"
	"io"
	"os"

	"github.com/cosiner/gohper/unsafe2"
)

var (
	Black   = New(FgBlack, Highlight)
	Red     = New(FgRed, Highlight)
	Green   = New(FgGreen, Highlight)
	Yellow  = New(FgYellow, Highlight)
	Blue    = New(FgBlue, Highlight)
	Magenta = New(FgMagenta, Highlight)
	Cyan    = New(FgCyan, Highlight)
	White   = New(FgWhite, Highlight)
)

// Renderer is a render for terminal string
type Renderer struct {
	begin []byte
	end   []byte
}

// New a terminal color render,
func New(codes ...Code) *Renderer {
	return &Renderer{
		begin: []byte(Begin(codes...)),
		end:   []byte(End()),
	}
}

func (r *Renderer) Writer(prefix []byte, w io.Writer) *Writer {
	return &Writer{
		prefix:   prefix,
		Writer:   w,
		Renderer: r,
	}
}

// Render rend a string
func (r *Renderer) Render(str []byte) []byte {
	bl, sl, el := len(r.begin), len(str), len(r.end)
	data := make([]byte, bl+sl+el)
	copy(data, r.begin)
	copy(data[bl:], str)
	copy(data[bl+sl:], r.end)
	return data
}

func (r *Renderer) RenderString(str string) string {
	return unsafe2.String(r.begin) + str + unsafe2.String(r.end)
}

// RenderTo render string to writer
func (r *Renderer) RenderTo(w io.Writer, str ...[]byte) (int, error) {
	err := r.Begin(w)
	var (
		n int
		c int
	)
	for i := 0; err == nil && i < len(str); i++ {
		c, err = w.Write(str[i])
		n += c
	}
	r.End(w)
	return n, err
}

func (r *Renderer) RenderStringTo(w io.Writer, str ...string) (int, error) {
	err := r.Begin(w)
	var (
		n int
		c int
	)
	for i := 0; err == nil && i < len(str); i++ {
		c, err = w.Write(unsafe2.Bytes(str[i]))
		n += c
	}
	r.End(w)
	return n, err
}

func (r *Renderer) Begin(w io.Writer) error {
	_, err := w.Write(r.begin)

	return err
}

func (r *Renderer) End(w io.Writer) error {
	_, err := w.Write(r.end)

	return err
}

func (r *Renderer) Fprint(w io.Writer, args ...interface{}) (int, error) {
	s := fmt.Sprint(args...)
	return r.RenderTo(w, unsafe2.Bytes(s))
}

func (r *Renderer) Fprintln(w io.Writer, args ...interface{}) (int, error) {
	s := fmt.Sprintln(args...)
	return r.RenderTo(w, unsafe2.Bytes(s))
}

func (r *Renderer) Fprintf(w io.Writer, format string, args ...interface{}) (int, error) {
	s := fmt.Sprintf(format, args...)
	return r.RenderTo(w, unsafe2.Bytes(s))
}

func (r *Renderer) Print(args ...interface{}) (int, error) {
	return r.Fprint(Stdout, args...)
}

func (r *Renderer) Println(args ...interface{}) (int, error) {
	return r.Fprintln(Stdout, args...)
}

func (r *Renderer) Printf(format string, args ...interface{}) (int, error) {
	return r.Fprintf(os.Stdout, format, args...)
}

func (r *Renderer) Error(args ...interface{}) (int, error) {
	return r.Fprint(Stderr, args...)
}

func (r *Renderer) Errorln(args ...interface{}) (int, error) {
	return r.Fprintln(Stderr, args...)
}

func (r *Renderer) Errorf(format string, args ...interface{}) (int, error) {
	return r.Fprintf(Stderr, format, args...)
}

func (r *Renderer) Sprint(args ...interface{}) string {
	return r.RenderString(fmt.Sprint(args...))
}

func (r *Renderer) Sprintln(args ...interface{}) string {
	return r.RenderString(fmt.Sprintln(args...))
}

func (r *Renderer) Sprintf(format string, args ...interface{}) string {
	return r.RenderString(fmt.Sprintf(format, args...))
}

type Writer struct {
	prefix []byte
	io.Writer
	*Renderer
}

func (w *Writer) Write(bs []byte) (int, error) {
	w.Renderer.Begin(w.Writer)
	n1, err := w.Writer.Write(w.prefix)
	if err != nil {
		return n1, err
	}
	n2, err := w.Writer.Write(bs)
	if err != nil {
	}
	w.Renderer.End(w.Writer)

	return n1 + n2, err
}
