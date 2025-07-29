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
	"sync/atomic"
)

var logger atomic.Pointer[DefaultLogger]

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
