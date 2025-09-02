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

package formatter

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/lang"
	"github.com/roma-glushko/frens/internal/log"
)

var ErrInvalidEntity = errors.New("invalid entity")

func init() {
	log.RegisterFormatter(log.FormatText, friend.Person{}, PersonTextFormatter{})
}

func wrapText(text string, width int) []string {
	words := strings.Fields(text)

	var lines []string

	var current string

	for _, word := range words {
		if len(current)+len(word)+1 > width {
			lines = append(lines, current)
			current = word
		} else {
			if current != "" {
				current += " "
			}

			current += word
		}
	}

	if current != "" {
		lines = append(lines, current)
	}

	return lines
}

type PersonTextFormatter struct{}

var _ log.Formatter = (*PersonTextFormatter)(nil)

func (p PersonTextFormatter) FormatSingle(e any) (string, error) {
	person, ok := e.(friend.Person)

	if !ok {
		return "", ErrInvalidEntity
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(" %s [%s]", labelStyle.Render(person.String()), person.ID))
	sb.WriteString("\n")

	if len(person.Tags) > 0 {
		sb.WriteString(" • " + tagStyle.Render(lang.RenderTags(person.Tags)))
		sb.WriteString(" ")
	}

	if len(person.Locations) > 0 {
		sb.WriteString(locationStyle.Render(lang.RenderLocMarkers(person.Locations)))
	}

	if person.Desc != "" {
		sb.WriteString("\n")

		wrapped := wrapText(person.Desc, 80)

		for _, line := range wrapped {
			sb.WriteString(" " + line + "\n")
		}

		sb.WriteString("\n")
	}

	return sb.String(), nil
}

func (p PersonTextFormatter) FormatList(el any) (string, error) {
	persons, ok := el.([]friend.Person)

	if !ok {
		return "", ErrInvalidEntity
	}

	var buf bytes.Buffer

	w := tabwriter.NewWriter(&buf, 0, 0, 3, ' ', 0)

	for _, person := range persons {
		_, _ = fmt.Fprintf(
			w,
			" %s\t%s\t%s\t%s\t%s\t%s\n",
			idStyle.Render(person.ID),
			labelStyle.Render(person.String()),
			tagStyle.Render(lang.RenderTags(person.Tags)),
			locationStyle.Render(lang.RenderLocMarkers(person.Locations)),
			countLabel.Render(fmt.Sprintf("✎ %d", person.Notes)),
			countLabel.Render(fmt.Sprintf("⚙ %d", person.Activities)),
		)
	}

	_ = w.Flush()

	return buf.String(), nil
}
