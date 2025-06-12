package logging

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

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

func GetLogger(ctx context.Context) Logger {
	return loggerFromContext(ctx)
}

func NewLogger() Logger {
	logrusLogger := logrus.New()
	logrusLogger.SetLevel(logrus.InfoLevel)
	logrusLogger.SetReportCaller(true)
	logrusLogger.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s:%d", filename, f.Line), fmt.Sprintf("%s()", f.Function)
		},
		DisableColors:   false,
		FullTimestamp:   true,
	}

	logrusLogger.SetOutput(os.Stdout)
	logrusLogger.SetReportCaller(true)

	 return&logger{
		Logger: logrusLogger,
	}

	
}

// SetLevel sets the logging level
func (l *logger) SetLevel(level logrus.Level) {
	l.Logger.SetLevel(level)
}

// GetLevel returns the current logging level
func (l *logger) GetLevel() logrus.Level {
	return l.Logger.GetLevel()
}

// WithField adds a field to the logger
func (l *logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.Logger.WithField(key, value)
}

// WithFields adds multiple fields to the logger
func (l *logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.Logger.WithFields(fields)
}

// WithError adds an error to the logger
func (l *logger) WithError(err error) *logrus.Entry {
	return l.Logger.WithError(err)
}

// WithContext adds context to the logger
func (l *logger) WithContext(ctx context.Context) *logrus.Entry {
	return l.Logger.WithContext(ctx)
}

// WithTime adds time to the logger
func (l *logger) WithTime(t time.Time) *logrus.Entry {
	return l.Logger.WithTime(t)
}

// Tracef logs a formatted trace message
func (l *logger) Tracef(format string, args ...interface{}) {
	l.Logger.Tracef(format, args...)
}

// Debugf logs a formatted debug message
func (l *logger) Debugf(format string, args ...interface{}) {
	l.Logger.Debugf(format, args...)
}

// Infof logs a formatted info message
func (l *logger) Infof(format string, args ...interface{}) {
	l.Logger.Infof(format, args...)
}

// Printf logs a formatted message
func (l *logger) Printf(format string, args ...interface{}) {
	l.Logger.Printf(format, args...)
}

// Warnf logs a formatted warning message
func (l *logger) Warnf(format string, args ...interface{}) {
	l.Logger.Warnf(format, args...)
}

// Warningf logs a formatted warning message
func (l *logger) Warningf(format string, args ...interface{}) {
	l.Logger.Warningf(format, args...)
}

// Errorf logs a formatted error message
func (l *logger) Errorf(format string, args ...interface{}) {
	l.Logger.Errorf(format, args...)
}

// Fatalf logs a formatted fatal message and exits
func (l *logger) Fatalf(format string, args ...interface{}) {
	l.Logger.Fatalf(format, args...)
}

// Panicf logs a formatted panic message and panics
func (l *logger) Panicf(format string, args ...interface{}) {
	l.Logger.Panicf(format, args...)
}

// Traceln logs a trace message with newline
func (l *logger) Traceln(args ...interface{}) {
	l.Logger.Traceln(args...)
}

// Debugln logs a debug message with newline
func (l *logger) Debugln(args ...interface{}) {
	l.Logger.Debugln(args...)
}

// Infoln logs an info message with newline
func (l *logger) Infoln(args ...interface{}) {
	l.Logger.Infoln(args...)
}

// Println logs a message with newline
func (l *logger) Println(args ...interface{}) {
	l.Logger.Println(args...)
}

// Warnln logs a warning message with newline
func (l *logger) Warnln(args ...interface{}) {
	l.Logger.Warnln(args...)
}

// Errorln logs an error message with newline
func (l *logger) Errorln(args ...interface{}) {
	l.Logger.Errorln(args...)
}

// Fatalln logs a fatal message with newline and exits
func (l *logger) Fatalln(args ...interface{}) {
	l.Logger.Fatalln(args...)
}

// Panicln logs a panic message with newline and panics
func (l *logger) Panicln(args ...interface{}) {
	l.Logger.Panicln(args...)
}






