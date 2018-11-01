package logger

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"runtime"
	"time"
)

type Level int

const (
	LDebug       Level = 1
	LInfo        Level = 2
	LWarn        Level = 3
	LErr         Level = 4
	LDebugPrefix       = "[DEBUG] "
	LInfoPrefix        = "[INFO] "
	LWarnPrefix        = "[WARN] "
	LErrPrefix         = "[ERR] "
)

const (
	DefaultLogFilePrefix   = "unicorn"
	DefaultLogFileSuffix   = "log"
	UnixDefaultLogFilePath = "/var/log/unicorn"
)

type Log struct {
	debug         bool
	info          bool
	warn          bool
	err           bool
	fileOut       bool
	filePath      string
	filePrefix    string
	fileSuffix    string
	fileSplit     bool
	fileLogger    *log.Logger
	consoleOut    bool
	consoleLogger *log.Logger
}

func (l *Log) InitLogger() (error) {
	flag := log.LstdFlags
	if l.consoleOut {
		l.consoleLogger = log.New(os.Stdout, "", flag)
	}

	if l.fileOut {
		if w, err := l.createFileWriter(time.Now().Format("2006-01-02")); err != nil {
			return err
		} else {
			l.fileLogger = log.New(w, "", flag)
		}
	}

	return nil
}

func (l *Log) SetDebugLevel() *Log {
	l.debug = true
	return l
}

func (l *Log) SetInfoLevel() *Log {
	l.info = true
	return l
}

func (l *Log) SetWarnLevel() *Log {
	l.warn = true
	return l
}

func (l *Log) SetErrLevel() *Log {
	l.err = true
	return l
}

func (l *Log) SetConsoleOut() *Log {
	l.consoleOut = true
	return l
}

func (l *Log) SetFileOut(path string) *Log {
	l.filePath = path
	l.fileOut = true
	return l
}

func (l *Log) SetFilePrefix(prefix string) *Log {
	l.filePrefix = prefix
	return l
}

func (l *Log) SetFileSuffix(suffix string) *Log {
	l.fileSuffix = suffix
	return l
}

func (l *Log) SetFileSplit(split bool) *Log {
	l.fileSplit = split
	return l
}

func (l *Log) GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()
		latency := end.Sub(start)
		l.Info("| %3d | %9v | %10s | %s  %s |",
			c.Writer.Status(),
			latency,
			c.ClientIP(),
			c.Request.Method,
			c.Request.URL.Path)
	}
}

func (l *Log) Split() error {
	d, err := createDayDuration()
	if err != nil {
		return err
	}
	tick := time.NewTicker(d)
	for {
		select {
		case <-tick.C:
			if w, err := l.createFileWriter(time.Now().Format("2006-01-02")); err != nil {
				break
			} else {
				l.fileLogger.SetOutput(w)
				if d, err = createDayDuration(); err != nil {
					break
				}
			}
		}
	}
}

func (l *Log) createFileWriter(date string) (io.Writer, error) {
	if l.filePath == "" {
		l.filePath = UnixDefaultLogFilePath
	}

	// create file path
	if _, err := os.Stat(l.filePath); os.IsNotExist(err) {
		if err = os.Mkdir(l.filePath, 0644); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	if l.filePrefix == "" {
		l.filePrefix = DefaultLogFilePrefix
	}

	if l.fileSuffix == "" {
		l.fileSuffix = DefaultLogFileSuffix
	}

	// xxx-2018-10-23.xxx
	filename := l.filePath + "/" + l.filePrefix + "-" + date + "." + l.fileSuffix

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

func createDayDuration() (time.Duration, error) {
	// now
	t1 := time.Now()
	t1Str := t1.Format("2006-01-02")
	// now zero
	t2, err := time.ParseInLocation("2006-01-02", t1Str, time.Local)
	if err != nil {
		return 0, err
	}
	// next day
	t3 := time.Date(t2.Year(), t2.Month(), t2.Day()+1, t2.Hour(), t2.Minute(), t2.Second(), t2.Nanosecond(), time.Local)
	return time.Duration(t3.Unix()-t2.Unix()) * time.Second, nil
}

func (l *Log) Debug(format string, v ...interface{}) {
	l.printf(LDebug, format, v...)
}

func (l *Log) Info(format string, v ...interface{}) {
	l.printf(LInfo, format, v...)
}

func (l *Log) Warn(format string, v ...interface{}) {
	l.printf(LWarn, format, v...)
}

func (l *Log) Err(format string, v ...interface{}) {
	l.printf(LErr, format, v...)
}

func (l *Log) Printf(level Level, format string, v ...interface{}) {
	l.printf(level, format, v...)
}

func (l *Log) printf(level Level, format string, v ...interface{}) {
	if _, file, line, ok := runtime.Caller(2); ok {
		format = "%s:%d: " + format
		v = append([]interface{}{
			file, line,
		}, v...)
	}

	if prefix := l.getPrefix(level); prefix != "" {
		if l.consoleLogger != nil {
			l.consoleLogger.SetPrefix(prefix)
			l.consoleLogger.Printf(format, v...)
		}

		if l.fileLogger != nil {
			l.fileLogger.SetPrefix(prefix)
			l.fileLogger.Printf(format, v...)
		}
	}
}

func (l *Log) getPrefix(level Level) string {
	var prefix string
	switch level {
	case LDebug:
		if l.debug {
			prefix = LDebugPrefix
		}
	case LInfo:
		if l.info {
			prefix = LInfoPrefix
		}
	case LWarn:
		if l.warn {
			prefix = LWarnPrefix
		}
	case LErr:
		if l.err {
			prefix = LErrPrefix
		}
	}

	return prefix
}
