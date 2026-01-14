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

func (p PersonTextFormatter) FormatSingle(ctx log.FormatterContext, e any) (string, error) {
	person, ok := e.(friend.Person)

	if !ok {
		return "", ErrInvalidEntity
	}

	if ctx.Density == log.DensityCompact {
		return p.formatCompact(person), nil
	}

	return p.formatRegular(person), nil
}

func (p PersonTextFormatter) formatCompact(person friend.Person) string {
	parts := []string{idStyle.Render(person.ID), person.String()}

	if len(person.Tags) > 0 {
		parts = append(parts, tagStyle.Render(lang.RenderTags(person.Tags)))
	}

	if len(person.Locations) > 0 {
		parts = append(parts, locationStyle.Render(lang.RenderLocMarkers(person.Locations)))
	}

	return strings.Join(parts, " ") + "\n"
}

func (p PersonTextFormatter) formatRegular(person friend.Person) string {
	var sb strings.Builder

	p.writeHeader(&sb, person)
	p.writeContacts(&sb, person.Contacts)
	p.writeDates(&sb, person.Dates)
	p.writeWishlist(&sb, person.Wishlist)

	sb.WriteString("\n")

	return sb.String()
}

func (p PersonTextFormatter) writeHeader(sb *strings.Builder, person friend.Person) {
	fmt.Fprintf(sb, "%s (%s)", labelStyle.Render(person.String()), idStyle.Render(person.ID))
	sb.WriteString("\n")

	if len(person.Tags) > 0 || len(person.Locations) > 0 {
		sb.WriteString("  ")

		if len(person.Tags) > 0 {
			sb.WriteString(tagStyle.Render(lang.RenderTags(person.Tags)))
		}

		if len(person.Locations) > 0 {
			if len(person.Tags) > 0 {
				sb.WriteString(" ")
			}

			sb.WriteString(locationStyle.Render(lang.RenderLocMarkers(person.Locations)))
		}

		sb.WriteString("\n")
	}

	if person.Desc != "" {
		wrapped := wrapText(person.Desc, 80)

		for _, line := range wrapped {
			sb.WriteString("  " + line + "\n")
		}
	}
}

func (p PersonTextFormatter) writeContacts(sb *strings.Builder, contacts []*friend.Contact) {
	if len(contacts) == 0 {
		return
	}

	sb.WriteString("\n")
	sb.WriteString("  " + labelStyle.Render("Contacts") + "\n")

	for _, c := range contacts {
		fmt.Fprintf(sb, "    %s %s: %s", log.BulletChar, c.Type, c.Value)

		if len(c.Tags) > 0 {
			sb.WriteString(" " + tagStyle.Render(lang.RenderTags(c.Tags)))
		}

		sb.WriteString("\n")
	}
}

func (p PersonTextFormatter) writeDates(sb *strings.Builder, dates []*friend.Date) {
	if len(dates) == 0 {
		return
	}

	sb.WriteString("\n")
	sb.WriteString("  " + labelStyle.Render("Dates") + "\n")

	for _, d := range dates {
		sb.WriteString("    " + log.BulletChar + " " + d.DateExpr)

		if d.Desc != "" {
			sb.WriteString(" - " + d.Desc)
		}

		if len(d.Tags) > 0 {
			sb.WriteString(" " + tagStyle.Render(lang.RenderTags(d.Tags)))
		}

		sb.WriteString("\n")
	}
}

func (p PersonTextFormatter) writeWishlist(sb *strings.Builder, wishlist []*friend.WishlistItem) {
	if len(wishlist) == 0 {
		return
	}

	sb.WriteString("\n")
	sb.WriteString("  " + labelStyle.Render("Wishlist") + "\n")

	for _, w := range wishlist {
		sb.WriteString("    " + log.BulletChar + " " + formatWishlistDesc(w.Desc, w.Link))

		if w.Price != "" {
			sb.WriteString(" ($" + w.Price + ")")
		}

		if len(w.Tags) > 0 {
			sb.WriteString(" " + tagStyle.Render(lang.RenderTags(w.Tags)))
		}

		sb.WriteString("\n")
	}
}

func (p PersonTextFormatter) FormatList(ctx log.FormatterContext, el any) (string, error) {
	persons, ok := el.([]friend.Person)

	if !ok {
		return "", ErrInvalidEntity
	}

	var buf bytes.Buffer

	w := tabwriter.NewWriter(&buf, 0, 0, 3, ' ', 0)

	for _, person := range persons {
		if ctx.Density == log.DensityCompact {
			_, _ = fmt.Fprintf(
				w,
				"%s\t%s\t%s\t%s\n",
				idStyle.Render(person.ID),
				person.String(),
				tagStyle.Render(lang.RenderTags(person.Tags)),
				locationStyle.Render(lang.RenderLocMarkers(person.Locations)),
			)
		} else {
			_, _ = fmt.Fprintf(
				w,
				"%s\t%s\t%s\t%s\t%s\n",
				idStyle.Render(person.ID),
				labelStyle.Render(person.String()),
				tagStyle.Render(lang.RenderTags(person.Tags)),
				locationStyle.Render(lang.RenderLocMarkers(person.Locations)),
				countLabel.Render(formatCounts(person.Notes, person.Activities)),
			)
		}
	}

	_ = w.Flush()

	return buf.String(), nil
}

// formatCounts returns a human-readable count string like "2 notes, 1 activity"
func formatCounts(notes, activities int) string {
	var parts []string

	if notes > 0 {
		if notes == 1 {
			parts = append(parts, "1 note")
		} else {
			parts = append(parts, fmt.Sprintf("%d notes", notes))
		}
	}

	if activities > 0 {
		if activities == 1 {
			parts = append(parts, "1 activity")
		} else {
			parts = append(parts, fmt.Sprintf("%d activities", activities))
		}
	}

	return strings.Join(parts, ", ")
}

// formatWishlistDesc formats a wishlist item description for text output
func formatWishlistDesc(desc, link string) string {
	if desc != "" {
		return desc
	}

	if link == "" {
		return ""
	}

	if len(link) > 50 {
		return link[:50] + "..."
	}

	return link
}
