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
	"io"
	"reflect"
)

// Printer provides a high-level API for printing entities in the configured format.
// It abstracts away format selection so commands don't need to know which format is active.
type Printer interface {
	Print(entity any) error
	PrintList(entities any) error
}

type printer struct {
	format  Format
	density Density
	w       io.Writer
}

// NewPrinter creates a Printer that outputs entities in the specified format with regular density.
func NewPrinter(format Format, w io.Writer) Printer {
	return &printer{format: format, density: DensityRegular, w: w}
}

// NewPrinterWithDensity creates a Printer with the specified format and density.
func NewPrinterWithDensity(format Format, density Density, w io.Writer) Printer {
	return &printer{format: format, density: density, w: w}
}

func (p *printer) Print(entity any) error {
	fmtr, err := GetFormatter(p.format, entity)
	if err != nil {
		return err
	}

	ctx := FormatterContext{Density: p.density}
	out, err := fmtr.FormatSingle(ctx, entity)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(p.w, out)

	return err
}

func (p *printer) PrintList(entities any) error {
	// Get element type from slice to look up the correct formatter
	t := reflect.TypeOf(entities)
	if t.Kind() != reflect.Slice {
		return fmt.Errorf("PrintList expects a slice, got %T", entities)
	}

	// Create a zero value of the element type for formatter lookup
	elemType := t.Elem()
	sample := reflect.New(elemType).Elem().Interface()

	fmtr, err := GetFormatter(p.format, sample)
	if err != nil {
		return err
	}

	ctx := FormatterContext{Density: p.density}
	out, err := fmtr.FormatList(ctx, entities)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(p.w, out)

	return err
}
