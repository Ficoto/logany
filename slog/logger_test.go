package slog

import (
	"os"
	"testing"
)

func TestLog(t *testing.T) {
	l := New(os.Stdout, SetLevel(LevelInfo), SetFormatterJson(), SetAddSource(), SetProjectName("logany"))
	var a struct {
		A string
		B float64
	}
	a.A = "t"
	a.B = 2.0
	l.WithField("a", a).WithFields(map[string]any{
		"b": a,
		"c": a,
	}).Info("test")
	l.Infoln("test2")
}
