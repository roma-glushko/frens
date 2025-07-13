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
	"strings"
	"time"

	"github.com/markusmobius/go-dateparser"
	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/tag"
)

var (
	DatePartition      = ":"
	FormatActivityInfo = "[DATE or RELATIVE DATE]: DESCRIPTION [#tag1, #tag2...] [@location1, @location2...]"
)

func ExtractDate(s string) time.Time {
	ts := time.Now().UTC()

	if s != "" {
		parsedDate, err := dateparser.Parse(nil, s)

		if err != nil {
			ts = time.Now().UTC()
		} else {
			ts = parsedDate.Time.UTC()
		}
	}

	return ts
}

func ExtractActivity(s string) (friend.Event, error) {
	if s == "" {
		return friend.Event{}, ErrNoInfo
	}

	parts := strings.SplitN(s, DatePartition, 2)

	dateStr := ""
	desc := parts[0]

	if len(parts) > 1 {
		dateStr = parts[0]
		desc = parts[1]
	}

	dateStr = strings.TrimSpace(dateStr)
	desc = strings.TrimSpace(desc)

	ts := ExtractDate(dateStr)

	if desc == "" {
		return friend.Event{}, ErrNoInfo
	}

	tags := tag.Tags(ExtractTags(desc)).ToNames()
	locations := ExtractLocMarkers(desc)

	desc = RemoveTags(desc)
	desc = RemoveLocMarkers(desc)

	return friend.Event{
		Type:      friend.EventTypeActivity,
		Date:      ts,
		Desc:      desc,
		Tags:      tags,
		Locations: locations,
	}, nil
}

func RenderActivity(a friend.Event) string {
	var sb strings.Builder

	if !a.Date.IsZero() {
		sb.WriteString(a.Date.Format("2006-01-02 15:04:05"))
		sb.WriteString(DatePartition)
		sb.WriteString(" ")
	}

	sb.WriteString(a.Desc)

	if len(a.Locations) > 0 {
		sb.WriteString(" ")
		sb.WriteString(RenderLocMarkers(a.Locations))
	}

	if len(a.Tags) > 0 {
		sb.WriteString(" ")
		sb.WriteString(RenderTags(a.Tags))
	}

	return sb.String()
}
