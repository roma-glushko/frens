// Copyright 2025 Roma Hlushko
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

	if l.Level > LogLevelVerbose {
		return
	}

	l.OutputHandler(l.Out, msg)
}

func (l *DefaultLogger) Debugf(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.Level > LogLevelVerbose {
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
