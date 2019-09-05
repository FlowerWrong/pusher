package log

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/FlowerWrong/pusher/config"
	"github.com/sirupsen/logrus"
)

var (
	logger     *logrus.Logger
	loggerOnce sync.Once
)

func init() {
	Logger()
}

func initLogger() {
	logger = logrus.New()
	if config.AppEnv == "production" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		// The TextFormatter is default, you don't actually have to do this.
		logger.SetFormatter(&logrus.TextFormatter{})
	}
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)
}

// Logger ...
func Logger() *logrus.Logger {
	if logger == nil {
		loggerOnce.Do(func() {
			initLogger()
		})
	}
	return logger
}

// Adds a field to the log entry, note that it doesn't log until you call
// Debug, Print, Info, Warn, Error, Fatal or Panic. It only creates a log entry.
// If you want multiple fields, use `WithFields`.
func WithField(key string, value interface{}) *logrus.Entry {
	return logger.WithField(key, value)
}

// Adds a struct of fields to the log entry. All it does is call `WithField` for
// each `Field`.
func WithFields(fields logrus.Fields) *logrus.Entry {
	return logger.WithFields(fields)
}

// Add an error as single field to the log entry.  All it does is call
// `WithError` for the given `error`.
func WithError(err error) *logrus.Entry {
	return logger.WithError(err)
}

// Add a context to the log entry.
func WithContext(ctx context.Context) *logrus.Entry {
	return logger.WithContext(ctx)
}

// Overrides the time of the log entry.
func WithTime(t time.Time) *logrus.Entry {
	return logger.WithTime(t)
}

func Logf(level logrus.Level, format string, args ...interface{}) {
	logger.Logf(level, format, args...)
}

func Tracef(format string, args ...interface{}) {
	logger.Logf(logrus.TraceLevel, format, args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Logf(logrus.DebugLevel, format, args...)
}

func Infof(format string, args ...interface{}) {
	logger.Logf(logrus.InfoLevel, format, args...)
}

func Printf(format string, args ...interface{}) {
	logger.Printf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Logf(logrus.WarnLevel, format, args...)
}

func Warningf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Logf(logrus.ErrorLevel, format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Logf(logrus.FatalLevel, format, args...)
	logger.Exit(1)
}

func Panicf(format string, args ...interface{}) {
	logger.Logf(logrus.PanicLevel, format, args...)
}

func Log(level logrus.Level, args ...interface{}) {
	logger.Log(level, args...)
}

func Trace(args ...interface{}) {
	logger.Log(logrus.TraceLevel, args...)
}

func Debug(args ...interface{}) {
	logger.Log(logrus.DebugLevel, args...)
}

func Info(args ...interface{}) {
	logger.Log(logrus.InfoLevel, args...)
}

func Print(args ...interface{}) {
	logger.Print(args...)
}

func Warn(args ...interface{}) {
	logger.Log(logrus.WarnLevel, args...)
}

func Warning(args ...interface{}) {
	logger.Warn(args...)
}

func Error(args ...interface{}) {
	logger.Log(logrus.ErrorLevel, args...)
}

func Fatal(args ...interface{}) {
	logger.Log(logrus.FatalLevel, args...)
	logger.Exit(1)
}

func Panic(args ...interface{}) {
	logger.Log(logrus.PanicLevel, args...)
}

func Logln(level logrus.Level, args ...interface{}) {
	logger.Logln(level, args...)
}

func Traceln(args ...interface{}) {
	logger.Logln(logrus.TraceLevel, args...)
}

func Debugln(args ...interface{}) {
	logger.Logln(logrus.DebugLevel, args...)
}

func Infoln(args ...interface{}) {
	logger.Logln(logrus.InfoLevel, args...)
}

func Println(args ...interface{}) {
	logger.Println(args...)
}

func Warnln(args ...interface{}) {
	logger.Logln(logrus.WarnLevel, args...)
}

func Warningln(args ...interface{}) {
	logger.Warnln(args...)
}

func Errorln(args ...interface{}) {
	logger.Logln(logrus.ErrorLevel, args...)
}

func Fatalln(args ...interface{}) {
	logger.Logln(logrus.FatalLevel, args...)
	logger.Exit(1)
}

func Panicln(args ...interface{}) {
	logger.Logln(logrus.PanicLevel, args...)
}

func Exit(code int) {
	logger.Exit(code)
}
