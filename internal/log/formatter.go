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

type formatterKey struct {
	Format Format
	Type   reflect.Type
}

type Formatter interface {
	FormatSingle(entity any) (string, error)
	FormatList(entities []any) (string, error)
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
