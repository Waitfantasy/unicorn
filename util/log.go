package util

import (
	"io"
	"log"
	"os"
)

const (
	LDebug = iota
	LInfo
	LWarn
	LErr
)

const (
	LOG_PREFIX_DEBUG   = "Debug:"
	LOG_PREFIX_INFO    = "Info:"
	LOG_PREFIX_WARNING = "Warning:"
	LOG_PREFIX_ERROR   = "Error:"
)

type Log struct {
	debug bool
	info  bool
	warn  bool
	err   bool
	path  string

	fileOut    bool
	consoleOut bool
	writer     io.Writer
}

func (l *Log) SetDebugLevel() {
	l.debug = true
}

func (l *Log) SetInfoLevel() {
	l.info = true
}

func (l *Log) SetWarnLevel() {
	l.warn = true
}

func (l *Log) SetErrLevel() {
	l.err = true
}

func (l *Log) SetConsoleOut() {
	l.consoleOut = true
}

func (l *Log) SetFileOut(path string) {
	l.path = path
	l.fileOut = true
}

func NewLogger(w io.Writer) *Log {
	l := new(Log)
	l.writer = w
	return l
}

func (l *Log) SetWriter(w io.Writer) {
	l.writer = w
}

func (l *Log) SetWriterByFile(filename string) error {
	if filename == "" {
		return nil
	}
	_, err := os.Stat(filename)
	if err == nil {
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		l.writer = f
		return err
	} else {
		if os.IsNotExist(err) {
			if fd, err := os.Create(filename); err != nil {
				return err
			} else {
				l.writer = fd
				return nil
			}
		} else if os.IsExist(err) {
			f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				return err
			}
			l.writer = f
			return nil
		} else {
			return err
		}
	}
}

func (l *Log) Printf(level int, format string, v ...interface{}) {
	l.createLevelLog(level).Printf(format, v...)
}

func (l *Log) Print(level int, v ...interface{}) {
	l.createLevelLog(level).Print(v...)
}

func (l *Log) Debug(format string, v ...interface{}) {
	l.createLevelLog(LDebug).Printf(format, v...)
}

func (l *Log) Info(format string, v ...interface{}) {
	l.createLevelLog(LInfo).Printf(format, v...)
}

func (l *Log) Warning(format string, v ...interface{}) {
	l.createLevelLog(LWarn).Printf(format, v...)
}

func (l *Log) Error(format string, v ...interface{}) {
	l.createLevelLog(LErr).Printf(format, v...)
	os.Exit(1)
}

func (l *Log) createLevelLog(level int) *log.Logger {
	flag := log.Ldate | log.Ltime | log.Lshortfile
	prefix := ""
	switch level {
	case LDebug:
		if l.debug {
			//prefix = LOG_PREFIX_DEBUG
			//return log.New(l.writer, LOG_PREFIX_DEBUG, flag)
		}
	case LInfo:
		if l.info {

		}
		//prefix = LOG_PREFIX_INFO
	case LWarn:
		if l.warn {

		}
		//prefix = LOG_PREFIX_WARNING
	case LErr:
		if l.err {

		}
		prefix = LOG_PREFIX_ERROR
	default:
		prefix = "yc-snowflake"
	}
	return log.New(l.writer, prefix, flag)
}

func (l *Log) createWriters() []io.Writer  {
	var writers []io.Writer
	if l.consoleOut {
		writers = append(writers, os.Stdout)
	}

	if l.fileOut {

	}
}

func (l *Log) createFileWriter() (io.Writer, error) {

}