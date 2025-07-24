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

package lang

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/tag"

	"github.com/roma-glushko/frens/internal/utils"
)

var (
	locMarkerRe *regexp.Regexp
	locRe       *regexp.Regexp
)

var (
	FormatLocationMarkers = "@location1[, @location2...]"
	FormatLocationInfo    = fmt.Sprintf(
		"NAME[, COUNTRY] [(aka ALIAS1[, ALIAS2...])] :: DESCRIPTION [%s] [$id:LOCATION_ID]",
		FormatTags,
	)
)

func init() {
	locMarkerRe = regexp.MustCompile(`@([\p{L}0-9_-]+)`)
	locRe = regexp.MustCompile(
		`(?m)^(?P<name>[\p{L}\s]+)(?:,\s*(?P<country>[\p{L}\s]+))?(?:\s*\((?:aka|a\.k\.a\.)\s+(?P<aliases>[^)]+)\))?\s*::\s*(?P<description>.*?)?\s*$`,
	)
}

type locProps struct {
	ID string `frentxt:"id"`
}

// ExtractLocation extracts location information from a string.
func ExtractLocation(s string) (friend.Location, error) {
	if s == "" {
		return friend.Location{}, ErrNoInfo
	}

	props, err := ExtractProps[locProps](s)
	if err != nil {
		return friend.Location{}, fmt.Errorf("failed to parse location properties: %w", err)
	}

	tags := tag.Tags(ExtractTags(s)).ToNames()

	s = RemoveTags(s)
	s = RemoveProps(s)

	matches := locRe.FindStringSubmatch(s)

	if matches == nil {
		return friend.Location{}, ErrNoInfo
	}

	name := strings.TrimSpace(matches[1])
	country := strings.TrimSpace(matches[2])
	aliases := extractNicknames(matches[3])
	desc := strings.TrimSpace(matches[4])

	return friend.Location{
		ID:      props.ID,
		Name:    name,
		Country: country,
		Aliases: aliases,
		Desc:    desc,
		Tags:    tags,
	}, nil
}

func ExtractLocationQuery(q string) (friend.ListLocationQuery, error) {
	tags := tag.Tags(ExtractTags(q)).ToNames()

	q = RemoveTags(q)

	search := strings.TrimSpace(q)

	return friend.ListLocationQuery{
		Search: search,
		Tags:   tags,
		// TODO: parse sorting options
	}, nil
}

func RenderLocation(l friend.Location) string {
	var sb strings.Builder

	sb.WriteString(l.Name)

	if l.Country != "" {
		sb.WriteString(", ")
		sb.WriteString(l.Country)
	}

	if len(l.Aliases) > 0 {
		sb.WriteString(" (a.k.a. ")
		sb.WriteString(strings.Join(l.Aliases, ", "))
		sb.WriteString(")")
	}

	if l.Desc != "" {
		sb.WriteString(" :: ")
		sb.WriteString(l.Desc)
	}

	if len(l.Tags) > 0 {
		sb.WriteString(" ")
		sb.WriteString(RenderTags(l.Tags))
	}

	if l.ID == "" {
		sb.WriteString(" ")
		sb.WriteString(RenderProps(locProps{ID: l.ID}))
	}

	return sb.String()
}

func ExtractLocMarkers(s string) []string {
	matches := locMarkerRe.FindAllString(s, -1)
	locationIDs := make([]string, len(matches))

	for i, match := range matches {
		locationIDs[i] = strings.TrimLeft(match, "@")
	}

	return utils.Unique(locationIDs)
}

func RemoveLocMarkers(s string) string {
	return locMarkerRe.ReplaceAllString(s, "")
}

func RenderLocMarkers(locations []string) string {
	if len(locations) == 0 {
		return ""
	}

	markers := make([]string, len(locations))

	for i, loc := range locations {
		loc = strings.TrimSpace(loc)

		if loc == "" {
			continue
		}

		markers[i] = "@" + loc
	}

	return strings.Join(markers, " ")
}
