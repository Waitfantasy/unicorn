package logger

import (
	"testing"
	"time"
)

func TestLog(t *testing.T)  {
	l := &Log{}
	l.SetDebugLevel()
	l.SetInfoLevel()
	l.SetWarnLevel()
	l.SetErrLevel()
	l.SetConsoleOut()
	l.SetFileOut("logs")
	l.SetFileSplit(true)
	err := l.InitLogger()
	if err != nil {
		t.Error(err)
		return
	}

	go func() {
		if err := l.Split(); err != nil {
			t.Error(err)
		}
	}()

	go func() {
		tick := time.Tick(time.Second)
		for {
			select {
			case <-tick:
				l.Debug("%s\n", "do some...")
				l.Info("%s\n", "do some...")
				l.Warn("%s\n", "do some...")
				l.Err("%s\n", "do some...")
				tick = time.Tick(time.Second)
			}
		}
	}()
	time.Sleep(time.Hour)
}
