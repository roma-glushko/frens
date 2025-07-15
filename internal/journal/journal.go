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
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/gosimple/slug"

	"github.com/roma-glushko/frens/internal/matcher"

	"github.com/segmentio/ksuid"

	"github.com/roma-glushko/frens/internal/lang"

	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/tag"
)

var ErrEventNotFound = errors.New("event not found")

type ListFriendQuery struct {
	Location string
	Tag      string
}

type Journal struct {
	DirPath    string
	Tags       tag.Tags
	Friends    []friend.Person
	Locations  friend.Locations
	Activities []friend.Event
	Notes      []friend.Event

	dirty           bool
	matcherMu       sync.Mutex
	friendMatcher   *matcher.Matcher[friend.Person]
	locationMatcher *matcher.Matcher[friend.Location]
}

func (j *Journal) Init() {
	// TODO: implement
}

func (j *Journal) Dirty() bool {
	return j.dirty
}

func (j *Journal) Path() string {
	return j.DirPath
}

func (j *Journal) AddFriend(f friend.Person) {
	if f.ID == "" {
		f.ID = slug.Make(f.Name)
	}

	// TODO: check for duplicated IDs
	// TODO: check for duplicated aliases

	j.Friends = append(j.Friends, f)
	j.dirty = true
}

func (j *Journal) GetFriend(q string) (*friend.Person, error) {
	matches := j.frenMatcher().Match(q)

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

	return m.Entities[0], nil
}

func (j *Journal) UpdateFriend(o, n friend.Person) {
	for i, f := range j.Friends {
		if f.Name == o.Name {
			j.Friends[i] = n
			j.dirty = true

			return
		}
	}

	// TODO: update friend references in activities and notes

	// If the friend was not found, add it as a new one
	j.AddFriend(n)
}

func (j *Journal) RemoveFriends(toRemove []friend.Person) {
	for _, fr := range toRemove {
		for i, f := range j.Friends {
			if f.Name == fr.Name {
				j.Friends = append(j.Friends[:i], j.Friends[i+1:]...)
				j.dirty = true

				break
			}
		}
	}
}

func (j *Journal) AddLocation(l friend.Location) {
	if l.ID == "" {
		l.ID = slug.Make(l.Name)
	}

	// TODO: check for duplicated IDs
	// TODO: check for duplicated aliases

	j.Locations = append(j.Locations, l)
	j.dirty = true
}

func (j *Journal) GetLocation(q string) (*friend.Location, error) {
	matches := j.locMatcher().Match(q)

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

		return nil, fmt.Errorf(
			"multiple locations found for '%s': %s",
			q,
			strings.Join(names, ", "),
		)
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

		return nil, fmt.Errorf(
			"multiple locations found for '%s': %s",
			q,
			strings.Join(names, ", "),
		)
	}

	return m.Entities[0], nil
}

func (j *Journal) UpdateLocation(o, n friend.Location) {
	for i, l := range j.Locations {
		if l.Name == o.Name {
			j.Locations[i] = n
			j.dirty = true

			return
		}
	}

	// TODO: update friend references in activities and notes

	// If the friend was not found, add it as a new one
	j.AddLocation(n)
}

func (j *Journal) RemoveLocations(toRemove []friend.Location) {
	for _, loc := range toRemove {
		for i, l := range j.Locations {
			if l.Name == loc.Name {
				j.Locations = append(j.Locations[:i], j.Locations[i+1:]...)
				j.dirty = true

				break
			}
		}
	}
}

func (j *Journal) AddTags(t []tag.Tag) {
	j.Tags = append(j.Tags, t...).Unique()

	j.dirty = true
}

