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
	"errors"
	"regexp"
	"strings"

	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/tag"
)

var (
	FormatPersonInfo = "NAME [(aka NICK1[, NICK2])] :: description [#tag1, #tag2] [@location1, @location2] [$id:FRIEND_ID]"
	ErrNoInfo        = errors.New("no information provided")
	regexPerson      *regexp.Regexp
)

func init() {
	regexPerson = regexp.MustCompile(`(?m)^(?P<name>[^\(\$:\n]+?)\s*(?:\(\s*a\.?k\.?a\.?\s+(?P<nicknames>[^)]*)\))?\s*(?:\$id:(?P<id>[^\s:]+))?\s*(?:::\s*(?P<desc>.+))?$`)
}

func parseNicknames(raw string) []string {
	raw = strings.ReplaceAll(raw, `"`, "")
	parts := strings.Split(raw, ",")

	var cleaned []string

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}

	return cleaned
}

func ParsePerson(s string) (friend.Person, error) {
	if s == "" {
		return friend.Person{}, ErrNoInfo
	}

	tags := tag.Tags(Parse(s)).ToNames()
	locations := ParseLocMarkers(s)

	s = RemoveTagMarkers(s)
	s = RemoveLocMarkers(s)

	// TODO: Remove them from the string after parsing.

	matches := regexPerson.FindStringSubmatch(s)

	if matches == nil {
		return friend.Person{}, ErrNoInfo
	}

	name := strings.TrimSpace(matches[1])
	nicknames := parseNicknames(matches[2])
	desc := strings.TrimSpace(matches[4])

	return friend.Person{
		Name:      name,
		Nicknames: nicknames,
		Desc:      desc,
		Tags:      tags,
		Locations: locations,
	}, nil
}

func RenderPerson(p friend.Person) string {
	var sb strings.Builder

	sb.WriteString(p.Name)

	if len(p.Nicknames) > 0 {
		sb.WriteString(" (a.k.a. ")
		sb.WriteString(strings.Join(p.Nicknames, ", "))
		sb.WriteString(")")
	}

	if p.Desc != "" {
		sb.WriteString(" :: ")
		sb.WriteString(p.Desc)
	}

	if len(p.Tags) > 0 {
		sb.WriteString(" ")
		sb.WriteString(RenderTags(p.Tags))
	}

	if len(p.Locations) > 0 {
		sb.WriteString(" ")
		sb.WriteString(RenderLocMarkers(p.Locations))
	}

	return sb.String()
}
