// Package termcolor enable color output for terminal use ansi escape code
//
// Usage:
// Create a instance by New(), configure through Bg, Fg, ..., finally Finish() configuration.
//
// Then Render a string or RenderTo a writer, or Begin, operations..., End.
package termcolor

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mattn/go-colorable"

	"github.com/cosiner/gohper/os2"
	"github.com/cosiner/gohper/strings2"
	"github.com/cosiner/gohper/unsafe2"
)

const (
	BLACK   = "black"
	RED     = "red"
	GREEN   = "green"
	YELLOW  = "yellow"
	BLUE    = "blue"
	MAGENTA = "magenta"
	CYAN    = "cyan"
	WHITE   = "white"
)

var (
	// compatible for Windows
	Stdout io.Writer = os.Stdout
	Stderr io.Writer = os.Stderr
)

func init() {
	if os2.IsWindows() {
		Stdout = colorable.NewColorableStdout()
		Stderr = colorable.NewColorableStderr()
	}
}

var (
	// frontground: 30:黑 31:红 32:绿 33:黄 34:蓝色 35:紫色 36:深绿 37:白色
	// background: 30:black 31:red 32:green 33:yellow 34:blue 35:purple 36:deep green 37:white
	// background: 40, 41, ...
	Colors = map[string]int{
		BLACK:   0,
		RED:     1,
		GREEN:   2,
		YELLOW:  3,
		BLUE:    4,
		MAGENTA: 5,
		CYAN:    6,
		WHITE:   7,
	}
	// only background
	LightBlack   = FgHl(BLACK)
	LightRed     = FgHl(RED)
	LightGreen   = FgHl(GREEN)
	LightYellow  = FgHl(YELLOW)
	LightBlue    = FgHl((BLUE))
	LightMagenta = FgHl(MAGENTA)
	LightCyan    = FgHl(CYAN)
	LightWhite   = FgHl(WHITE)

	Black   = Fg(BLACK)
	Red     = Fg(RED)
	Green   = Fg(GREEN)
	Yellow  = Fg(YELLOW)
	Blue    = Fg((BLUE))
	Magenta = Fg(MAGENTA)
	Cyan    = Fg(CYAN)
	White   = Fg(WHITE)
)

// TermColor is a render for terminal string
type TermColor struct {
	enable        bool
	fg            int
	bg            int
	underline     bool
	blink         bool
	highlight     bool
	inverse       bool
	hidden        bool
	settingsCount int
	settings      string
}

type Writer struct {
	Prefix string
	io.Writer
	*TermColor
}

func (w Writer) Write(bs []byte) (int, error) {
	i, err := w.Writer.Write(unsafe2.Bytes(w.Prefix))
	if err == nil {
		n, e := w.RenderTo(w.Writer, unsafe2.String(bs))
		if e == nil {
			return n + i, err
		}
		err = e
	}
	return 0, err
}

// New create a new terminal color render
func New() *TermColor {
	return &TermColor{
		fg:     -1,
		bg:     -1,
		enable: true,
	}
}

// FgHl is a quick way to New().Fg(color).Highlight().Finish()
func FgHl(color string) *TermColor {
	return New().Fg(color).Highlight().Finish()
}

// Fg is a quick way to New().Fg(color).Finish()
func Fg(color string) *TermColor {
	return New().Fg(color).Finish()
}

// Disable disable color render
func (tc *TermColor) Disable() {
	tc.enable = false
}

// Render rend a string
func (tc *TermColor) Render(str string) string {
	if str == "" || !tc.enable {
		return str
	}
	return "\033[" + tc.settings + "m" + str + "\033[0m"
	// return fmt.Sprintf("\033[%sm%s\033[0m", tc.settings, str)
}

// RenderTo render string to writer
func (tc *TermColor) RenderTo(w io.Writer, str string) (int, error) {
	if str == "" || !tc.enable {
		return w.Write(unsafe2.Bytes(str))
	}

	if err := tc.Begin(w); err == nil {
		c, err := w.Write(unsafe2.Bytes(str))
		tc.End(w)
		return c, err
	} else {
		return 0, err
	}
}

func (tc *TermColor) Begin(w io.Writer) error {
	_, err := w.Write(unsafe2.Bytes("\033[" + tc.settings + "m"))
	return err
}

func (tc *TermColor) End(w io.Writer) error {
	_, err := w.Write(unsafe2.Bytes("\033[0m"))
	return err
}

// Bg set render's background color
func (tc *TermColor) Bg(bg string) *TermColor {
	tc.settingsCount++
	tc.bg = Colors[strings.ToLower(bg)]
	return tc
}

// Fg set render's foreground color
func (tc *TermColor) Fg(fg string) *TermColor {
	tc.settingsCount++
	tc.fg = Colors[strings.ToLower(fg)]
	return tc
}

// Highlight enable render to highlight
func (tc *TermColor) Highlight() *TermColor {
	tc.settingsCount++
	tc.highlight = true
	return tc
}

// Underline enable render to underline
func (tc *TermColor) Underline() *TermColor {
	tc.settingsCount++
	tc.underline = true
	return tc
}

// Blink enable render to blink
func (tc *TermColor) Blink() *TermColor {
	tc.settingsCount++
	tc.blink = true
	return tc
}

// Inverse enable render to inverse color
func (tc *TermColor) Inverse() *TermColor {
	tc.settingsCount++
	tc.inverse = true
	return tc
}

// Hidden enable render to hidden color
func (tc *TermColor) Hidden() *TermColor {
	tc.settingsCount++
	tc.hidden = true
	return tc
}

// Finish complete color settings
func (tc *TermColor) Finish() *TermColor {
	color := make([]int, tc.settingsCount)
	index := 0
	if tc.fg != -1 {
		color[index] = tc.fg + 30
		index++
	}
	if tc.bg != -1 {
		color[index] = tc.bg + 40
		index++
	}
	if tc.highlight {
		color[index] = 1
		index++
	}
	if tc.underline {
		color[index] = 4
		index++
	}
	if tc.blink {
		color[index] = 5
		index++
	}
	if tc.inverse {
		color[index] = 7
		index++
	}
	if tc.hidden {
		color[index] = 8
		index++
	}
	tc.settings = strings2.JoinInt(color, ";")
	return tc
}

func (tc *TermColor) Writer(prefix string, w io.Writer) io.Writer {
	return Writer{
		Prefix:    prefix,
		Writer:    w,
		TermColor: tc,
	}
}

func (tc *TermColor) Fprint(w io.Writer, args ...interface{}) (int, error) {
	return tc.RenderTo(w, fmt.Sprint(args...))
}

func (tc *TermColor) Fprintln(w io.Writer, args ...interface{}) (int, error) {
	return tc.RenderTo(w, fmt.Sprintln(args...))
}

func (tc *TermColor) Fprintf(w io.Writer, format string, args ...interface{}) (int, error) {
	return tc.RenderTo(w, fmt.Sprintf(format, args...))
}

func (tc *TermColor) Print(args ...interface{}) (int, error) {
	return tc.RenderTo(Stdout, fmt.Sprint(args...))
}

func (tc *TermColor) Println(args ...interface{}) (int, error) {
	return tc.RenderTo(Stdout, fmt.Sprintln(args...))
}

func (tc *TermColor) Printf(format string, args ...interface{}) (int, error) {
	return tc.RenderTo(Stdout, fmt.Sprintf(format, args...))
}
