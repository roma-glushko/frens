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
	log.RegisterFormatter(log.FormatText, friend.Event{}, EventTextFormatter{})
}

func CutStr(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}

	return strings.TrimSpace(s[:maxLen]) + "…"
}

type EventTextFormatter struct{}

var _ log.Formatter = (*EventTextFormatter)(nil)

func (f EventTextFormatter) FormatSingle(ctx log.FormatterContext, entity any) (string, error) {
	e, ok := entity.(friend.Event)

	if !ok {
		return "", ErrInvalidEntity
	}

	if ctx.Density == log.DensityCompact {
		return f.formatCompact(e), nil
	}

	return f.formatRegular(e), nil
}

func (f EventTextFormatter) formatCompact(e friend.Event) string {
	date := e.Date.Format("2006-01-02 15:04")
	parts := []string{idStyle.Render(e.ID), date, CutStr(e.Desc, 50)}

	if len(e.FriendIDs) > 0 {
		parts = append(parts, friendStyle.Render(strings.Join(e.FriendIDs, " ")))
	}

	if len(e.Tags) > 0 {
		parts = append(parts, tagStyle.Render(lang.RenderTags(e.Tags)))
	}

	if len(e.LocationIDs) > 0 {
		parts = append(parts, locationStyle.Render(lang.RenderLocMarkers(e.LocationIDs)))
	}

	return strings.Join(parts, " ") + "\n"
}

func (f EventTextFormatter) formatRegular(e friend.Event) string {
	date := e.Date.Format("Mon Jan 2, 2006 15:04 MST")

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(" %s [%s]\n", labelStyle.Render(date), e.ID))

	if len(e.Tags) > 0 || len(e.LocationIDs) > 0 {
		sb.WriteString(" * ")
	}

	if len(e.Tags) > 0 {
		sb.WriteString(tagStyle.Render(lang.RenderTags(e.Tags)))
		sb.WriteString(" ")
	}

	if len(e.LocationIDs) > 0 {
		sb.WriteString(locationStyle.Render(lang.RenderLocMarkers(e.LocationIDs)))
	}

	if len(e.FriendIDs) > 0 {
		sb.WriteString("\n")
		sb.WriteString(" + " + strings.Join(e.FriendIDs, " "))
	}

	sb.WriteString("\n\n")
	sb.WriteString(" " + e.Desc)

	return sb.String()
}

func (f EventTextFormatter) FormatList(ctx log.FormatterContext, el any) (string, error) {
	events, ok := el.([]friend.Event)

	if !ok {
		return "", ErrInvalidEntity
	}

	var buf bytes.Buffer

	w := tabwriter.NewWriter(&buf, 0, 0, 3, ' ', 0)

	for _, e := range events {
		if ctx.Density == log.DensityCompact {
			_, _ = fmt.Fprintf(
				w,
				"%s\t%s\t%s\t%s\n",
				idStyle.Render(e.ID),
				e.Date.Format("2006-01-02"),
				CutStr(e.Desc, 50),
				friendStyle.Render(strings.Join(e.FriendIDs, " ")),
			)
		} else {
			_, _ = fmt.Fprintf(
				w,
				" %s\t%s\t%s\t%s\t%s\t%s\n",
				idStyle.Render(e.ID),
				e.Date.Format("Mon Jan 2, 2006 15:04 MST"),
				labelStyle.Render(CutStr(e.Desc, 80)),
				friendStyle.Render(strings.Join(e.FriendIDs, " ")),
				tagStyle.Render(lang.RenderTags(e.Tags)),
				locationStyle.Render(lang.RenderLocMarkers(e.LocationIDs)),
			)
		}
	}

	_ = w.Flush()

	return buf.String(), nil
}
