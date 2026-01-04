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

package lang

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/tag"
)

var (
	FormatPersonInfo = fmt.Sprintf(
		"NAME [(aka NICK1[, NICK2...])] %s DESCRIPTION [%s] [%s] [$id:FRIEND_ID]",
		Separator,
		FormatTags,
		FormatLocationMarkers,
	)
	FormatPersonQuery = fmt.Sprintf(
		"[SEARCH TERM] [%s] [%s] [$sort:SORT_OPTION] [$order:ORDER_OPTION]",
		FormatTags,
		FormatLocationMarkers,
	)
	ErrNoInfo = errors.New("no information provided")
	personRe  *regexp.Regexp
)

func init() {
	personRe = regexp.MustCompile(
		`(?m)^(?P<name>[^\(\$:\n]+?)\s*(?:\(\s*a\.?k\.?a\.?\s+(?P<nicknames>[^)]*)\))?\s*(?:\$id:(?P<id>[^\s:]+))?\s*(?:::\s*(?P<desc>.+))?$`,
	)
}

type personProps struct {
	ID string `frentxt:"id"`
}

type orderProps struct {
	SortBy    string `frentxt:"sort"`
	SortOrder string `frentxt:"order"`
}

func extractNicknames(raw string) []string {
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

func ExtractPerson(s string) (friend.Person, error) {
	if s == "" {
		return friend.Person{}, ErrNoInfo
	}

	props, err := ExtractProps[personProps](s)
	if err != nil {
		return friend.Person{}, fmt.Errorf("failed to parse person properties: %w", err)
	}

	tags := tag.Tags(ExtractTags(s)).ToNames()
	locations := ExtractLocMarkers(s)

	s = RemoveTags(s)
	s = RemoveLocMarkers(s)
	s = RemoveProps(s)

	matches := personRe.FindStringSubmatch(s)

	if matches == nil {
		return friend.Person{}, ErrNoInfo
	}

	name := strings.TrimSpace(matches[1])
	nicknames := extractNicknames(matches[2])
	desc := strings.TrimSpace(matches[4])

	return friend.Person{
		ID:        props.ID,
		Name:      name,
		Nicknames: nicknames,
		Desc:      desc,
		Tags:      tags,
		Locations: locations,
	}, nil
}

func ExtractPersonQuery(q string) (friend.ListFriendQuery, error) {
	tags := tag.Tags(ExtractTags(q)).ToNames()
	locations := ExtractLocMarkers(q)

	props, err := ExtractProps[orderProps](q)
	if err != nil {
		return friend.ListFriendQuery{}, fmt.Errorf(
			"failed to parse friend list query properties: %w",
			err,
		)
	}

	q = RemoveTags(q)
	q = RemoveLocMarkers(q)
	q = RemoveProps(q)

	search := strings.TrimSpace(q)

	return friend.ListFriendQuery{
		Keyword:   search,
		Locations: locations,
		Tags:      tags,
		SortBy:    friend.SortOption(props.SortBy),
		SortOrder: friend.SortOrderOption(props.SortOrder),
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
		sb.WriteString(" ")
		sb.WriteString(Separator)
		sb.WriteString(" ")
		sb.WriteString(p.Desc)
	}

	if len(p.Locations) > 0 {
		sb.WriteString(" ")
		sb.WriteString(RenderLocMarkers(p.Locations))
	}

	if len(p.Tags) > 0 {
		sb.WriteString(" ")
		sb.WriteString(RenderTags(p.Tags))
	}

	if p.ID != "" {
		sb.WriteString(" ")
		sb.WriteString(RenderProps(personProps{ID: p.ID}))
	}

	return sb.String()
}
