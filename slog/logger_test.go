package slog

import (
	"github.com/ficoto/logany"
	"os"
	"testing"
)

func TestLog(t *testing.T) {
	l := New(os.Stdout, SetLevel(logany.LevelInfo), SetFormatterJson(), SetAddSource(), SetProjectName("logany"))
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
	l.SetLevel(logany.LevelError)
	l.Info("test3")
	l.Error("test4")
}
