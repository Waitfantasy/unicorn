package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	Ltime  = 1 << iota //time format "2006/01/02 15:04:05"
	Lfile              //file.go:123
	Llevel             //[Debug|Info...]
)

const (
	LevelDebug = iota
	LevelInfo
	LevelWarning
	LevelError
)

const (
	TimeFormat           = "2006/01/02 15:04:05"
	DefaultLogFilePrefix = "unicorn"
	DefaultLogFileSuffix = "log"
)

var LevelSlice = [4]string{"DEBUG", "INFO", "WARN", "ERROR"}

type Config struct {
	Level      string
	Output     string
	Split      bool
	FilePath   string
	FilePrefix string
	FileSuffix string
}

type Logger struct {
	cfg           *Config
	flag          int
	debug         bool
	info          bool
	warning       bool
	error         bool
	fileOut       bool
	consoleOut    bool
	fileWriter    io.Writer
	consoleWriter io.Writer
}

var GlobalLogger *Logger

func New(cfg *Config) *Logger {
	GlobalLogger = new(Logger)

	GlobalLogger.cfg = cfg

	GlobalLogger.flag = Ltime | Lfile | Llevel

	return GlobalLogger
}

func (l *Logger) Run() error {
	GlobalLogger.setLevel()

	if err := GlobalLogger.setWriter(); err != nil {
		return err
	}

	if GlobalLogger.cfg.Split {
		go GlobalLogger.RunSplit()
		GlobalLogger.Debug("run log split goroutine")
	}

	return nil
}

func (l *Logger) setLevel() {
	if l.cfg.Level == "" {
		l.debug = true
		l.info = true
		l.warning = true
		l.error = true

	} else {
		levels := strings.Split(l.cfg.Level, "|")
		for _, level := range levels {
			switch strings.ToLower(level) {
			case "debug":
				l.debug = true
			case "info":
				l.info = true
			case "warning":
				l.warning = true
			case "error":
				l.error = true
			}
		}
	}
}

func (l *Logger) setWriter() error {
	outputs := strings.Split(l.cfg.Output, "|")
	for _, output := range outputs {
		switch strings.ToLower(output) {
		case "console":
			l.consoleOut = true
			l.consoleWriter = os.Stderr
		case "file":
			l.fileOut = true
			if w, err := l.createFileWriter(formatDate(time.Now())); err != nil {
				return err
			} else {
				l.fileWriter = w
			}
		}
	}
	return nil
}

func (l *Logger) createFileWriter(date string) (io.Writer, error) {
	if l.cfg.FilePath == "" {
		switch runtime.GOOS {
		case "windows":
			l.cfg.FilePath = os.Getenv("userprofile") + "\\.unicorn\\log"
		case "linux":
			l.cfg.FilePath = "/var/log/unicorn"
		default:
			l.cfg.FilePath = "."
		}
	}
	fmt.Println(l.cfg.FilePath)

	// create file path
	if _, err := os.Stat(l.cfg.FilePath); os.IsNotExist(err) {
		if err = os.MkdirAll(l.cfg.FilePath, 0644); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	if l.cfg.FilePrefix == "" {
		l.cfg.FilePrefix = DefaultLogFilePrefix
	}

	if l.cfg.FileSuffix == "" {
		l.cfg.FileSuffix = DefaultLogFileSuffix
	}

	// xxx-201810231955.xxx
	filename := l.cfg.FilePath + "/" + l.cfg.FilePrefix + "-" + date + "." + l.cfg.FileSuffix

	// file exists
	if _, err := os.Stat(filename); err == nil {
		if f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0666); err != nil {
			return nil, err
		} else {
			return f, nil
		}
		// new file
	} else if os.IsNotExist(err) {
		if f, err := os.Create(filename); err != nil {
			return nil, err
		} else {
			return f, nil
		}
	} else {
		return nil, err
	}
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if l.debug {
		l.output(LevelDebug, format, v...)
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	if l.info {
		l.output(LevelInfo, format, v...)
	}
}

func (l *Logger) Warning(format string, v ...interface{}) {
	if l.warning {
		l.output(LevelWarning, format, v...)
	}
}

func (l *Logger) Error(format string, v ...interface{}) {
	if l.error {
		l.output(LevelError, format, v...)
	}
}

func (l *Logger) output(level int, format string, v ...interface{}) {
	var buf bytes.Buffer

	if l.flag&Ltime > 0 {
		now := time.Now().Format(TimeFormat)
		buf.WriteString(now)
		buf.WriteString(" - ")
	}

	if l.flag&Llevel > 0 {
		buf.WriteString(LevelSlice[level])
		buf.WriteString(" - ")
	}

	if l.flag&Lfile > 0 {
		_, file, line, ok := runtime.Caller(2)
		if !ok {
			file = "???"
			line = 0
		} else {
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					file = file[i+1:]
					break
				}
			}
		}
		buf.WriteString(file)
		buf.WriteString(":[")
		buf.WriteString(strconv.FormatInt(int64(line), 10))
		buf.WriteString("] - ")
	}

	s := fmt.Sprintf(format, v...)
	buf.WriteString(s)
	if s[len(s)-1] != '\n' {
		buf.WriteString("\n")
	}

	msg := buf.String()
	if l.consoleWriter != nil {
		fmt.Fprintf(l.consoleWriter, msg)
	}

	if l.fileWriter != nil {
		fmt.Fprintf(l.fileWriter, msg)
	}

}

func (l *Logger) RunSplit() {
	d, _ := getDuration()
	tick := time.Tick(d)
	for {
		select {
		case <-tick:
			w, err := l.createFileWriter(formatDate(time.Now()))
			if err != nil {
				tick = time.Tick(time.Second)
				break
			}

			l.fileWriter = w

			if d, err = getDuration(); err != nil {
				tick = time.Tick(time.Second)
				break
			} else {
				tick = time.Tick(d)
			}
		}
	}
}

func getDuration() (time.Duration, error) {
	t1 := time.Now()
	t2 := t1.AddDate(0, 0, 1)

	t2Zero, err := time.ParseInLocation("2006-01-02", formatDate(t2), time.Local)
	if err != nil {
		return 0, err
	}

	// next day
	return time.Duration(t2Zero.Unix()-t1.Unix()) * time.Second, nil
}

func formatDate(t time.Time) string { return t.Format("2006-01-02") }
