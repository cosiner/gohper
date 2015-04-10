package log

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"

	t "github.com/cosiner/gohper/lib/time"

	"github.com/cosiner/gohper/config"

	"github.com/cosiner/gohper/lib/types"
)

const (
	// log dir permission when create
	_LOGDIR_PERM = 0755
)

var timeFormat = t.FormatLayout("yyyymmdd-HHMMSS")

//==============================================================================
//                           Log Buffer
//==============================================================================
// logBuffer represent a log writer for a special level
type logBuffer struct {
	file *os.File
	*bufio.Writer
	nbytes  uint64
	logdir  string
	level   string
	bufsize uint64
	maxsize uint64
}

// newLogBuffer create a new log buffer
func newLogBuffer(logdir, level string, bufsize, maxsize uint64) (*logBuffer, error) {
	buf := &logBuffer{
		logdir:  logdir,
		level:   level,
		bufsize: bufsize,
		maxsize: maxsize,
	}
	return buf, buf.newLogFile()
}

// newLogFile create a new log file
func (buf *logBuffer) newLogFile() (err error) {
	if buf.file != nil {
		buf.Flush()
		buf.file.Close()
	}
	filename := fmt.Sprintf("%s.log.%s.%d",
		buf.level, time.Now().Format(timeFormat), os.Getpid())
	if buf.file, err = os.Create(filepath.Join(buf.logdir, filename)); err == nil {
		buf.nbytes = 0
		buf.Writer = bufio.NewWriterSize(buf.file, int(buf.bufsize))
	}
	return
}

// flush flush log buffer
func (buf *logBuffer) flush() (err error) {
	if err = buf.Flush(); err == nil {
		err = buf.file.Sync()
	}
	return
}

// close close the log buffer
func (buf *logBuffer) close() {
	buf.Flush()
	buf.file.Close()
	return
}

// write write log message to log file
func (buf *logBuffer) write(msg string) (err error) {
	if buf.nbytes+uint64(len(msg)) >= buf.maxsize {
		if err = buf.newLogFile(); err != nil {
			return
		}
	}
	n, err := buf.WriteString(msg)
	buf.nbytes += uint64(n)
	return
}

//==============================================================================
//                          File Log Writer
//==============================================================================
// logWrite is actuall log writer, output is local file
type FileWriter struct {
	level Level
	files []*logBuffer
}

// Config resolv config, format like bufsize=xxx&maxsize=xxx&logdir=xxx&level=info
func (writer *FileWriter) Config(conf string) (err error) {
	c := config.NewConfig(config.LINE)
	if err = c.ParseString(conf); err != nil {
		return
	}
	writer.level = ParseLevel(c.ValDef("level", "info"))
	if writer.level == LEVEL_OFF {
		return
	}

	logdir := c.ValDef("logdir", filepath.Join(os.TempDir(), "gologs"))
	bufsize, err := types.BytesCount(c.ValDef("bufsize", "10K"))
	maxsize, err := types.BytesCount(c.ValDef("maxsize", "10M"))
	err = os.MkdirAll(logdir, _LOGDIR_PERM)

	writer.files = make([]*logBuffer, _LEVEL_MAX+1)
	for l := writer.level; l <= _LEVEL_MAX && err == nil; l++ {
		writer.files[l], err = newLogBuffer(logdir, l.String(), bufsize, maxsize)
	}
	return
}

// Write write log to log file, higher level log will simultaneously
// output to all lower level log file
func (writer *FileWriter) Write(log *Log) (err error) {
	for l := writer.level; l <= log.Level; l++ {
		if writer.files[l].write(log.String()) != nil {
			return
		}
	}
	return
}

// Flush flush log writer
func (writer *FileWriter) Flush() {
	for l := writer.level; l <= _LEVEL_MAX; l++ {
		writer.files[l].flush()
	}
}

// Close close log writer
func (writer *FileWriter) Close() {
	for l := writer.level; l <= _LEVEL_MAX; l++ {
		writer.files[l].close()
		writer.files[l] = nil
	}
}
