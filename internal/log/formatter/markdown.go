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
	"fmt"
	"strings"

	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/log"
)

func init() {
	log.RegisterFormatter(log.FormatMarkdown, friend.Person{}, PersonMarkdownFormatter{})
	log.RegisterFormatter(log.FormatMarkdown, friend.Contact{}, ContactMarkdownFormatter{})
	log.RegisterFormatter(log.FormatMarkdown, friend.Event{}, EventMarkdownFormatter{})
	log.RegisterFormatter(log.FormatMarkdown, friend.Location{}, LocationMarkdownFormatter{})
	log.RegisterFormatter(log.FormatMarkdown, friend.Date{}, DateMarkdownFormatter{})
	log.RegisterFormatter(
		log.FormatMarkdown,
		friend.WishlistItem{},
		WishlistItemMarkdownFormatter{},
	)
}

// Helper to render tags as markdown
func renderTagsMd(tags []string) string {
	if len(tags) == 0 {
		return ""
	}

	result := make([]string, 0, len(tags))

	for _, tag := range tags {
		result = append(result, "`"+tag+"`")
	}

	return strings.Join(result, " ")
}

// formatWishlistItemDesc formats a wishlist item description for markdown
func formatWishlistItemDesc(desc, link string) string {
	if link == "" {
		return desc
	}

	if desc != "" {
		return fmt.Sprintf("[%s](%s)", desc, link)
	}

	return fmt.Sprintf("[link](%s)", link)
}

// ============================================================================
// Person Markdown Formatter
// ============================================================================

type PersonMarkdownFormatter struct{}

var _ log.Formatter = (*PersonMarkdownFormatter)(nil)

func (p PersonMarkdownFormatter) FormatSingle(_ log.FormatterContext, e any) (string, error) {
	person, ok := e.(friend.Person)
	if !ok {
		return "", ErrInvalidEntity
	}

	var sb strings.Builder

	p.writeHeader(&sb, person)
	p.writeMetadata(&sb, person)
	p.writeContacts(&sb, person.Contacts)
	p.writeDates(&sb, person.Dates)
	p.writeWishlist(&sb, person.Wishlist)

	return sb.String(), nil
}

func (p PersonMarkdownFormatter) writeHeader(sb *strings.Builder, person friend.Person) {
	fmt.Fprintf(sb, "## %s\n\n", person.String())
	fmt.Fprintf(sb, "- **ID:** `%s`\n", person.ID)

	if len(person.Nicknames) > 0 {
		fmt.Fprintf(sb, "- **Nicknames:** %s\n", strings.Join(person.Nicknames, ", "))
	}

	if len(person.Locations) > 0 {
		fmt.Fprintf(sb, "- **Locations:** %s\n", strings.Join(person.Locations, ", "))
	}

	if len(person.Tags) > 0 {
		fmt.Fprintf(sb, "- **Tags:** %s\n", renderTagsMd(person.Tags))
	}

	fmt.Fprintf(sb, "- **Notes:** %d\n", person.Notes)
	fmt.Fprintf(sb, "- **Activities:** %d\n", person.Activities)
}

func (p PersonMarkdownFormatter) writeMetadata(sb *strings.Builder, person friend.Person) {
	if person.Desc != "" {
		fmt.Fprintf(sb, "\n%s\n", person.Desc)
	}
}

func (p PersonMarkdownFormatter) writeContacts(sb *strings.Builder, contacts []*friend.Contact) {
	if len(contacts) == 0 {
		return
	}

	sb.WriteString("\n### Contacts\n\n")

	for _, c := range contacts {
		fmt.Fprintf(sb, "- **%s:** %s", c.Type, c.Value)

		if len(c.Tags) > 0 {
			sb.WriteString(" " + renderTagsMd(c.Tags))
		}

		sb.WriteString("\n")
	}
}

