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

package friend

import (
	"sort"
	"strings"
	"time"

	"github.com/roma-glushko/frens/internal/tag"

	"github.com/markusmobius/go-dateparser"
)

const (
	DatePartition = ":"
)

type EventType string

const (
	EventTypeActivity EventType = "activity"
	EventTypeNote     EventType = "note"
)

type Event struct {
	Type EventType
	Date time.Time
	Desc string

	Friends   []string
	Locations []string
	Tags      []string
}

var _ tag.Tagged = (*Event)(nil)

// NewEvent creates a new event (activity, note, etc.) with date and description.
//
//	The original event description comes in the following formats:
//		"<date>: <description>"
//		"<description>" (The current date will be used by default.)
func NewEvent(etype EventType, rawDesc string) Event {
	parts := strings.SplitN(rawDesc, DatePartition, 2)

	dateStr := ""
	desc := parts[0]

	if len(parts) > 1 {
		dateStr = parts[0]
		desc = parts[1]
	}

	desc = strings.TrimSpace(desc)

	ts := time.Now().UTC()

	if dateStr != "" {
		parsedDate, err := dateparser.Parse(nil, dateStr)

		if err != nil {
			ts = time.Now().UTC()
		} else {
			ts = parsedDate.Time.UTC()
		}
	}

	return Event{
		Type: etype,
		Date: ts,
		Desc: desc,
	}
}

func (f *Event) SetTags(tags []string) {
	f.Tags = tags
}

func (f *Event) GetTags() []string {
	return f.Tags
}

type Events []Event

var _ sort.Interface = (*Events)(nil)

func (e Events) Len() int           { return len(e) }
func (e Events) Less(i, j int) bool { return e[i].Date.Before(e[j].Date) }
func (e Events) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
