// Package color enable color output for terminal use ansi escape code
package color

import (
	"strings"

	"github.com/cosiner/gohper/terminal/ansi"
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
	fgColors = map[string]string{
		BLACK:   ansi.FgBlack,
		RED:     ansi.FgRed,
		GREEN:   ansi.FgGreen,
		YELLOW:  ansi.FgYellow,
		BLUE:    ansi.FgBlue,
		MAGENTA: ansi.FgMagenta,
		CYAN:    ansi.FgCyan,
		WHITE:   ansi.FgWhite,
	}

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

// FgHl is a quick way to New().Fg(color).Highlight().Finish()
func FgHl(color string) *Renderer {
	return New(fgColors[color], ansi.Highlight)
}

// Fg is a quick way to New().Fg(color).Finish()
func Fg(color string) *Renderer {
	return New(fgColors[color])
}

// Has check whether a color can be use
func Has(c string) bool {
	_, has := fgColors[strings.ToLower(c)]

	return has
}
