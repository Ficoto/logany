package logruslog

import (
	"github.com/ficoto/logany"
	"github.com/sirupsen/logrus"
	"io"
	"runtime"
	"strings"
	"sync"
)

type AddSourceHook struct{}

func (h AddSourceHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
		logrus.TraceLevel,
	}
}

var (

	// qualified package name, cached at first use
	packageName string

	// Positions in the call stack when tracing to report the calling method
	minimumCallerDepth int

	// Used for caller information initialisation
	callerInitOnce sync.Once
)

const (
	maximumCallerDepth int = 25
	knownLogrusFrames  int = 9
)

func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}

func getCaller() *runtime.Frame {
	// cache this package's fully-qualified name
	callerInitOnce.Do(func() {
		pcs := make([]uintptr, maximumCallerDepth)
		_ = runtime.Callers(0, pcs)

		// dynamic get the package name and the minimum caller depth
		for i := 0; i < maximumCallerDepth; i++ {
			funcName := runtime.FuncForPC(pcs[i]).Name()
			if strings.Contains(funcName, "getCaller") {
				packageName = getPackageName(funcName)
				break
			}
		}

		minimumCallerDepth = knownLogrusFrames

	})

	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)

		// If the caller isn't part of this package, we're done
		if pkg != packageName {
			return &f //nolint:scopelint
		}
	}

	// if we got here, we failed to find the caller's context
	return nil
}

func (h AddSourceHook) Fire(e *logrus.Entry) error {
	e.Caller = getCaller()
	return nil
}

type Config struct {
	formatter           logrus.Formatter
	writer              io.Writer
	level               logrus.Level
	projectName         string
	includeReportCaller bool
	hooks               []logrus.Hook
}

type Logrus struct {
	*logrus.Entry
}

func New(writer io.Writer, setters ...Setter) logany.Logger {
	var config Config
	for _, setter := range setters {
		setter(&config)
	}
	if config.formatter == nil {
		config.formatter = &logrus.JSONFormatter{}
	}
	logrus.SetFormatter(config.formatter)
	logrus.SetOutput(writer)
	logrus.SetLevel(config.level)
	if config.includeReportCaller {
		logrus.SetReportCaller(config.includeReportCaller)
		logrus.AddHook(AddSourceHook{})
	}
	for _, hook := range config.hooks {
		logrus.AddHook(hook)
	}
	return &Logrus{
		Entry: logrus.WithField("project", config.projectName),
	}
}

func (l *Logrus) WithError(err error) logany.Logger {
	return &Logrus{
		Entry: l.Entry.WithError(err),
	}
}

func (l *Logrus) WithField(key string, value any) logany.Logger {
	return &Logrus{
		Entry: l.Entry.WithField(key, value),
	}
}

func (l *Logrus) WithFields(fields map[string]any) logany.Logger {
	return &Logrus{
		Entry: l.Entry.WithFields(fields),
	}
}

func (l *Logrus) Trace(args ...any) {
	l.Entry.Trace(args...)
}

func (l *Logrus) Debug(args ...any) {
	l.Entry.Debug(args...)
}

func (l *Logrus) Print(args ...any) {
	l.Entry.Print(args...)
}

func (l *Logrus) Info(args ...any) {
	l.Entry.Info(args...)
}

func (l *Logrus) Warn(args ...any) {
	l.Entry.Warn(args...)
}

func (l *Logrus) Error(args ...any) {
	l.Entry.Error(args...)
}

func (l *Logrus) Fatal(args ...any) {
	l.Entry.Fatal(args...)
}

func (l *Logrus) Panic(args ...any) {
	l.Entry.Panic(args...)
}

func (l *Logrus) Tracef(format string, args ...any) {
	l.Entry.Tracef(format, args...)
}

func (l *Logrus) Debugf(format string, args ...any) {
	l.Entry.Debugf(format, args...)
}

func (l *Logrus) Infof(format string, args ...any) {
	l.Entry.Infof(format, args...)
}

func (l *Logrus) Printf(format string, args ...any) {
	l.Entry.Printf(format, args...)
}

func (l *Logrus) Warnf(format string, args ...any) {
	l.Entry.Warnf(format, args...)
}

func (l *Logrus) Errorf(format string, args ...any) {
	l.Entry.Errorf(format, args...)
}

func (l *Logrus) Fatalf(format string, args ...any) {
	l.Entry.Fatalf(format, args...)
}

func (l *Logrus) Panicf(format string, args ...any) {
	l.Entry.Panicf(format, args...)
}

func (l *Logrus) Traceln(args ...any) {
	l.Entry.Traceln(args...)
}

func (l *Logrus) Debugln(args ...any) {
	l.Entry.Debugln(args...)
}

func (l *Logrus) Println(args ...any) {
	l.Entry.Println(args...)
}

func (l *Logrus) Infoln(args ...any) {
	l.Entry.Infoln(args...)
}

func (l *Logrus) Warnln(args ...any) {
	l.Entry.Warnln(args...)
}

func (l *Logrus) Errorln(args ...any) {
	l.Entry.Errorln(args...)
}

func (l *Logrus) Fatalln(args ...any) {
	l.Entry.Fatalln(args...)
}

func (l *Logrus) Panicln(args ...any) {
	l.Entry.Panicln(args...)
}
