package log

import (
	"fmt"
	"strings"

	"github.com/cosiner/golib/termcolor"
	"github.com/cosiner/gomodule/config"
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

// ConsoleLogWriter output log to console
type ConsoleLogWriter struct {
	termColor [5]*termcolor.TermColor
}

// Config config console log writer
// parameter conf can use to config color for each log level, such as
// warn="black"&info="green"&error="red"...
func (clw *ConsoleLogWriter) Config(conf string) error {
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
func (clw *ConsoleLogWriter) DisableColor() {
	for _, tc := range clw.termColor {
		tc.Disable()
	}
}

// Write write
func (clw *ConsoleLogWriter) Write(log *Log) error {
	_, err := fmt.Print(clw.termColor[log.Level].Render(log.String()))
	return err
}

func (clw *ConsoleLogWriter) ResetLevel(level Level) error { return nil }
func (clw *ConsoleLogWriter) Flush()                       {}
func (clw *ConsoleLogWriter) Close()                       {}
