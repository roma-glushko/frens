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
	"slices"
	"strings"
	"sync"

	"github.com/roma-glushko/frens/internal/lang"

	"github.com/roma-glushko/frens/internal/utils"

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
		return nil, fmt.Errorf("no friend found for '%s'", q)
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

func (d *Data) UpdateFriend(o, n friend.Person) {
	for i, f := range d.Friends {
		if f.Name == o.Name {
			d.Friends[i] = n
			d.dirty = true

			return
		}
	}

	// TODO: update friend references in activities and notes

	// If the friend was not found, add it as a new one
	d.AddFriend(n)
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

func (d *Data) UpdateLocation(o, n friend.Location) {
	for i, l := range d.Locations {
		if l.Name == o.Name {
			d.Locations[i] = n
			d.dirty = true

			return
		}
	}

	// TODO: update friend references in activities and notes

	// If the friend was not found, add it as a new one
	d.AddLocation(n)
}

func (d *Data) AddTags(t []tag.Tag) {
	d.Tags = append(d.Tags, t...).Unique()

	d.dirty = true
}

func (d *Data) AddActivity(e friend.Event) { //nolint:cyclop
	matches := d.frenMatcher().Match(e.Desc)

	certainPersons := make([]friend.Person, 0, len(matches))
	ambiguitiesMatches := make([]utils.Match[friend.Person], 0, len(matches))

	for _, m := range matches {
		if len(m.Entities) == 1 {
			certainPersons = append(certainPersons, m.Entities[0])
			continue
		}

		shortestNameFriend := slices.MinFunc(m.Entities, func(a, b friend.Person) int {
			return strings.Compare(a.Name, b.Name)
		})

		shortestName := shortestNameFriend.Name
		allContains := true

		for _, e := range m.Entities {
			if !strings.Contains(e.Name, shortestName) {
				allContains = false
				break
			}
		}

		if allContains {
			certainPersons = append(certainPersons, shortestNameFriend)
		} else {
			ambiguitiesMatches = append(ambiguitiesMatches, m)
		}
	}

	type friendPair struct {
		KnownPerson       friend.Person
		AmbiguitiesPerson friend.Person
	}

	guessedPersons := make([]friend.Person, 0, len(ambiguitiesMatches))

	if len(ambiguitiesMatches) > 0 {
		rankPairs := make([]friendPair, 0, len(certainPersons)*len(ambiguitiesMatches))

		for _, cp := range certainPersons {
			for _, am := range ambiguitiesMatches {
				for _, ap := range am.Entities {
					rankPairs = append(rankPairs, friendPair{
						KnownPerson:       cp,
						AmbiguitiesPerson: ap,
					})
				}
			}
		}

		for _, act := range d.Activities {
			for _, pair := range rankPairs {
				if slices.Contains(act.Friends, pair.KnownPerson.Name) &&
					slices.Contains(act.Friends, pair.AmbiguitiesPerson.Name) {
					pair.AmbiguitiesPerson.Score++
				}
			}
		}

		for _, am := range ambiguitiesMatches {
			guessedPerson := slices.MaxFunc(am.Entities, func(a, b friend.Person) int {
				if a.Score != b.Score {
					return b.Score - a.Score
				}

				return b.Activities - a.Activities
			})

			guessedPersons = append(guessedPersons, guessedPerson)
		}
	}

	_ = d.locMatcher().Match(e.Desc)

	// TODO: record locs/friends

	tags := lang.ExtractTags(e.Desc)

	if len(tags) > 0 {
		d.AddTags(tags)
		tag.Add(&e, tags)
	}

	e.Friends = make([]string, 0, len(certainPersons)+len(guessedPersons))

	for _, p := range certainPersons {
		e.Friends = append(e.Friends, p.Name)
		p.Activities++
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
