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

package friend

import (
	"errors"
	"sort"
	"time"

	"github.com/roma-glushko/frens/internal/tag"
)

var ErrEventDescEmpty = errors.New("event description must be provided")

type EventType string

const (
	EventTypeActivity EventType = "activity"
	EventTypeNote     EventType = "note"
)

type Event struct {
	ID   string    `toml:"id"   json:"id"`
	Type EventType `toml:"type" json:"type"`
	Date time.Time `toml:"date" json:"date"`
	Desc string    `toml:"desc" json:"description"`

	FriendIDs   []string `toml:"friends,omitempty"   json:"friendIds,omitempty"`
	LocationIDs []string `toml:"locations,omitempty" json:"locationIds,omitempty"`
	Tags        []string `toml:"tags,omitempty"      json:"tags,omitempty"`
}

type EventView struct {
	Event
	Friends   Persons   `toml:"-"`
	Locations Locations `toml:"-"`
}

var _ tag.Tagged = (*Event)(nil)

func (e *Event) Validate() error {
	if e.Desc == "" {
		return ErrEventDescEmpty
	}

	return nil
}

func (e *Event) SetTags(tags []string) {
	e.Tags = tags
}

func (e *Event) GetTags() []string {
	return e.Tags
}

type Events []Event

var _ sort.Interface = (*Events)(nil)

func (e Events) Len() int           { return len(e) }
func (e Events) Less(i, j int) bool { return e[i].Date.Before(e[j].Date) }
func (e Events) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
