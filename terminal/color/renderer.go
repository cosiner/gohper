package color

import (
	"fmt"
	"io"

	"github.com/cosiner/gohper/terminal/ansi"
	"github.com/cosiner/gohper/terminal/color/output"
	"github.com/cosiner/gohper/unsafe2"
)

type Writer struct {
	Prefix string
	io.Writer
	*Renderer
}

func (w Writer) Write(bs []byte) (int, error) {
	i, err := w.Writer.Write(unsafe2.Bytes(w.Prefix))
	if err != nil {
		return 0, err
	}

	n, err := w.RenderTo(w.Writer, unsafe2.String(bs))
	if err != nil {
		return 0, err
	}

	return n + i, nil
}

// Renderer is a render for terminal string
type Renderer struct {
	enable   bool
	settings string
	end      string
}

// New a terminal color render,
func New(codes ...string) *Renderer {
	return &Renderer{
		enable:   true,
		settings: ansi.Begin(codes...),
		end:      ansi.End(),
	}
}

// Disable color render
func (r *Renderer) Disable() {
	r.enable = false
}

// Render rend a string
func (r *Renderer) Render(str string) string {
	if str == "" || !r.enable {
		return str
	}

	return r.settings + str + r.end
}

// RenderTo render string to writer
func (r *Renderer) RenderTo(w io.Writer, str string) (int, error) {
	if str == "" || !r.enable {
		return w.Write(unsafe2.Bytes(str))
	}

	if err := r.Begin(w); err == nil {
		c, err := w.Write(unsafe2.Bytes(str))
		r.End(w)

		return c, err
	} else {
		return 0, err
	}
}

func (r *Renderer) Begin(w io.Writer) error {
	_, err := w.Write(unsafe2.Bytes(r.settings))

	return err
}

func (r *Renderer) End(w io.Writer) error {
	_, err := w.Write(unsafe2.Bytes(r.end))

	return err
}

func (r *Renderer) Writer(prefix string, w io.Writer) io.Writer {
	return Writer{
		Prefix:   prefix,
		Writer:   w,
		Renderer: r,
	}
}

func (r *Renderer) Fprint(w io.Writer, args ...interface{}) (int, error) {
	return r.RenderTo(w, fmt.Sprint(args...))
}

func (r *Renderer) Fprintln(w io.Writer, args ...interface{}) (int, error) {
	return r.RenderTo(w, fmt.Sprintln(args...))
}

func (r *Renderer) Fprintf(w io.Writer, format string, args ...interface{}) (int, error) {
	return r.RenderTo(w, fmt.Sprintf(format, args...))
}

func (r *Renderer) Print(args ...interface{}) (int, error) {
	return r.RenderTo(output.Stdout, fmt.Sprint(args...))
}

func (r *Renderer) Println(args ...interface{}) (int, error) {
	return r.RenderTo(output.Stdout, fmt.Sprintln(args...))
}

func (r *Renderer) Printf(format string, args ...interface{}) (int, error) {
	return r.RenderTo(output.Stdout, fmt.Sprintf(format, args...))
}

func (r *Renderer) Error(args ...interface{}) (int, error) {
	return r.RenderTo(output.Stderr, fmt.Sprint(args...))
}

func (r *Renderer) Errorln(args ...interface{}) (int, error) {
	return r.RenderTo(output.Stderr, fmt.Sprintln(args...))
}

func (r *Renderer) Errorf(format string, args ...interface{}) (int, error) {
	return r.RenderTo(output.Stderr, fmt.Sprintf(format, args...))
}

func (r *Renderer) Sprint(args ...interface{}) string {
	return r.Render(fmt.Sprint(args...))
}

func (r *Renderer) Sprintln(args ...interface{}) string {
	return r.Render(fmt.Sprintln(args...))
}

func (r *Renderer) Sprintf(format string, args ...interface{}) string {
	return r.Render(fmt.Sprintf(format, args...))
}
