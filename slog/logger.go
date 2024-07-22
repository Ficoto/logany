package slog

import (
	"context"
	"fmt"
	"github.com/ficoto/logany"
	"io"
	"log/slog"
	"os"
	"runtime"
)

type Formatter uint8

const (
	FormatterText Formatter = 0
	FormatterJson Formatter = 1
)

func LevelReplaceAttr(groups []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case slog.LevelKey:
		l := a.Value.Any().(slog.Level)
		return slog.String(slog.LevelKey, level2Str(l))
	default:
		return a
	}
}

func AddSourceReplaceAttr(groups []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case slog.SourceKey:
		var pcs [1]uintptr
		// skip [runtime.Callers, this function, this function's caller]
		runtime.Callers(9, pcs[:])
		fs := runtime.CallersFrames([]uintptr{pcs[0]})
		f, _ := fs.Next()
		return slog.Any(slog.SourceKey, &slog.Source{
			Function: f.Function,
			File:     f.File,
			Line:     f.Line,
		})
	default:
		return a
	}
}

func replaceAttrForList(handlers ...func(groups []string, a slog.Attr) slog.Attr) func(groups []string, a slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		for _, handler := range handlers {
			a = handler(groups, a)
		}
		return a
	}
}

type Config struct {
	formatter       Formatter
	level           slog.Level
	projectName     string
	addSource       bool
	replaceAttrList []func(groups []string, a slog.Attr) slog.Attr
}

type Logger struct {
	*slog.Logger
}

func New(writer io.Writer, setters ...Setter) logany.Logger {
	var c = new(Config)
	for _, setter := range setters {
		setter(c)
	}
	var (
		h              slog.Handler
		handlerOptions = &slog.HandlerOptions{
			AddSource: c.addSource,
			Level:     c.level,
		}
	)
	if len(c.replaceAttrList) != 0 {
		handlerOptions.ReplaceAttr = replaceAttrForList(c.replaceAttrList...)
	}
	switch c.formatter {
	case FormatterText:
		h = slog.NewTextHandler(writer, handlerOptions)
	case FormatterJson:
		h = slog.NewJSONHandler(writer, handlerOptions)
	}
	l := slog.New(h)
	if len(c.projectName) != 0 {
		l = l.With(slog.String("project", c.projectName))
	}
	return &Logger{
		Logger: l,
	}
}

func (l *Logger) WithError(err error) logany.Logger {
	return &Logger{
		Logger: l.Logger.With(slog.Any("error", err)),
	}
}

func (l *Logger) WithField(key string, value any) logany.Logger {
	return &Logger{
		Logger: l.Logger.With(slog.Any(key, value)),
	}
}

func (l *Logger) WithFields(fields map[string]any) logany.Logger {
	var (
		attrList = make([]any, len(fields))
		index    int
	)
	for k, v := range fields {
		attrList[index] = slog.Any(k, v)
		index++
	}
	return &Logger{
		Logger: l.Logger.With(attrList...),
	}
}

func (l *Logger) Trace(args ...any) {
	l.Logger.Log(context.Background(), LevelTrace, fmt.Sprint(args...))
}

func (l *Logger) Debug(args ...any) {
	l.Logger.Log(context.Background(), LevelDebug, fmt.Sprint(args...))
}

func (l *Logger) Print(args ...any) {
	l.Info(args...)
}

func (l *Logger) Info(args ...any) {
	l.Logger.Log(context.Background(), LevelInfo, fmt.Sprint(args...))
}

func (l *Logger) Warn(args ...any) {
	l.Logger.Log(context.Background(), LevelWarn, fmt.Sprint(args...))
}

func (l *Logger) Error(args ...any) {
	l.Logger.Log(context.Background(), LevelError, fmt.Sprint(args...))
}

func (l *Logger) Fatal(args ...any) {
	l.Logger.Log(context.Background(), LevelFatal, fmt.Sprint(args...))
	if !l.Logger.Enabled(context.Background(), LevelFatal) {
		return
	}
	os.Exit(1)
}

func (l *Logger) Panic(args ...any) {
	l.Logger.Log(context.Background(), LevelPanic, fmt.Sprint(args...))
	if !l.Logger.Enabled(context.Background(), LevelPanic) {
		return
	}
	panic(fmt.Sprint(args...))
}

func (l *Logger) Tracef(format string, args ...any) {
	l.Logger.Log(context.Background(), LevelTrace, fmt.Sprintf(format, args...))
}

func (l *Logger) Debugf(format string, args ...any) {
	l.Logger.Log(context.Background(), LevelDebug, fmt.Sprintf(format, args...))
}

func (l *Logger) Infof(format string, args ...any) {
	l.Logger.Log(context.Background(), LevelInfo, fmt.Sprintf(format, args...))
}

func (l *Logger) Printf(format string, args ...any) {
	l.Infof(format, args...)
}

func (l *Logger) Warnf(format string, args ...any) {
	l.Logger.Log(context.Background(), LevelWarn, fmt.Sprintf(format, args...))
}

func (l *Logger) Errorf(format string, args ...any) {
	l.Logger.Log(context.Background(), LevelError, fmt.Sprintf(format, args...))
}

func (l *Logger) Fatalf(format string, args ...any) {
	l.Logger.Log(context.Background(), LevelFatal, fmt.Sprintf(format, args...))
	if !l.Logger.Enabled(context.Background(), LevelFatal) {
		return
	}
	os.Exit(1)
}

func (l *Logger) Panicf(format string, args ...any) {
	l.Logger.Log(context.Background(), LevelPanic, fmt.Sprintf(format, args...))
	if !l.Logger.Enabled(context.Background(), LevelPanic) {
		return
	}
	panic(fmt.Sprintf(format, args...))
}

func sprintlnn(args ...interface{}) string {
	msg := fmt.Sprintln(args...)
	return msg[:len(msg)-1]
}

func (l *Logger) Traceln(args ...any) {
	l.Logger.Log(context.Background(), LevelTrace, sprintlnn(args...))
}

func (l *Logger) Debugln(args ...any) {
	l.Logger.Log(context.Background(), LevelDebug, sprintlnn(args...))
}

func (l *Logger) Println(args ...any) {
	l.Infoln(args...)
}

func (l *Logger) Infoln(args ...any) {
	l.Logger.Log(context.Background(), LevelInfo, sprintlnn(args...))
}

func (l *Logger) Warnln(args ...any) {
	l.Logger.Log(context.Background(), LevelWarn, sprintlnn(args...))
}

func (l *Logger) Errorln(args ...any) {
	l.Logger.Log(context.Background(), LevelError, sprintlnn(args...))
}

func (l *Logger) Fatalln(args ...any) {
	l.Logger.Log(context.Background(), LevelFatal, sprintlnn(args...))
	if !l.Logger.Enabled(context.Background(), LevelFatal) {
		return
	}
	os.Exit(1)
}

func (l *Logger) Panicln(args ...any) {
	l.Logger.Log(context.Background(), LevelPanic, sprintlnn(args...))
	if !l.Logger.Enabled(context.Background(), LevelPanic) {
		return
	}
	panic(sprintlnn(args...))
}
