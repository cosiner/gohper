package log

//TODO: error handling
//
// Inspired by google glog
import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/cosiner/gohper/lib/defval"
	"github.com/cosiner/gohper/lib/types"
)

const (
	// log dir permission when create
	_LOGDIR_PERM = 0755
)

type (
	// logBuffer represent a log w for a special level
	logBuffer struct {
		file *os.File
		*bufio.Writer
		nbytes  uint64
		logdir  string
		level   string
		bufsize uint64
		maxsize uint64
	}

	FileWriterOption struct {
		Bufsize string
		Maxsize string
		Daily   bool
		Logdir  string
	}

	// logWrite is actuall log w, output is local file
	FileWriter struct {
		level Level
		files []*logBuffer
		lock  int32
		quit  chan struct{}
	}
)

// newLogBuffer create a new log buffer
func newLogBuffer(logdir, level string, bufsize, maxsize uint64) (*logBuffer, error) {
	buf := &logBuffer{
		logdir:  logdir,
		level:   level,
		bufsize: bufsize,
		maxsize: maxsize,
	}
	now := time.Now()
	return buf, buf.newLogFile(&now)
}

// newLogFile create a new log file
func (buf *logBuffer) newLogFile(now *time.Time) (err error) {
	if buf.file != nil {
		buf.Flush()
		buf.file.Close()
	}
	file := fmt.Sprintf("%s.log.%s.%d",
		buf.level, now.Format(timeFormat), os.Getpid())
	if buf.file, err = os.Create(filepath.Join(buf.logdir, file)); err == nil {
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
		now := time.Now()
		if err = buf.newLogFile(&now); err != nil {
			return
		}
	}
	n, err := buf.WriteString(msg)
	buf.nbytes += uint64(n)
	return
}

func (o *FileWriterOption) init() {
	defval.String(&o.Bufsize, "10K")
	defval.String(&o.Maxsize, "20M")
	defval.String(&o.Logdir, "logs")
}

// Config resolv config, format like bufsize=xxx&maxsize=xxx&logdir=xxx&level=info
func (w *FileWriter) Config(conf interface{}) (err error) {
	var opt *FileWriterOption
	if conf == nil {
		opt = &FileWriterOption{}
	} else {
		switch c := conf.(type) {
		case *FileWriterOption:
			opt = c
		case FileWriterOption:
			opt = &c
		default:
			return ErrInvalidConfig
		}
	}
	opt.init()
	w.lock = _UNLOCKED
	w.quit = make(chan struct{})
	if w.level == LEVEL_OFF {
		return
	}

	err = os.MkdirAll(opt.Logdir, _LOGDIR_PERM)
	if err != nil {
		return
	}
	bufsize, err := types.BytesCount(opt.Bufsize)
	if err != nil {
		return
	}
	maxsize, err := types.BytesCount(opt.Maxsize)
	if err != nil {
		return
	}

	w.files = make([]*logBuffer, _LEVEL_MAX+1)
	for l := w.level; l <= _LEVEL_MAX && err == nil; l++ {
		w.files[l], err = newLogBuffer(opt.Logdir, l.String(), bufsize, maxsize)
	}

	if opt.Daily {
		go w.enableDaily()
	}
	return
}

func (w *FileWriter) SetLevel(l Level) {
	w.level = l
}

// Write write log to log file, higher level log will simultaneously
// output to all lower level log file
func (w *FileWriter) Write(log *Log) (err error) {
	w.Lock()
	for l := w.level; l <= log.Level && err == nil; l++ {
		err = log.WriteTo(w.files[l])
	}
	w.Unlock()
	return
}

// Flush flush log w
func (w *FileWriter) Flush() {
	for l := w.level; l <= _LEVEL_MAX; l++ {
		w.files[l].flush()
	}
}

func (w *FileWriter) enableDaily() {
	go func(w *FileWriter) {
		h, m, s := time.Now().Clock()
		sec := 24*3600 - (h*3600 + m*60 + s)

		if sec != 0 {
			select {
			case t := <-time.After(time.Duration(sec) * time.Second):
				w.changeLogFile(&t)
			case <-w.quit:
				return
			}
		}

		ticker := time.NewTicker(24 * time.Hour).C
		for {
			select {
			case t := <-ticker:
				w.changeLogFile(&t)
			case <-w.quit:
				return
			}
		}
	}(w)
}

func (w *FileWriter) changeLogFile(t *time.Time) {
	w.Lock()
	for l := w.level; l <= _LEVEL_MAX; l++ {
		w.files[l].newLogFile(t)
	}
	w.Unlock()
}

// Close close log w
func (w *FileWriter) Close() {
	w.quit <- struct{}{}
	w.Lock()
	for l := w.level; l <= _LEVEL_MAX; l++ {
		w.files[l].close()
		w.files[l] = nil
	}
	w.Unlock()
}

const (
	_UNLOCKED = 0
	_LOCKED   = 1
)

// Spinlock
func (w *FileWriter) Lock() {
	for !atomic.CompareAndSwapInt32(&w.lock, _UNLOCKED, _LOCKED) {
	}
}

func (w *FileWriter) Unlock() {
	for !atomic.CompareAndSwapInt32(&w.lock, _LOCKED, _UNLOCKED) {
	}
}
