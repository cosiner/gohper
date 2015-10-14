package color

import "bytes"

type Code string

const (
	FgBlack   Code = "30"
	FgRed     Code = "31"
	FgGreen   Code = "32"
	FgYellow  Code = "33"
	FgBlue    Code = "34"
	FgMagenta Code = "35"
	FgCyan    Code = "36"
	FgWhite   Code = "37"

	BgBlack   Code = "40"
	BgRed     Code = "41"
	BgGreen   Code = "42"
	BgYellow  Code = "43"
	BgBlue    Code = "44"
	BgMagenta Code = "45"
	BgCyan    Code = "46"
	BgWhite   Code = "47"

	Highlight Code = "1"
	UnderLine Code = "4"
	Blink     Code = "5"
	Inverse   Code = "6"
	Hidden    Code = "7"

	_Reset Code = "0"
)

func (c Code) String() string {
	return "\033[" + string(c) + "m"
}

func Begin(codes ...Code) string {
	buf := bytes.NewBuffer(make([]byte, 0, len(codes)*3))
	buf.WriteString("\033[")
	for i := 0; i < len(codes); i++ {
		if i != 0 {
			buf.WriteByte(';')
		}
		buf.WriteString(string(codes[i]))
	}
	buf.WriteByte('m')
	return buf.String()
}

func End() string {
	return "\033[0m"
}
