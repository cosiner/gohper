package log

import (
	"io"
	"strings"

	"github.com/cosiner/gohper/lib/termcolor"
	"github.com/mattn/go-colorable"
)

// defTermColor define default color for each log level
var defTermColor = [5]*termcolor.TermColor{
	bgColor(termcolor.GREEN),  //debug
	bgColor(termcolor.WHITE),  //info
	bgColor(termcolor.YELLOW), //warn
	bgColor(termcolor.BLUE),   //error
	bgColor(termcolor.RED),    //fatal
}

type (
	ConsoleWriterOption struct {
		DisableColor bool
		Colors       map[string]string
	}

	// ConsoleWriter output log to console
	ConsoleWriter struct {
		termColor [5]*termcolor.TermColor
		out       io.Writer
		err       io.Writer
	}
)

// bgColor create color render use given background color, default highlight
func bgColor(bg string) *termcolor.TermColor {
	return termcolor.New().Highlight().Bg(bg).Finish()
}

// Config config console log writer
// parameter conf can use to config color for each log level, such as
// warn="black"&info="green"&error="red"...
func (w *ConsoleWriter) Config(conf interface{}) error {
	var opt *ConsoleWriterOption
	if conf == nil {
		opt = &ConsoleWriterOption{}
	} else {
		switch c := conf.(type) {
		case *ConsoleWriterOption:
			opt = c
		case ConsoleWriterOption:
			opt = &c
		default:
			return ErrInvalidConfig
		}
	}

	w.out = colorable.NewColorableStdout()
	w.err = colorable.NewColorableStderr()

	w.termColor = defTermColor
	if opt.DisableColor {
		w.DisableColor()
	} else if len(opt.Colors) != 0 {
		for l := _LEVEL_MIN; l < _LEVEL_MAX; l++ {
			s := strings.ToLower(l.String())
			if color := opt.Colors[s]; color != "" {
				w.termColor[l] = bgColor(color)
			}
		}
	}
	return nil
}

func (w *ConsoleWriter) SetLevel(l Level) {}

// DisableColor disable color output
func (w *ConsoleWriter) DisableColor() {
	for _, tc := range w.termColor {
		tc.Disable()
	}
}

// Write write
func (w *ConsoleWriter) Write(log *Log) error {
	out := w.out
	if log.Level >= LEVEL_PANIC {
		out = w.err
	}
	tc := w.termColor[log.Level]
	tc.Begin(out)
	log.WriteTo(out)
	return tc.End(out)
}

func (w *ConsoleWriter) Flush() {}
func (w *ConsoleWriter) Close() {}