func (p PersonMarkdownFormatter) writeDates(sb *strings.Builder, dates []*friend.Date) {
	if len(dates) == 0 {
		return
	}

	sb.WriteString("\n### Dates\n\n")

	for _, d := range dates {
		fmt.Fprintf(sb, "- **%s**", d.DateExpr)

		if d.Desc != "" {
			sb.WriteString(" — " + d.Desc)
		}

		if len(d.Tags) > 0 {
			sb.WriteString(" " + renderTagsMd(d.Tags))
		}

		sb.WriteString("\n")
	}
}

func (p PersonMarkdownFormatter) writeWishlist(
	sb *strings.Builder,
	wishlist []*friend.WishlistItem,
) {
	if len(wishlist) == 0 {
		return
	}

	sb.WriteString("\n### Wishlist\n\n")

	for _, w := range wishlist {
		sb.WriteString("- " + formatWishlistItemDesc(w.Desc, w.Link))

		if w.Price != "" {
			sb.WriteString(" ($" + w.Price + ")")
		}

		if len(w.Tags) > 0 {
			sb.WriteString(" " + renderTagsMd(w.Tags))
		}

		sb.WriteString("\n")
	}
}

func (p PersonMarkdownFormatter) FormatList(_ log.FormatterContext, el any) (string, error) {
	persons, ok := el.([]friend.Person)
	if !ok {
		return "", ErrInvalidEntity
	}

	var sb strings.Builder

	sb.WriteString("| ID | Name | Tags | Locations | Notes | Activities |\n")
	sb.WriteString("|---|---|---|---|---|---|\n")

	for _, person := range persons {
		sb.WriteString(fmt.Sprintf("| `%s` | %s | %s | %s | %d | %d |\n",
			person.ID,
			person.String(),
			renderTagsMd(person.Tags),
			strings.Join(person.Locations, ", "),
			person.Notes,
			person.Activities,
		))
	}

	return sb.String(), nil
}

// ============================================================================
// Contact Markdown Formatter
// ============================================================================

type ContactMarkdownFormatter struct{}

var _ log.Formatter = (*ContactMarkdownFormatter)(nil)

func (f ContactMarkdownFormatter) FormatSingle(_ log.FormatterContext, e any) (string, error) {
	var c friend.Contact

	switch v := e.(type) {
	case friend.Contact:
		c = v
	case *friend.Contact:
		c = *v
	default:
		return "", ErrInvalidEntity
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("## %s: %s\n\n", c.Type, c.Value))
	sb.WriteString(fmt.Sprintf("- **ID:** `%s`\n", c.ID))

	if c.Person != "" {
		sb.WriteString(fmt.Sprintf("- **Person:** %s\n", c.Person))
	}

	if len(c.Tags) > 0 {
		sb.WriteString(fmt.Sprintf("- **Tags:** %s\n", renderTagsMd(c.Tags)))
	}

	return sb.String(), nil
}

func (f ContactMarkdownFormatter) FormatList(_ log.FormatterContext, el any) (string, error) {
	contacts, ok := el.([]friend.Contact)
	if !ok {
		return "", ErrInvalidEntity
	}

	var sb strings.Builder

	sb.WriteString("| ID | Person | Type | Value | Tags |\n")
	sb.WriteString("|---|---|---|---|---|\n")

	for _, c := range contacts {
		sb.WriteString(fmt.Sprintf("| `%s` | %s | %s | %s | %s |\n",
			c.ID,
			c.Person,
			c.Type,
			c.Value,
			renderTagsMd(c.Tags),
		))
	}

	return sb.String(), nil
}

// ============================================================================
// Event Markdown Formatter
// ============================================================================

type EventMarkdownFormatter struct{}

var _ log.Formatter = (*EventMarkdownFormatter)(nil)

