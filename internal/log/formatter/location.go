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

func (l LocationTextFormatter) FormatSingle(e any) (string, error) {
	location, ok := e.(friend.Location)

	if !ok {
		return "", fmt.Errorf("expected 'friend.Location'")
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(" %s [%s]", labelStyle.Render(location.String()), location.ID))
	sb.WriteString("\n")

	if len(location.Tags) > 0 {
		sb.WriteString(" • " + tagStyle.Render(lang.RenderTags(location.Tags)))
		sb.WriteString(" ")
	}

	if location.Desc != "" {
		sb.WriteString("\n")

		wrapped := wrapText(location.Desc, 80)

		for _, line := range wrapped {
			sb.WriteString(" " + line + "\n")
		}

		sb.WriteString("\n")
	}

	return sb.String(), nil
}

func (l LocationTextFormatter) FormatList(el any) (string, error) {
	locations, ok := el.([]*friend.Location)

	if !ok {
		return "", fmt.Errorf("expected '[]friend.Location'")
	}

	var buf bytes.Buffer

	w := tabwriter.NewWriter(&buf, 0, 0, 3, ' ', 0)

	for _, l := range locations {
		_, _ = fmt.Fprintf(
			w,
			" %s\t%s\t%s\n",
			idStyle.Render(l.ID),
			labelStyle.Render(l.String()),
			tagStyle.Render(lang.RenderTags(l.Tags)),
		)
	}

	_ = w.Flush()

	return buf.String(), nil
}
