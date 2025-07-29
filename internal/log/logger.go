package log

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

type LogLevel int

const (
	LogLevelVerbose LogLevel = iota
	LogLevelStandard
	LogLevelQuiet
)

type Logger interface {
	Debugf(format string, args ...interface{})
}

type DefaultLogger struct {
	Level         LogLevel
	OutputHandler OutputHandler
	Out           io.Writer
	Err           io.Writer
	mu            sync.Mutex
}

func New() *DefaultLogger {
	return &DefaultLogger{
		Level:         LogLevelStandard,
		OutputHandler: TextOutputHandler,
		Out:           os.Stdout,
		Err:           os.Stderr,
	}
}

func (l *DefaultLogger) PrintEntity(e any) {
}

func (l *DefaultLogger) PrintList(el []any) {
}

func (l *DefaultLogger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.Level = level
}

func (l *DefaultLogger) SetOutputHandler(h OutputHandler) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if h == nil {
		return errors.New("output handler cannot be nil")
	}

	l.OutputHandler = h

	return nil
}

func (l *DefaultLogger) Debug(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.Level < LogLevelVerbose {
		return
	}

	l.OutputHandler(l.Out, msg)
}

func (l *DefaultLogger) Debugf(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.Level < LogLevelVerbose {
		return
	}

	l.OutputHandler(l.Out, fmt.Sprintf(format, args...))
}

func (l *DefaultLogger) Info(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.Level > LogLevelStandard {
		return
	}

	l.OutputHandler(l.Out, msg)
}

func (l *DefaultLogger) Infof(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.Level > LogLevelStandard {
		return
	}

	l.OutputHandler(l.Out, fmt.Sprintf(format, args...))
}

func (l *DefaultLogger) Warn(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.Level > LogLevelStandard {
		return
	}

	l.OutputHandler(l.Err, msg)
}

func (l *DefaultLogger) Warnf(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.Level > LogLevelStandard {
		return
	}

	l.OutputHandler(l.Err, fmt.Sprintf(format, args...))
}

func (l *DefaultLogger) Error(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.OutputHandler(l.Err, msg)
}

func (l *DefaultLogger) Errorf(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.OutputHandler(l.Err, fmt.Sprintf(format, args...))
}
