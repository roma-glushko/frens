// Copyright 2026 Roma Hlushko
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
	"fmt"
	"reflect"
	"sync"
)

type Format int

const (
	FormatText Format = iota
	FormatJSON
	FormatMarkdown
)

func (f Format) String() string {
	switch f {
	case FormatText:
		return "text"
	case FormatJSON:
		return "json"
	case FormatMarkdown:
		return "markdown"
	default:
		return "unknown"
	}
}

type Density int

const (
	DensityRegular Density = iota
	DensityCompact
)

func (d Density) String() string {
	switch d {
	case DensityRegular:
		return "regular"
	case DensityCompact:
		return "compact"
	default:
		return "regular"
	}
}

// FormatterContext provides formatting options to formatters
type FormatterContext struct {
	Density Density
}

type formatterKey struct {
	Format Format
	Type   reflect.Type
}

type Formatter interface {
	FormatSingle(ctx FormatterContext, entity any) (string, error)
	FormatList(ctx FormatterContext, entities any) (string, error)
}

var (
	formatterMu sync.RWMutex
	formatters  map[formatterKey]Formatter
)

func RegisterFormatter(f Format, e any, formatter Formatter) {
	if formatter == nil {
		panic("formatter cannot be nil")
	}

	t := reflect.TypeOf(e)

	if t == nil {
		panic("entity type cannot be nil")
	}

	key := formatterKey{Format: f, Type: t}

	formatterMu.Lock()
	defer formatterMu.Unlock()

	if formatters == nil {
		formatters = make(map[formatterKey]Formatter)
	}

	if _, exists := formatters[key]; exists {
		panic("formatter already registered for this type and format")
	}

	formatters[key] = formatter
}

func GetFormatter(f Format, e any) (Formatter, error) {
	t := reflect.TypeOf(e)

	if t == nil {
		panic("entity type cannot be nil")
	}

	key := formatterKey{Format: f, Type: t}

	formatterMu.RLock()
	defer formatterMu.RUnlock()

	formatter, exists := formatters[key]

	if !exists {
		return nil, fmt.Errorf("no formatter registered for %s type and %s format", t.Name(), f)
	}

	return formatter, nil
}
