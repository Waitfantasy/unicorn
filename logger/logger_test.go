package logger

import "testing"

func TestLogger(t *testing.T) {
	l := New(&Config{
		Level:  "debug|info|error|warning",
		Output: "console|file",
		Split:  true,
	})
	if err := l.Run(); err != nil {
		t.Error("error")
		return
	}
	l.Debug("debug %s...", "do something")
	l.Info("info %s...", "do something")
	l.Warning("warning %s...", "do something")
	l.Error("error %s...", "do something")
}
