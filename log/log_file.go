package log

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	// log dir permission when create
	_LOGDIR_PERM = 0755
)

// logBuffer represent a log writer for a special level
type logBuffer struct {
	lw   *logWriter
	file *os.File
	*bufio.Writer
	level  Level
	nbytes uint64
}

// newLogBuffer create a new log buffer
func newLogBuffer(lw *logWriter, level Level) (*logBuffer, error) {
	lb := &logBuffer{lw, nil, nil, level, 0}
	return lb, lb.rotateFile()
}

// rotateFile create a new log file
func (lb *logBuffer) rotateFile() (err error) {
	if lb.file != nil {
		lb.Flush()
		lb.file.Close()
	}
	lb.file, err = createLogFile(lb.lw.logDir, lb.level.Name())
	lb.nbytes = 0
	if err == nil {
		lb.Writer = bufio.NewWriterSize(lb.file, int(lb.lw.bufSize))
	}
	return
}

// flush flush log buffer
func (lb *logBuffer) flush() (err error) {
	err = lb.Flush()
	err = lb.file.Sync()
	return
}

// close close the log buffer
func (lb *logBuffer) close() {
	lb.Flush()
	lb.file.Close()
	return
}

// write write log message to log file
func (lb *logBuffer) write(msg string) (err error) {
	if lb.nbytes+uint64(len(msg)) >= lb.lw.maxSize {
		if err = lb.rotateFile(); err != nil {
			return
		}
	}
	n, err := lb.WriteString(msg)
	if msg[len(msg)-1] != '\n' {
		lb.WriteByte('\n')
	}
	lb.nbytes += uint64(n)
	return
}

// logWrite is actuall log writer, output is local file
type logWriter struct {
	level   Level
	bufSize uint64
	maxSize uint64
	logDir  string
	files   [_LEVEL_NUM]*logBuffer
}

// newLogWriter create a new log writer
func newLogWriter(level Level, bufSize, maxSize uint64, logDir string) (lw *logWriter, err error) {
	err = os.Mkdir(logDir, _LOGDIR_PERM)
	if err == nil || strings.Contains(err.Error(), "file exists") {
		err = nil
		lw = &logWriter{level: level,
			logDir:  logDir,
			bufSize: bufSize,
			maxSize: maxSize}
		for l := level; l < LEVEL_MAX; l++ {
			lw.files[l], err = newLogBuffer(lw, l)
			if err != nil {
				return nil, err
			}
		}
	}
	return lw, err
}

// Write write log to log file
func (lw *logWriter) Write(log *Log) (err error) {
	for l := lw.level; l <= log.level; l++ {
		if err = lw.files[l].write(log.msg); err != nil {
			return
		}
	}
	return
}

// Flush flush log writer
func (lw *logWriter) Flush() {
	for l := lw.level; l < LEVEL_MAX; l++ {
		lw.files[l].flush()
	}
}

// Close close log writer
func (lw *logWriter) Close() {
	for l := lw.level; l < LEVEL_MAX; l++ {
		lw.files[l].close()
		lw.files[l] = nil
	}
}

// ResetLevel reset log level
func (lw *logWriter) ResetLevel(level Level) error {
	var l Level
	for l = LEVEL_MIN; l < level; l++ {
		if lb := lw.files[l]; lb != nil {
			lb.close()
		}
		lw.files[l] = nil
	}
	var err error
	for ; l < LEVEL_MAX; l++ {
		if lw.files[l] == nil {
			if lw.files[l], err = newLogBuffer(lw, l); err != nil {
				return err
			}
		}
	}
	lw.level = level
	return err
}

// createLogFile create log file in the format:level.log.yyyymmdd-HHMMSS.pid
func createLogFile(dir string, filename string) (*os.File, error) {
	t := timeNow()
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
