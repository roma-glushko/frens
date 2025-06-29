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

package journal

import (
	"fmt"
	"github.com/roma-glushko/frens/internal/utils"
	"strings"
	"sync"

	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/tag"
)

type ListFriendQuery struct {
	Location string
	Tag      string
}

type Data struct {
	DirPath    string
	Tags       tag.Tags
	Friends    []friend.Person
	Locations  friend.Locations
	Activities []friend.Event

	dirty           bool
	matcherMu       sync.Mutex
	friendMatcher   *utils.Matcher[friend.Person]
	locationMatcher *utils.Matcher[friend.Location]
}

func (d *Data) Init() {
	// TODO: implement
}

func (d *Data) Dirty() bool {
	return d.dirty
}

func (d *Data) Path() string {
	return d.DirPath
}

func (d *Data) AddFriend(f friend.Person) {
	d.Friends = append(d.Friends, f)

	d.dirty = true
}

func (d *Data) GetFriend(q string) (*friend.Person, error) {
	matches := d.frenMatcher().Match(q)

	if len(matches) == 0 {
		return nil, fmt.Errorf("no friends found for '%s'", q)
	}

	if len(matches) > 1 {
		names := make([]string, 0, len(matches))

		for _, m := range matches {
			for _, f := range m.Entities {
				names = append(names, f.Name)
			}
		}

		return nil, fmt.Errorf("multiple friends found for '%s': %s", q, strings.Join(names, ", "))
	}

	m := matches[0]

	if len(m.Entities) == 0 {
		return nil, fmt.Errorf("no friends found for '%s'", q)
	}

	if len(m.Entities) > 1 {
		names := make([]string, 0, len(m.Entities))

		for _, f := range m.Entities {
			names = append(names, f.Name)
		}

		return nil, fmt.Errorf("multiple friends found for '%s': %s", q, strings.Join(names, ", "))
	}

	f := m.Entities[0]

	return &f, nil
}

func (d *Data) AddLocation(l friend.Location) {
	d.Locations = append(d.Locations, l)

	d.dirty = true
}

func (d *Data) GetLocation(q string) (*friend.Location, error) {
	matches := d.locMatcher().Match(q)

	if len(matches) == 0 {
		return nil, fmt.Errorf("no locations found for '%s'", q)
	}

	if len(matches) > 1 {
		names := make([]string, 0, len(matches))

		for _, m := range matches {
			for _, f := range m.Entities {
				names = append(names, f.Name)
			}
		}

		return nil, fmt.Errorf("multiple locations found for '%s': %s", q, strings.Join(names, ", "))
	}

	m := matches[0]

	if len(m.Entities) == 0 {
		return nil, fmt.Errorf("no locations found for '%s'", q)
	}

	if len(m.Entities) > 1 {
		names := make([]string, 0, len(m.Entities))

		for _, f := range m.Entities {
			names = append(names, f.Name)
		}

		return nil, fmt.Errorf("multiple locations found for '%s': %s", q, strings.Join(names, ", "))
	}

	l := m.Entities[0]

	return &l, nil
}

func (d *Data) AddTags(t []tag.Tag) {
	d.Tags = append(d.Tags, t...).Unique()

	d.dirty = true
}

func (d *Data) AddActivity(e friend.Event) {
	_ = d.frenMatcher().Match(e.Desc)
	_ = d.locMatcher().Match(e.Desc)

	// TODO: record locs/friends

	tags := tag.Match(e.Desc)

	if len(tags) > 0 {
		d.AddTags(tags)
		tag.Add(&e, tags)
		// TODO: update tag stats
	}

	d.Activities = append(d.Activities, e)

	d.dirty = true
}

func (d *Data) ListFriends(q ListFriendQuery) []friend.Person {
	v := make([]friend.Person, 0, 5)

	for _, f := range d.Friends {
		if q.Location != "" && !f.HasLocation(q.Location) {
			continue
		}

		if q.Tag != "" && !tag.HasTag(&f, q.Tag) {
			continue
		}

		v = append(v, f)
	}

	return v
}

func (d *Data) locMatcher() *utils.Matcher[friend.Location] {
	if d.locationMatcher == nil {
		d.matcherMu.Lock()
		defer d.matcherMu.Unlock()

		d.locationMatcher = utils.NewMatcher[friend.Location]()

		for _, l := range d.Locations {
			d.locationMatcher.Add(l)
		}
	}

	return d.locationMatcher
}

func (d *Data) frenMatcher() *utils.Matcher[friend.Person] {
	if d.friendMatcher == nil {
		d.matcherMu.Lock()
		defer d.matcherMu.Unlock()

		d.friendMatcher = utils.NewMatcher[friend.Person]()

		for _, f := range d.Friends {
			d.friendMatcher.Add(f)
		}
	}

	return d.friendMatcher
}
