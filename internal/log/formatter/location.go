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
	log.RegisterFormatter(log.FormatText, friend.Location{}, LocationTextFormatter{})
}

type LocationTextFormatter struct{}

var _ log.Formatter = (*LocationTextFormatter)(nil)

func (l LocationTextFormatter) FormatSingle(ctx log.FormatterContext, e any) (string, error) {
	location, ok := e.(friend.Location)

	if !ok {
		return "", ErrInvalidEntity
	}

	if ctx.Density == log.DensityCompact {
		return l.formatCompact(location), nil
	}

	return l.formatRegular(location), nil
}

func (l LocationTextFormatter) formatCompact(location friend.Location) string {
	parts := []string{idStyle.Render(location.ID), location.String()}

	if len(location.Tags) > 0 {
		parts = append(parts, tagStyle.Render(lang.RenderTags(location.Tags)))
	}

	return strings.Join(parts, " ") + "\n"
}

func (l LocationTextFormatter) formatRegular(location friend.Location) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s (%s)", labelStyle.Render(location.String()), idStyle.Render(location.ID)))
	sb.WriteString("\n")

	if len(location.Tags) > 0 {
		sb.WriteString("  " + log.BulletChar + " " + tagStyle.Render(lang.RenderTags(location.Tags)) + "\n")
	}

	if location.Desc != "" {
		sb.WriteString("\n")

		wrapped := wrapText(location.Desc, 78)

		for _, line := range wrapped {
			sb.WriteString("  " + line + "\n")
		}
	}

	return sb.String()
}

func (l LocationTextFormatter) FormatList(ctx log.FormatterContext, el any) (string, error) {
	locations, ok := el.([]friend.Location)

	if !ok {
		return "", ErrInvalidEntity
	}

	var buf bytes.Buffer

	w := tabwriter.NewWriter(&buf, 0, 0, 3, ' ', 0)

	for _, loc := range locations {
		if ctx.Density == log.DensityCompact {
			_, _ = fmt.Fprintf(
				w,
				"%s\t%s\n",
				idStyle.Render(loc.ID),
				loc.String(),
			)
		} else {
			_, _ = fmt.Fprintf(
				w,
				"%s\t%s\t%s\n",
				idStyle.Render(loc.ID),
				labelStyle.Render(loc.String()),
				tagStyle.Render(lang.RenderTags(loc.Tags)),
			)
		}
	}

	_ = w.Flush()

	return buf.String(), nil
}
