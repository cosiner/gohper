package log

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cosiner/golib/errors"

	"github.com/cosiner/golib/sys"
	"github.com/cosiner/golib/types"
)

const (
	// log dir permission when create
	_LOGDIR_PERM = 0755
	CONF_BUFSIZE = "bufsize"
	CONF_MAXSIZE = "maxsize"
	CONF_LOGDIR  = "logdir"
)

//==============================================================================
//                           Log Buffer
//==============================================================================
// logBuffer represent a log writer for a special level
type logBuffer struct {
	writer *FileLogWriter
	file   *os.File
	*bufio.Writer
	level  Level
	nbytes uint64
}

// newLogBuffer create a new log buffer
func newLogBuffer(writer *FileLogWriter, level Level) (*logBuffer, error) {
	buf := &logBuffer{writer, nil, nil, level, 0}
	return buf, buf.newLogFile()
}

// newLogFile create a new log file
func (buf *logBuffer) newLogFile() (err error) {
	if buf.file != nil {
		buf.Flush()
		buf.file.Close()
	}
	buf.file, err = createLogFile(buf.writer.logdir, buf.level.String())
	buf.nbytes = 0
	if err == nil {
		buf.Writer = bufio.NewWriterSize(buf.file, int(buf.writer.bufsize))
	}
	return
}

// flush flush log buffer
func (buf *logBuffer) flush() (err error) {
	err = buf.Flush()
	err = buf.file.Sync()
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
	if buf.nbytes+uint64(len(msg)) >= buf.writer.maxsize {
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
type FileLogWriter struct {
	level   Level
	bufsize uint64
	maxsize uint64
	logdir  string
	files   [_LEVEL_MAX + 1]*logBuffer
}

// parseConf parse a pair of config
func (writer *FileLogWriter) parseConf(pair *types.Pair) (err error) {
	err = errors.Errorf("Wrong config format %s", pair.String())
	if pair.HasKey() && pair.HasValue() {
		key := strings.ToLower(pair.Key)
		if key == CONF_LOGDIR {
			writer.logdir, err = pair.Value, nil
		} else if size, e := types.Str2Bytes(pair.Value); e == nil {
			if key == CONF_BUFSIZE {
				writer.bufsize, err = size, nil
			} else if key == CONF_MAXSIZE {
				writer.maxsize, err = size, nil
			}
		}
	}
	return
}

// Config resolv config
func (writer *FileLogWriter) Config(conf string) (err error) {
	confs := strings.FieldsFunc(conf, func(r rune) bool {
		return r == '&'
	})
	if len(confs) == 0 {
		return errors.Err("No config found")
	} else {
		writer.logdir = filepath.Join(os.TempDir(), "gologs")
		for _, c := range confs {
			if err := writer.parseConf(types.ParsePair(c, "=")); err != nil {
				return err
			}
		}
		return sys.MkdirWithParent(writer.logdir)
	}
}

// Write write log to log file
func (writer *FileLogWriter) Write(log *Log) (err error) {
	for l := writer.level; l <= log.Level; l++ {
		err = writer.files[l].write(
			fmt.Sprintf("[%s] %s %s", log.Level.String(), datetime(), log.Message))
		if err != nil {
			return
		}
	}
	return
}

// Flush flush log writer
func (writer *FileLogWriter) Flush() {
	for l := writer.level; l <= _LEVEL_MAX; l++ {
		writer.files[l].flush()
	}
}

// Close close log writer
func (writer *FileLogWriter) Close() {
	for l := writer.level; l <= _LEVEL_MAX; l++ {
		writer.files[l].close()
		writer.files[l] = nil
	}
}

// ResetLevel reset log level
func (writer *FileLogWriter) ResetLevel(level Level) (err error) {
	for l := _LEVEL_MIN; l < level; l++ {
		if buf := writer.files[l]; buf != nil {
			buf.close()
		}
		writer.files[l] = nil
	}
	for l := level; l <= _LEVEL_MAX; l++ {
		if writer.files[l] == nil {
			writer.files[l], err = newLogBuffer(writer, l)
			if err != nil {
				return
			}
		}
	}
	writer.level = level
	return
}

// createLogFile create log file in the format:level.log.yyyymmdd-HHMMSS.pid
func createLogFile(dir string, filename string) (*os.File, error) {
	t := timenow()
	filename = fmt.Sprintf("%s.log.%04d%02d%02d-%02d%02d%02d.%d",
		filename,
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		os.Getpid())
	fname := filepath.Join(dir, filename)
	f, err := os.Create(fname)
	return f, err
}
