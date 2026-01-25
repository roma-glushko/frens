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

package formatter

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/lang"
	"github.com/roma-glushko/frens/internal/log"
)

func init() {
	log.RegisterFormatter(log.FormatText, friend.Date{}, DateTextFormatter{})
}

type DateTextFormatter struct{}

var _ log.Formatter = (*DateTextFormatter)(nil)

func (f DateTextFormatter) FormatSingle(ctx log.FormatterContext, e any) (string, error) {
	var dt *friend.Date

	switch v := e.(type) {
	case friend.Date:
		dt = &v
	case *friend.Date:
		dt = v
	default:
		return "", ErrInvalidEntity
	}

	if ctx.Density == log.DensityCompact {
		return f.formatCompact(dt), nil
	}

	return f.formatRegular(dt), nil
}

func (f DateTextFormatter) formatCompact(dt *friend.Date) string {
	parts := []string{idStyle.Render(dt.ID), dt.DateExpr}

	if dt.Person != "" {
		parts = append(parts, dt.Person)
	}

	if len(dt.Tags) > 0 {
		parts = append(parts, tagStyle.Render(lang.RenderTags(dt.Tags)))
	}

	return strings.Join(parts, " ") + "\n"
}

func (f DateTextFormatter) formatRegular(dt *friend.Date) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("[%s] %s", idStyle.Render(dt.ID), labelStyle.Render(dt.DateExpr)))

	if len(dt.Tags) > 0 {
		sb.WriteString("\n")
		sb.WriteString(" * " + tagStyle.Render(lang.RenderTags(dt.Tags)))
		sb.WriteString(" ")
	}

	if dt.Desc != "" {
		sb.WriteString("\n")

		wrapped := wrapText(dt.Desc, 80)

		for _, line := range wrapped {
			sb.WriteString(" " + line + "\n")
		}
	}

	return sb.String()
}

func (f DateTextFormatter) FormatList(ctx log.FormatterContext, el any) (string, error) {
	dates, ok := el.([]*friend.Date)

	if !ok {
		return "", ErrInvalidEntity
	}

	var buf bytes.Buffer

	w := tabwriter.NewWriter(&buf, 0, 0, 3, ' ', 0)

	for _, dt := range dates {
		if ctx.Density == log.DensityCompact {
			_, _ = fmt.Fprintf(
				w,
				"%s\t%s\t%s\n",
				idStyle.Render(dt.ID),
				dt.Person,
				dt.DateExpr,
			)
		} else {
			_, _ = fmt.Fprintf(
				w,
				" %s\t%s\t%s\t%s\n",
				idStyle.Render(dt.ID),
				labelStyle.Render(dt.Person),
				labelStyle.Render(dt.DateExpr),
				tagStyle.Render(lang.RenderTags(dt.Tags)),
			)
		}
	}

	_ = w.Flush()

	return buf.String(), nil
}
