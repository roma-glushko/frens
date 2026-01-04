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
	log.RegisterFormatter(log.FormatText, friend.WishlistItem{}, WishlistItemTextFormatter{})
}

type WishlistItemTextFormatter struct{}

var _ log.Formatter = (*WishlistItemTextFormatter)(nil)

func (f WishlistItemTextFormatter) FormatSingle(e any) (string, error) {
	w, ok := e.(*friend.WishlistItem)

	if !ok {
		return "", ErrInvalidEntity
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("[%s]", idStyle.Render(w.ID)))

	if w.Desc != "" {
		sb.WriteString(" " + labelStyle.Render(w.Desc))
	}

	if w.Link != "" {
		sb.WriteString("\n")
		sb.WriteString(" 🔗 " + w.Link)
	}

	if w.Price != "" {
		sb.WriteString("\n")
		sb.WriteString(" 💰 " + w.Price)
	}

	if len(w.Tags) > 0 {
		sb.WriteString("\n")
		sb.WriteString(" • " + tagStyle.Render(lang.RenderTags(w.Tags)))
	}

	return sb.String(), nil
}

func (f WishlistItemTextFormatter) FormatList(el any) (string, error) {
	items, ok := el.([]friend.WishlistItem)

	if !ok {
		return "", ErrInvalidEntity
	}

	var buf bytes.Buffer

	w := tabwriter.NewWriter(&buf, 0, 0, 3, ' ', 0)

	for _, item := range items {
		desc := item.Desc
		if desc == "" && item.Link != "" {
			// Truncate long URLs for list view
			desc = item.Link
			if len(desc) > 40 {
				desc = desc[:40] + "..."
			}
		}

		price := ""
		if item.Price != "" {
			price = "💰 " + item.Price
		}

		_, _ = fmt.Fprintf(
			w,
			" %s\t%s\t%s\t%s\t%s\n",
			idStyle.Render(item.ID),
			labelStyle.Render(item.Person),
			desc,
			price,
			tagStyle.Render(lang.RenderTags(item.Tags)),
		)
	}

	_ = w.Flush()

	return buf.String(), nil
}
