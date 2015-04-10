package log

import (
	"fmt"
	"io"
	"strings"

	"github.com/cosiner/gohper/config"
	"github.com/cosiner/gohper/lib/termcolor"
	"github.com/mattn/go-colorable"
)

// bgColor create color render use given background color, default highlight
func bgColor(bg string) *termcolor.TermColor {
	return termcolor.NewColor().Highlight().Bg(bg)
}

// defTermColor define default color for each log level
var defTermColor = [5]*termcolor.TermColor{
	bgColor(termcolor.GREEN),  //debug
	bgColor(termcolor.WHITE),  //info
	bgColor(termcolor.YELLOW), //warn
	bgColor(termcolor.BLUE),   //error
	bgColor(termcolor.RED),    //fatal
}

// ConsoleWriter output log to console
type ConsoleWriter struct {
	termColor [5]*termcolor.TermColor
	out       io.Writer
	err       io.Writer
}

// Config config console log writer
// parameter conf can use to config color for each log level, such as
// warn="black"&info="green"&error="red"...
func (clw *ConsoleWriter) Config(conf string) error {
	clw.out = colorable.NewColorableStdout()
	clw.err = colorable.NewColorableStderr()
	clw.termColor = defTermColor
	if conf != "" {
		c := config.NewConfig(config.LINE)
		c.ParseString(conf)
		if _, has := c.Val("disableColor"); has {
			clw.DisableColor()
		} else {
			for l := _LEVEL_MIN; l < _LEVEL_MAX; l++ {
				s := strings.ToLower(l.String())
				if color := c.ValDef(s, ""); color != "" {
					clw.termColor[l] = bgColor(color)
				}
			}
		}
	}
	return nil
}

// DisableColor disable color output
func (clw *ConsoleWriter) DisableColor() {
	for _, tc := range clw.termColor {
		tc.Disable()
	}
}

// Write write
func (clw *ConsoleWriter) Write(log *Log) error {
	out := clw.out
	if log.Level >= LEVEL_ERROR {
		out = clw.err
	}
	_, err := fmt.Fprint(out, clw.termColor[log.Level].Render(log.String()))
	return err
}

func (clw *ConsoleWriter) Flush() {}
func (clw *ConsoleWriter) Close() {}
