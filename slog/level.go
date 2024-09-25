package slog

import (
	"github.com/ficoto/logany"
	"log/slog"
)

const (
	LevelTrace            = -8
	LevelDebug            = slog.LevelDebug
	LevelInfo             = slog.LevelInfo
	LevelWarn             = slog.LevelWarn
	LevelError            = slog.LevelError
	LevelFatal slog.Level = 12
	LevelPanic slog.Level = 16
)

func toSlogLevel(l logany.Level) slog.Level {
	return slog.Level(l)
}

func level2Str(l slog.Level) string {
	switch l {
	case LevelFatal:
		return "FATAL"
	case LevelPanic:
		return "PANIC"
	case LevelTrace:
		return "TRACE"
	default:
		return l.String()
	}
}
