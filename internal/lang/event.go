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
	"strings"
	"time"

	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/tag"
)

var (
	FormatEventInfo = fmt.Sprintf(
		"[DATE or RELATIVE DATE %s] DESCRIPTION [%s] [%s]",
		Separator,
		FormatTags,
		FormatLocationMarkers,
	)
	FormatEventQuery = fmt.Sprintf(
		"[SEARCH TERM] [%s] [%s] [$since:DATE] [$until:DATE] [$sort:SORT_OPTION] [$order:ORDER_OPTION]",
		FormatTags,
		FormatLocationMarkers,
	)
)

type eventProps struct {
	SortBy    string `frentxt:"sort"`
	SortOrder string `frentxt:"order"`
	Since     string `frentxt:"since"`
	Until     string `frentxt:"until"`
}

func ExtractEvent(t friend.EventType, s string) (friend.Event, error) {
	if s == "" {
		return friend.Event{}, ErrNoInfo
	}

	parts := strings.SplitN(s, Separator, 2)

	dateStr := ""
	desc := parts[0]

	if len(parts) > 1 {
		dateStr = parts[0]
		desc = parts[1]
	}

	dateStr = strings.TrimSpace(dateStr)
	desc = strings.TrimSpace(desc)

	ts := ExtractDate(dateStr, time.Now().UTC())

	if desc == "" {
		return friend.Event{}, ErrNoInfo
	}

	tags := tag.Tags(ExtractTags(desc)).ToNames()
	locations := ExtractLocMarkers(desc)

	desc = RemoveTags(desc)
	desc = RemoveLocMarkers(desc)

	return friend.Event{
		Type:      t,
		Date:      ts,
		Desc:      desc,
		Tags:      tags,
		Locations: locations,
	}, nil
}

func RenderEvent(e *friend.Event) string {
	var sb strings.Builder

	if !e.Date.IsZero() {
		sb.WriteString(e.Date.Format("2006-01-02 15:04:05"))
		sb.WriteString(" ")
		sb.WriteString(Separator)
		sb.WriteString(" ")
	}

	sb.WriteString(e.Desc)

	if len(e.Locations) > 0 {
		sb.WriteString(" ")
		sb.WriteString(RenderLocMarkers(e.Locations))
	}

	if len(e.Tags) > 0 {
		sb.WriteString(" ")
		sb.WriteString(RenderTags(e.Tags))
	}

	return sb.String()
}

func ExtractEventQuery(q string) (friend.ListEventQuery, error) {
	props, err := ExtractProps[eventProps](q)
	if err != nil {
		return friend.ListEventQuery{}, fmt.Errorf(
			"failed to parse event list query properties: %w",
			err,
		)
	}

	tags := tag.Tags(ExtractTags(q)).ToNames()
	locations := ExtractLocMarkers(q)

	q = RemoveTags(q)
	q = RemoveLocMarkers(q)
	q = RemoveProps(q)

	search := strings.TrimSpace(q)

	return friend.ListEventQuery{
		Keyword:   search,
		Tags:      tags,
		Locations: locations,
		Since:     ExtractDate(props.Since),
		Until:     ExtractDate(props.Until),
		SortBy:    friend.SortOption(props.SortBy),
		SortOrder: friend.SortOrderOption(props.SortOrder),
	}, nil
}
