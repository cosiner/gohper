package ansi

import "strings"

func String(code string) string {
	return "\033[" + code + "m"
}

const (
	FgBlack   = "30"
	FgRed     = "31"
	FgGreen   = "32"
	FgYellow  = "33"
	FgBlue    = "34"
	FgMagenta = "35"
	FgCyan    = "36"
	FgWhite   = "37"

	BgBlack   = "40"
	BgRed     = "41"
	BgGreen   = "42"
	BgYellow  = "43"
	BgBlue    = "44"
	BgMagenta = "45"
	BgCyan    = "46"
	BgWhite   = "47"

	Highlight = "1"
	UnderLine = "4"
	Blink     = "5"
	Inverse   = "6"
	Hidden    = "7"

	Reset = "0"
)

func Begin(codes ...string) string {
	return "\033[" + strings.Join(codes, ";") + "m"
}

func End() string {
	return "\033[0m"
}

func Render(s string, codes ...string) string {
	return Begin(codes...) + s + End()
}