func (f EventMarkdownFormatter) FormatSingle(_ log.FormatterContext, e any) (string, error) {
	event, ok := e.(friend.Event)
	if !ok {
		return "", ErrInvalidEntity
	}

	var sb strings.Builder

	date := event.Date.Format("Mon Jan 2, 2006 15:04 MST")

	sb.WriteString(fmt.Sprintf("## %s\n\n", date))
	sb.WriteString(fmt.Sprintf("- **ID:** `%s`\n", event.ID))
	sb.WriteString(fmt.Sprintf("- **Type:** %s\n", event.Type))

	if len(event.FriendIDs) > 0 {
		sb.WriteString(fmt.Sprintf("- **Friends:** %s\n", strings.Join(event.FriendIDs, ", ")))
	}

	if len(event.LocationIDs) > 0 {
		sb.WriteString(fmt.Sprintf("- **Locations:** %s\n", strings.Join(event.LocationIDs, ", ")))
	}

	if len(event.Tags) > 0 {
		sb.WriteString(fmt.Sprintf("- **Tags:** %s\n", renderTagsMd(event.Tags)))
	}

	if event.Desc != "" {
		sb.WriteString(fmt.Sprintf("\n%s\n", event.Desc))
	}

	return sb.String(), nil
}

func (f EventMarkdownFormatter) FormatList(_ log.FormatterContext, el any) (string, error) {
	events, ok := el.([]friend.Event)
	if !ok {
		return "", ErrInvalidEntity
	}

	var sb strings.Builder

	sb.WriteString("| ID | Date | Description | Friends | Tags | Locations |\n")
	sb.WriteString("|---|---|---|---|---|---|\n")

	for _, e := range events {
		desc := CutStr(e.Desc, 50)
		sb.WriteString(fmt.Sprintf("| `%s` | %s | %s | %s | %s | %s |\n",
			e.ID,
			e.Date.Format("2006-01-02 15:04"),
			desc,
			strings.Join(e.FriendIDs, ", "),
			renderTagsMd(e.Tags),
			strings.Join(e.LocationIDs, ", "),
		))
	}

	return sb.String(), nil
}

// ============================================================================
// Location Markdown Formatter
// ============================================================================

type LocationMarkdownFormatter struct{}

var _ log.Formatter = (*LocationMarkdownFormatter)(nil)

func (l LocationMarkdownFormatter) FormatSingle(_ log.FormatterContext, e any) (string, error) {
	location, ok := e.(friend.Location)
	if !ok {
		return "", ErrInvalidEntity
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("## %s\n\n", location.String()))
	sb.WriteString(fmt.Sprintf("- **ID:** `%s`\n", location.ID))

	if location.Country != "" {
		sb.WriteString(fmt.Sprintf("- **Country:** %s\n", location.Country))
	}

	if len(location.Aliases) > 0 {
		sb.WriteString(fmt.Sprintf("- **Aliases:** %s\n", strings.Join(location.Aliases, ", ")))
	}

	if len(location.Tags) > 0 {
		sb.WriteString(fmt.Sprintf("- **Tags:** %s\n", renderTagsMd(location.Tags)))
	}

	if location.Lat != nil && location.Lng != nil {
		sb.WriteString(fmt.Sprintf("- **Coordinates:** %.4f, %.4f\n", *location.Lat, *location.Lng))
	}

	if location.Desc != "" {
		sb.WriteString(fmt.Sprintf("\n%s\n", location.Desc))
	}

	return sb.String(), nil
}

func (l LocationMarkdownFormatter) FormatList(_ log.FormatterContext, el any) (string, error) {
	locations, ok := el.([]friend.Location)
	if !ok {
		return "", ErrInvalidEntity
	}

	var sb strings.Builder

	sb.WriteString("| ID | Name | Country | Tags |\n")
	sb.WriteString("|---|---|---|---|\n")

	for _, loc := range locations {
		sb.WriteString(fmt.Sprintf("| `%s` | %s | %s | %s |\n",
			loc.ID,
			loc.String(),
			loc.Country,
			renderTagsMd(loc.Tags),
		))
	}

	return sb.String(), nil
}

// ============================================================================
// Date Markdown Formatter
// ============================================================================

type DateMarkdownFormatter struct{}

var _ log.Formatter = (*DateMarkdownFormatter)(nil)

