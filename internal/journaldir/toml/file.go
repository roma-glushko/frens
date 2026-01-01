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

package toml

import (
	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/tag"
)

const (
	FileNameFriends    = "friends.toml"
	FileNameActivities = "activities.toml"
)

// Data represents the persisted journal data.
// This is the pure data structure that storage operates on,
// separate from the Journal service which adds behavior (matching, dirty tracking, etc).
type Data struct {
	DirPath    string
	Tags       tag.Tags
	Friends    []*friend.Person
	Locations  []*friend.Location
	Activities []*friend.Event
	Notes      []*friend.Event
}

// FriendsFile represents the structure of friends.toml
type FriendsFile struct {
	Tags      []tag.Tag          `toml:"tags"`
	Friends   []*friend.Person   `toml:"friends"`
	Locations []*friend.Location `toml:"locations"`
}

// EventsFile represents the structure of activities.toml
type EventsFile struct {
	Activities []*friend.Event `toml:"activities"`
	Notes      []*friend.Event `toml:"notes"`
}
