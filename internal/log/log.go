package log

import (
	"sync/atomic"
)

var (
	logger atomic.Pointer[DefaultLogger]
)

func init() {
	logger.Store(New())
}

func SetLogger(l *DefaultLogger) {
	logger.Store(l)
}

func SetLevel(level LogLevel) {
	if l := logger.Load(); l != nil {
		l.SetLevel(level)
	}
}

func SetOutputHandler(handler OutputHandler) {
	if l := logger.Load(); l != nil {
		_ = l.SetOutputHandler(handler)
	}
}

func Debug(msg string) {
	logger.Load().Debug(msg)
}

func Debugf(format string, args ...any) {
	logger.Load().Debugf(format, args...)
}

func Infof(format string, args ...any) {
	logger.Load().Infof(format, args...)
}

func Info(msg string) {
	logger.Load().Info(msg)
}

func PrintEntity(entity any) {
	logger.Load().PrintEntity(entity)
}

func Warnf(format string, args ...any) {
	logger.Load().Warnf(format, args...)
}

func Warn(msg string) {
	logger.Load().Warn(msg)
}

func Errorf(format string, args ...any) {
	logger.Load().Errorf(format, args...)
}

func Error(msg string) {
	logger.Load().Error(msg)
}