func (j *Journal) GuessFriends(q string) []*friend.Person { //nolint:cyclop
	matches := j.frenMatcher().Match(q)

	certainPersons := make([]*friend.Person, 0, len(matches))
	ambiguitiesMatches := make([]matcher.Match[friend.Person], 0, len(matches))

	for _, m := range matches {
		if len(m.Entities) == 1 {
			certainPersons = append(certainPersons, m.Entities[0])
			continue
		}

		shortestNameFriend := slices.MinFunc(m.Entities, func(a, b *friend.Person) int {
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

	guessedPersons := make([]*friend.Person, 0, len(ambiguitiesMatches))

	if len(ambiguitiesMatches) > 0 {
		rankPairs := make([]friendPair, 0, len(certainPersons)*len(ambiguitiesMatches))

		for _, cp := range certainPersons {
			for _, am := range ambiguitiesMatches {
				for _, ap := range am.Entities {
					rankPairs = append(rankPairs, friendPair{
						KnownPerson:       *cp,
						AmbiguitiesPerson: *ap,
					})
				}
			}
		}

		for _, act := range j.Activities {
			for _, pair := range rankPairs {
				if slices.Contains(act.Friends, pair.KnownPerson.Name) &&
					slices.Contains(act.Friends, pair.AmbiguitiesPerson.Name) {
					pair.AmbiguitiesPerson.Score++
				}
			}
		}

		for _, am := range ambiguitiesMatches {
			guessedPerson := slices.MaxFunc(am.Entities, func(a, b *friend.Person) int {
				if a.Score != b.Score {
					return b.Score - a.Score
				}

				return b.Activities - a.Activities
			})

			guessedPersons = append(guessedPersons, guessedPerson)
		}
	}

	return append(certainPersons, guessedPersons...)
}

func (j *Journal) AddEvent(e friend.Event) (friend.Event, error) {
	e.ID = ksuid.New().String()

	guessedPersons := j.GuessFriends(e.Desc)

	_ = j.locMatcher().Match(e.Desc)

	// TODO: record locs/friends

	tags := lang.ExtractTags(e.Desc)

	if len(tags) > 0 {
		j.AddTags(tags)
		tag.Add(&e, tags)
	}

	e.Friends = make([]string, 0, len(guessedPersons))

	for _, p := range guessedPersons {
		e.Friends = append(e.Friends, p.ID)

		if e.Type == friend.EventTypeActivity {
			p.Activities++
		} else {
			p.Notes++
		}
	}

	switch e.Type {
	case friend.EventTypeActivity:
		j.Activities = append(j.Activities, e)
	case friend.EventTypeNote:
		j.Notes = append(j.Notes, e)
	default:
		return friend.Event{}, fmt.Errorf("unknown event type: %s", e.Type)
	}

	j.dirty = true

	return e, nil
}

func (j *Journal) GetEvent(t friend.EventType, q string) (friend.Event, error) {
	if t == friend.EventTypeActivity {
		for _, act := range j.Activities {
			if act.ID == q {
				return act, nil
			}
		}
	}

	if t == friend.EventTypeNote {
		for _, note := range j.Notes {
			if note.ID == q {
				return note, nil
			}
		}
	}

	return friend.Event{}, ErrEventNotFound
}

func (j *Journal) UpdateEvent(o, n friend.Event) (friend.Event, error) {
	n.ID = o.ID

	if o.Type == friend.EventTypeActivity {
		for i, act := range j.Activities {
			if act.ID == o.ID {
				j.Activities[i] = n
				j.dirty = true

				return n, nil
			}
		}
	}

	if o.Type == friend.EventTypeNote {
		for i, note := range j.Notes {
			if note.ID == o.ID {
				j.Notes[i] = n
				j.dirty = true

				return n, nil
			}
		}
	}

	// TODO: update friend & location references

	// If the activity was not found, add it as a new one
	return j.AddEvent(n)
}

func (j *Journal) RemoveEvents(t friend.EventType, toRemove []friend.Event) {
	for _, act := range toRemove {
		if t == friend.EventTypeActivity {
			for i, a := range j.Activities {
				if a.ID == act.ID {
					j.Activities = append(j.Activities[:i], j.Activities[i+1:]...)
					j.dirty = true

					break
				}
			}
		}

		if t == friend.EventTypeNote {
			for i, n := range j.Notes {
				if n.ID == act.ID {
					j.Notes = append(j.Notes[:i], j.Notes[i+1:]...)
					j.dirty = true

					break
				}
			}
		}
	}
}

func (j *Journal) ListFriends(q ListFriendQuery) []friend.Person {
	v := make([]friend.Person, 0, 5)

	for _, f := range j.Friends {
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

func (j *Journal) locMatcher() *matcher.Matcher[friend.Location] {
	if j.locationMatcher == nil {
		j.matcherMu.Lock()
		defer j.matcherMu.Unlock()

		j.locationMatcher = matcher.NewMatcher[friend.Location]()

		for _, l := range j.Locations {
			j.locationMatcher.Add(&l)
		}
	}

	return j.locationMatcher
}

func (j *Journal) frenMatcher() *matcher.Matcher[friend.Person] {
	if j.friendMatcher == nil {
		j.matcherMu.Lock()
		defer j.matcherMu.Unlock()

		j.friendMatcher = matcher.NewMatcher[friend.Person]()

		for _, f := range j.Friends {
			j.friendMatcher.Add(&f)
		}
	}

	return j.friendMatcher
}
