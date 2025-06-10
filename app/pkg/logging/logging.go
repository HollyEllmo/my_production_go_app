package logging

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var defLogger = NewLogger()

type logger struct {
	*logrus.Logger
}

type Logger interface {
	SetLevel(level logrus.Level)
	GetLevel() logrus.Level
	WithField(key string, value interface{}) *logrus.Entry
	WithFields(fields logrus.Fields) *logrus.Entry
	WithError(err error) *logrus.Entry
	WithContext(ctx context.Context) *logrus.Entry
	WithTime(t time.Time) *logrus.Entry
	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Trace(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
	Println(args ...interface{})
}

func GetLogger() Logger {
	return defLogger
}

// Init инициализирует глобальный логгер с заданным уровнем логирования
func Init(levelStr string) {
	level := parseLogLevel(levelStr)
	defLogger.SetLevel(level)
}

// parseLogLevel парсит строку уровня логирования в logrus.Level
func parseLogLevel(levelStr string) logrus.Level {
	switch strings.ToLower(levelStr) {
	case "trace":
		return logrus.TraceLevel
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn", "warning":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.InfoLevel
	}
}

func NewLogger() Logger {
	logrusLogger := logrus.New()
	logrusLogger.SetLevel(logrus.InfoLevel)
	logrusLogger.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s:%d", filename, f.Line), fmt.Sprintf("%s()", f.Function)
		},
		DisableColors: false,
		FullTimestamp: true,
	}
	logrusLogger.SetOutput(os.Stdout)
	// TODO не работает, починить
	logrusLogger.SetReportCaller(true)

	return &logger{
		Logger: logrusLogger,
	}
}

func (l *logger) SetLevel(level logrus.Level) {
	l.Logger.SetLevel(level)
}

func (l *logger) GetLevel() logrus.Level {
	return l.Logger.GetLevel()
}

func WithField(ctx context.Context, key string, value interface{}) *logrus.Entry {
	return loggerFromContext(ctx).WithField(key, value)
}

func WithFields(ctx context.Context, fields logrus.Fields) *logrus.Entry {
	return loggerFromContext(ctx).WithFields(fields)
}

func WithError(ctx context.Context, err error) *logrus.Entry {
	return loggerFromContext(ctx).WithError(err)
}

func (l *logger) WithContext(ctx context.Context) *logrus.Entry {
	return l.Logger.WithContext(ctx)
}

func WithTime(ctx context.Context, t time.Time) *logrus.Entry {
	return loggerFromContext(ctx).WithTime(t)
}

func Tracef(ctx context.Context, format string, args ...interface{}) {
	loggerFromContext(ctx).Tracef(format, args...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	loggerFromContext(ctx).Debugf(format, args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	loggerFromContext(ctx).Infof(format, args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	loggerFromContext(ctx).Warnf(format, args...)
}

func Warningf(ctx context.Context, format string, args ...interface{}) {
	loggerFromContext(ctx).Warningf(format, args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	loggerFromContext(ctx).Errorf(format, args...)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	loggerFromContext(ctx).Fatalf(format, args...)
}

func Panicf(ctx context.Context, format string, args ...interface{}) {
	loggerFromContext(ctx).Panicf(format, args...)
}

func Trace(ctx context.Context, args ...interface{}) {
	loggerFromContext(ctx).Trace(args...)
}

func Debug(ctx context.Context, args ...interface{}) {
	loggerFromContext(ctx).Debug(args...)
}

func Info(ctx context.Context, args ...interface{}) {
	loggerFromContext(ctx).Info(args...)
}

func Warning(ctx context.Context, args ...interface{}) {
	loggerFromContext(ctx).Warning(args...)
}

func Error(ctx context.Context, args ...interface{}) {
	loggerFromContext(ctx).Error(args...)
}

func Fatal(ctx context.Context, args ...interface{}) {
	loggerFromContext(ctx).Fatal(args...)
}

func Panic(ctx context.Context, args ...interface{}) {
	loggerFromContext(ctx).Panic(args...)
}

func Println(ctx context.Context, args ...interface{}) {
	loggerFromContext(ctx).Print(args...)
}