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
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	
	Traceln(args ...interface{})
	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Println(args ...interface{})
	Warnln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})
}

// func GetLogger(ctx context.Context) Logger {
// 	return loggerFromContext(ctx)
// }

func GetLogger() Logger {
	return defLogger
}

// CallerHook это хук для исправления caller'а в логах
type CallerHook struct{}

func (hook *CallerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *CallerHook) Fire(entry *logrus.Entry) error {
	// Получаем стек вызовов
	pc := make([]uintptr, 15)
	n := runtime.Callers(8, pc) // Начинаем с 8-го фрейма, чтобы пропустить логгер
	frames := runtime.CallersFrames(pc[:n])
	
	// Ищем первый фрейм, который не относится к логгеру или logrus
	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		
		// Пропускаем фреймы из пакета logging и logrus
		if !strings.Contains(frame.File, "/pkg/logging/") && 
		   !strings.Contains(frame.File, "/logrus/") &&
		   !strings.Contains(frame.File, "github.com/sirupsen/logrus") {
			entry.Caller = &runtime.Frame{
				File: frame.File,
				Line: frame.Line,
				Function: frame.Function,
				PC: frame.PC,
			}
			break
		}
	}
	
	return nil
}

func NewLogger() Logger {
	logrusLogger := logrus.New()
	logrusLogger.SetLevel(logrus.InfoLevel)
	logrusLogger.SetReportCaller(true)
	
	// Добавляем хук для исправления caller'а
	logrusLogger.AddHook(&CallerHook{})
	
	logrusLogger.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s:%d", filename, f.Line), ""
		},
		DisableColors:   false,
		FullTimestamp:   true,
	}

	logrusLogger.SetOutput(os.Stdout)

	return &logger{
		Logger: logrusLogger,
	}
}

// SetLevel sets the logging level
func SetLevel(level logrus.Level) {
	defLogger.SetLevel(level)
}

// GetLevel returns the current logging level
func GetLevel() logrus.Level {
	return defLogger.GetLevel()
}

// WithField adds a field to the logger
func WithField(ctx context.Context, key string, value interface{}) *logrus.Entry {
	return LoggerFromContext(ctx).WithField(key, value)
}

// WithFields adds multiple fields to the logger
func WithFields(ctx context.Context, fields logrus.Fields) *logrus.Entry {
	return LoggerFromContext(ctx).WithFields(fields)
}

// WithError adds an error to the logger
func WithError(ctx context.Context, err error) *logrus.Entry {
	return LoggerFromContext(ctx).WithError(err)
}

// WithContext adds context to the logger
func WithContext(ctx context.Context) *logrus.Entry {
	return LoggerFromContext(ctx).WithContext(ctx)
}

// WithTime adds time to the logger
func WithTime(ctx context.Context, t time.Time) *logrus.Entry {
	return LoggerFromContext(ctx).WithTime(t)
}

// Tracef logs a formatted trace message
func Tracef(ctx context.Context, format string, args ...interface{}) {
	LoggerFromContext(ctx).Tracef(format, args...)
}

// Debugf logs a formatted debug message
func Debugf(ctx context.Context, format string, args ...interface{}) {
	LoggerFromContext(ctx).Debugf(format, args...)
}

// Infof logs a formatted info message
func Infof(ctx context.Context, format string, args ...interface{}) {
	LoggerFromContext(ctx).Infof(format, args...)
}

// Printf logs a formatted message
func Printf(ctx context.Context, format string, args ...interface{}) {
	LoggerFromContext(ctx).Printf(format, args...)
}

// Warnf logs a formatted warning message
func Warnf(ctx context.Context, format string, args ...interface{}) {
	LoggerFromContext(ctx).Warnf(format, args...)
}

// Warningf logs a formatted warning message
func Warningf(ctx context.Context, format string, args ...interface{}) {
	LoggerFromContext(ctx).Warningf(format, args...)
}

// Errorf logs a formatted error message
func Errorf(ctx context.Context, format string, args ...interface{}) {
	LoggerFromContext(ctx).Errorf(format, args...)
}

// Fatalf logs a formatted fatal message and exits
func Fatalf(ctx context.Context, format string, args ...interface{}) {
	LoggerFromContext(ctx).Fatalf(format, args...)
}

// Panicf logs a formatted panic message and panics
func Panicf(ctx context.Context, format string, args ...interface{}) {
	LoggerFromContext(ctx).Panicf(format, args...)
}

// Traceln logs a trace message with newline
func Traceln(ctx context.Context, args ...interface{}) {
	LoggerFromContext(ctx).Traceln(args...)
}

// Debugln logs a debug message with newline
func Debugln(ctx context.Context, args ...interface{}) {
	LoggerFromContext(ctx).Debugln(args...)
}

// Infoln logs an info message with newline
func Infoln(ctx context.Context, args ...interface{}) {
	LoggerFromContext(ctx).Infoln(args...)
}

// Println logs a message with newline
func Println(ctx context.Context, args ...interface{}) {
	LoggerFromContext(ctx).Println(args...)
}

// Warnln logs a warning message with newline
func Warnln(ctx context.Context, args ...interface{}) {
	LoggerFromContext(ctx).Warnln(args...)
}

// Errorln logs an error message with newline
func Errorln(ctx context.Context, args ...interface{}) {
	LoggerFromContext(ctx).Errorln(args...)
}

// Fatalln logs a fatal message with newline and exits
func Fatalln(ctx context.Context, args ...interface{}) {
	LoggerFromContext(ctx).Fatalln(args...)
}

// Panicln logs a panic message with newline and panics
func Panicln(ctx context.Context, args ...interface{}) {
	LoggerFromContext(ctx).Panicln(args...)
}