func (f DateMarkdownFormatter) FormatSingle(_ log.FormatterContext, e any) (string, error) {
	var dt friend.Date

	switch v := e.(type) {
	case friend.Date:
		dt = v
	case *friend.Date:
		dt = *v
	default:
		return "", ErrInvalidEntity
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("## %s\n\n", dt.DateExpr))
	sb.WriteString(fmt.Sprintf("- **ID:** `%s`\n", dt.ID))

	if dt.Person != "" {
		sb.WriteString(fmt.Sprintf("- **Person:** %s\n", dt.Person))
	}

	if dt.Calendar != "" {
		sb.WriteString(fmt.Sprintf("- **Calendar:** %s\n", dt.Calendar))
	}

	if len(dt.Tags) > 0 {
		sb.WriteString(fmt.Sprintf("- **Tags:** %s\n", renderTagsMd(dt.Tags)))
	}

	if dt.Desc != "" {
		sb.WriteString(fmt.Sprintf("\n%s\n", dt.Desc))
	}

	return sb.String(), nil
}

func (f DateMarkdownFormatter) FormatList(_ log.FormatterContext, el any) (string, error) {
	dates, ok := el.([]*friend.Date)
	if !ok {
		return "", ErrInvalidEntity
	}

	var sb strings.Builder

	sb.WriteString("| ID | Person | Date | Tags |\n")
	sb.WriteString("|---|---|---|---|\n")

	for _, dt := range dates {
		sb.WriteString(fmt.Sprintf("| `%s` | %s | %s | %s |\n",
			dt.ID,
			dt.Person,
			dt.DateExpr,
			renderTagsMd(dt.Tags),
		))
	}

	return sb.String(), nil
}

// ============================================================================
// WishlistItem Markdown Formatter
// ============================================================================

type WishlistItemMarkdownFormatter struct{}

var _ log.Formatter = (*WishlistItemMarkdownFormatter)(nil)

func (f WishlistItemMarkdownFormatter) FormatSingle(_ log.FormatterContext, e any) (string, error) {
	var w friend.WishlistItem

	switch v := e.(type) {
	case friend.WishlistItem:
		w = v
	case *friend.WishlistItem:
		w = *v
	default:
		return "", ErrInvalidEntity
	}

	var sb strings.Builder

	title := w.Desc
	if title == "" && w.Link != "" {
		title = w.Link
	}

	sb.WriteString(fmt.Sprintf("## %s\n\n", title))
	sb.WriteString(fmt.Sprintf("- **ID:** `%s`\n", w.ID))

	if w.Person != "" {
		sb.WriteString(fmt.Sprintf("- **Person:** %s\n", w.Person))
	}

	if w.Link != "" {
		sb.WriteString(fmt.Sprintf("- **Link:** [%s](%s)\n", w.Link, w.Link))
	}

	if w.Price != "" {
		sb.WriteString(fmt.Sprintf("- **Price:** %s\n", w.Price))
	}

	if len(w.Tags) > 0 {
		sb.WriteString(fmt.Sprintf("- **Tags:** %s\n", renderTagsMd(w.Tags)))
	}

	return sb.String(), nil
}

func (f WishlistItemMarkdownFormatter) FormatList(_ log.FormatterContext, el any) (string, error) {
	items, ok := el.([]friend.WishlistItem)
	if !ok {
		return "", ErrInvalidEntity
	}

	var sb strings.Builder

	sb.WriteString("| ID | Person | Description | Price | Tags |\n")
	sb.WriteString("|---|---|---|---|---|\n")

	for _, item := range items {
		desc := item.Desc
		if desc == "" && item.Link != "" {
			desc = CutStr(item.Link, 40)
		}

		sb.WriteString(fmt.Sprintf("| `%s` | %s | %s | %s | %s |\n",
			item.ID,
			item.Person,
			desc,
			item.Price,
			renderTagsMd(item.Tags),
		))
	}

	return sb.String(), nil
}
