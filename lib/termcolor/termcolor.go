// Package termcolor enable color output for terminal use ansi escape code
//
// Usage:
// Create a instance by New(), configure through Bg, Fg, ..., finally Finish() configuration.
//
// Then Render a string or RenderTo a writer, or Begin, operations..., End.
package termcolor

import (
	"io"
	"strings"

	"github.com/cosiner/gohper/lib/types"
)

const (
	BLACK      = "black"
	RED        = "red"
	GREEN      = "green"
	YELLOW     = "yellow"
	BLUE       = "blue"
	PURPLE     = "purple"
	DEEP_GREEN = "deep_green"
	WHITE      = "white"
)

// background: 30:黑 31:红 32:绿 33:黄 34:蓝色 35:紫色 36:深绿 37:白色
// background: 30:black 31:red 32:green 33:yellow 34:blue 35:purple 36:deep green 37:white
// frontground: 40, 41, ...
var Colors = map[string]int{
	BLACK:      0,
	RED:        1,
	GREEN:      2,
	YELLOW:     3,
	BLUE:       4,
	PURPLE:     5,
	DEEP_GREEN: 6,
	WHITE:      7,
}

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

// New create a new terminal color render
func New() *TermColor {
	return &TermColor{
		fg:     -1,
		bg:     -1,
		enable: true,
	}
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
func (tc *TermColor) RenderTo(w io.Writer, str string) error {
	var err error
	if str == "" || !tc.enable {
		_, err = w.Write(types.UnsafeBytes(str))
	} else {
		if err = tc.Begin(w); err == nil {
			_, err = w.Write(types.UnsafeBytes(str))
			tc.End(w)
		}
	}
	return err
}

func (tc *TermColor) Begin(w io.Writer) error {
	_, err := w.Write(types.UnsafeBytes("\033[" + tc.settings + "m"))
	return err
}

func (tc *TermColor) End(w io.Writer) error {
	_, err := w.Write(types.UnsafeBytes("\033[0m"))
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
		color[index] = tc.fg + 40
		index++
	}
	if tc.bg != -1 {
		color[index] = tc.bg + 30
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
	tc.settings = types.JoinInt(color, ";")
	return tc
}
