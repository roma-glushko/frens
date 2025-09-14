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
	"sort"
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

type Stats struct {
	Friends    int `json:"friends"`
	Locations  int `json:"locations"`
	Activities int `json:"activities"`
	Notes      int `json:"notes"`
}

type Journal struct {
	DirPath    string
	Tags       tag.Tags
	Friends    []*friend.Person
	Locations  friend.Locations
	Activities []*friend.Event
	Notes      []*friend.Event

	dirty           bool
	matcherMu       sync.Mutex
	friendMatcher   *matcher.Matcher[friend.Person]
	locationMatcher *matcher.Matcher[friend.Location]
}

func (j *Journal) Init() {
	j.matcherMu.Lock()
	defer j.matcherMu.Unlock()

	if j.friendMatcher == nil {
		j.friendMatcher = matcher.NewMatcher[friend.Person]()

		for _, f := range j.Friends {
			j.friendMatcher.Add(f)
		}
	}

	if j.locationMatcher == nil {
		j.locationMatcher = matcher.NewMatcher[friend.Location]()

		for _, l := range j.Locations {
			j.locationMatcher.Add(l)
		}
	}
}

func (j *Journal) IsDirty() bool {
	return j.dirty
}

func (j *Journal) SetDirty(d bool) {
	j.dirty = d
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

	j.Friends = append(j.Friends, &f)
	j.SetDirty(true)
}

func (j *Journal) GetFriend(q string) (friend.Person, error) {
	matches := j.frenMatcher().Match(q)

	if len(matches) == 0 {
		return friend.Person{}, fmt.Errorf("no friends found for '%s'", q)
	}

	if len(matches) > 1 {
		names := make([]string, 0, len(matches))

		for _, m := range matches {
			for _, f := range m.Entities {
				names = append(names, f.Name)
			}
		}

		return friend.Person{}, fmt.Errorf(
			"multiple friends found for '%s': %s",
			q,
			strings.Join(names, ", "),
		)
	}

	m := matches[0]

	if len(m.Entities) == 0 {
		return friend.Person{}, fmt.Errorf("no friend found for '%s'", q)
	}

	if len(m.Entities) > 1 {
		names := make([]string, 0, len(m.Entities))

		for _, f := range m.Entities {
			names = append(names, f.Name)
		}

		return friend.Person{}, fmt.Errorf(
			"multiple friends found for '%s': %s",
			q,
			strings.Join(names, ", "),
		)
	}

	return *m.Entities[0], nil
}

func (j *Journal) ListFriends(q friend.ListFriendQuery) []friend.Person { //nolint:cyclop
	fl := make([]friend.Person, 0, 10)

	for _, f := range j.Friends {
		if q.Keyword != "" &&
			!strings.Contains(strings.ToLower(f.Name), strings.ToLower(q.Keyword)) &&
			!strings.Contains(strings.ToLower(f.Desc), strings.ToLower(q.Keyword)) {
			continue
		}

		if len(q.Locations) > 0 && !f.HasLocations(q.Locations) {
			continue
		}

		if len(q.Tags) > 0 && !tag.HasTags(f, q.Tags) {
			continue
		}

		fl = append(fl, *f)
	}

	if len(fl) == 0 {
		return fl
	}

	// sort by and order by friends
	sort.SliceStable(fl, func(i, j int) bool {
		switch q.SortBy {
		case friend.SortAlpha:
			if q.SortOrder == friend.SortOrderDirect {
				return strings.ToLower(fl[i].Name) < strings.ToLower(fl[j].Name)
			}

			return strings.ToLower(fl[i].Name) > strings.ToLower(fl[j].Name)
		case friend.SortActivities:
			if q.SortOrder == friend.SortOrderDirect {
				return fl[i].Activities < fl[j].Activities
			}

			return fl[i].Activities > fl[j].Activities
		case friend.SortRecency:
			if q.SortOrder == friend.SortOrderDirect {
				return fl[i].MostRecentActivity.After(fl[j].MostRecentActivity)
			}

			return fl[i].MostRecentActivity.Before(fl[j].MostRecentActivity)
		default:
			return false
		}
	})

	return fl
}

func (j *Journal) UpdateFriend(o, n friend.Person) {
	if n.ID == "" {
		n.ID = o.ID
	}

	n.Activities = o.Activities
	n.Notes = o.Notes
	n.MostRecentActivity = o.MostRecentActivity

	for i, f := range j.Friends {
		if f.Name == o.Name {
			j.Friends[i] = &n
			j.SetDirty(true)

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
				j.SetDirty(true)

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

	j.Locations = append(j.Locations, &l)
	j.SetDirty(true)
}

func (j *Journal) GetLocation(q string) (friend.Location, error) {
	matches := j.locMatcher().Match(q)

	if len(matches) == 0 {
		return friend.Location{}, fmt.Errorf("no locations found for '%s'", q)
	}

	if len(matches) > 1 {
		names := make([]string, 0, len(matches))

		for _, m := range matches {
			for _, f := range m.Entities {
				names = append(names, f.Name)
			}
		}

		return friend.Location{}, fmt.Errorf(
			"multiple locations found for '%s': %s",
			q,
			strings.Join(names, ", "),
		)
	}

	m := matches[0]

	if len(m.Entities) == 0 {
		return friend.Location{}, fmt.Errorf("no locations found for '%s'", q)
	}

	if len(m.Entities) > 1 {
		names := make([]string, 0, len(m.Entities))

		for _, f := range m.Entities {
			names = append(names, f.Name)
		}

		return friend.Location{}, fmt.Errorf(
			"multiple locations found for '%s': %s",
			q,
			strings.Join(names, ", "),
		)
	}

	return *m.Entities[0], nil
}

func (j *Journal) UpdateLocation(o, n friend.Location) {
	if n.ID == "" {
		n.ID = o.ID
	}

	n.Activities = o.Activities
	n.MostRecentActivity = o.MostRecentActivity

	for i, l := range j.Locations {
		if l.Name == o.Name {
			j.Locations[i] = &n
			j.SetDirty(true)

			return
		}
	}

	// TODO: update friend references in activities and notes

	// If the friend was not found, add it as a new one
	j.AddLocation(n)
}

func (j *Journal) ListLocations(q friend.ListLocationQuery) []friend.Location { //nolint:cyclop
	locations := make([]friend.Location, 0, 10)

	for _, l := range j.Locations {
		if q.Keyword != "" &&
			!strings.Contains(strings.ToLower(l.Name), strings.ToLower(q.Keyword)) &&
			!strings.Contains(strings.ToLower(l.Desc), strings.ToLower(q.Keyword)) {
			continue
		}

		if len(q.Countries) > 0 && !slices.Contains(q.Countries, l.Country) {
			continue
		}

		if len(q.Tags) > 0 && !tag.HasTags(l, q.Tags) {
			continue
		}

		locations = append(locations, *l)
	}

	if len(locations) == 0 {
		return locations
	}

	sort.SliceStable(locations, func(i, j int) bool {
		switch q.SortBy {
		case friend.SortAlpha:
			if q.SortOrder == friend.SortOrderReverse {
				return strings.ToLower(locations[i].Name) > strings.ToLower(locations[j].Name)
			}

			return strings.ToLower(locations[i].Name) < strings.ToLower(locations[j].Name)
		case friend.SortActivities:
			if q.SortOrder == friend.SortOrderReverse {
				return locations[i].Activities < locations[j].Activities
			}

			return locations[i].Activities > locations[j].Activities
		case friend.SortRecency:
			if q.SortOrder == friend.SortOrderReverse {
				return locations[i].MostRecentActivity.After(locations[j].MostRecentActivity)
			}

			return locations[i].MostRecentActivity.Before(locations[j].MostRecentActivity)
		default:
			return false
		}
	})

	return locations
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

	j.SetDirty(true)
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
		KnownPerson       *friend.Person
		AmbiguitiesPerson *friend.Person
	}

	guessedPersons := make([]*friend.Person, 0, len(ambiguitiesMatches))

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

		for _, act := range j.Activities {
			for _, pair := range rankPairs {
				if slices.Contains(act.FriendIDs, pair.KnownPerson.Name) &&
					slices.Contains(act.FriendIDs, pair.AmbiguitiesPerson.Name) {
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

	e.FriendIDs = make([]string, 0, len(guessedPersons))

	for _, p := range guessedPersons {
		e.FriendIDs = append(e.FriendIDs, p.ID)

		if e.Type == friend.EventTypeActivity {
			p.MostRecentActivity = e.Date
			p.Activities++
		} else {
			p.Notes++
		}
	}

	switch e.Type {
	case friend.EventTypeActivity:
		j.Activities = append(j.Activities, &e)
	case friend.EventTypeNote:
		j.Notes = append(j.Notes, &e)
	default:
		return friend.Event{}, fmt.Errorf("unknown event type: %s", e.Type)
	}

	j.SetDirty(true)

	return e, nil
}

func (j *Journal) GetEvent(t friend.EventType, q string) (friend.Event, error) {
	if t == friend.EventTypeActivity {
		for _, act := range j.Activities {
			if act.ID == q {
				return *act, nil
			}
		}
	}

	if t == friend.EventTypeNote {
		for _, note := range j.Notes {
			if note.ID == q {
				return *note, nil
			}
		}
	}

	return friend.Event{}, ErrEventNotFound
}

func (j *Journal) UpdateEvent(o, n friend.Event) (friend.Event, error) {
	n.ID = o.ID

	guessedPersons := j.GuessFriends(n.Desc)
	tags := lang.ExtractTags(n.Desc)
	_ = j.locMatcher().Match(n.Desc) // TODO: parse location

	if len(tags) > 0 {
		j.AddTags(tags)
		tag.Add(&n, tags)
	}

	for _, p := range guessedPersons {
		n.FriendIDs = append(n.FriendIDs, p.ID)

		if n.Type == friend.EventTypeActivity {
			p.MostRecentActivity = n.Date
			p.Activities++
		} else {
			p.Notes++
		}
	}

	if o.Type == friend.EventTypeActivity {
		for i, act := range j.Activities {
			if act.ID == o.ID {
				j.Activities[i] = &n
				j.SetDirty(true)

				return n, nil
			}
		}
	}

	if o.Type == friend.EventTypeNote {
		for i, note := range j.Notes {
			if note.ID == o.ID {
				j.Notes[i] = &n
				j.SetDirty(true)

				return n, nil
			}
		}
	}

	// If the activity was not found, add it as a new one
	return j.AddEvent(n)
}

func (j *Journal) ListEvents(q friend.ListEventQuery) ([]friend.Event, error) { //nolint:cyclop
	events := make([]friend.Event, 0, 10)

	var source []*friend.Event

	switch q.Type {
	case friend.EventTypeActivity:
		source = j.Activities
	case friend.EventTypeNote:
		source = j.Notes
	default:
		return events, fmt.Errorf("unknown event type: %s", q.Type)
	}

	for _, note := range source {
		if q.Keyword != "" &&
			!strings.Contains(strings.ToLower(note.Desc), strings.ToLower(q.Keyword)) {
			continue
		}

		if len(q.Tags) > 0 && !tag.HasTags(note, q.Tags) {
			continue
		}

		if !q.Since.IsZero() && note.Date.Before(q.Since) {
			continue
		}

		if !q.Until.IsZero() && note.Date.After(q.Until) {
			continue
		}

		events = append(events, *note)
	}

	if len(events) == 0 {
		return events, nil
	}

	sort.SliceStable(events, func(i, j int) bool {
		switch q.SortBy { //nolint:exhaustive
		case friend.SortAlpha:
			if q.SortOrder == friend.SortOrderReverse {
				return strings.ToLower(events[i].Desc) > strings.ToLower(events[j].Desc)
			}

			return strings.ToLower(events[i].Desc) < strings.ToLower(events[j].Desc)
		case friend.SortRecency:
			if q.SortOrder == friend.SortOrderReverse {
				return events[i].Date.After(events[j].Date)
			}

			return events[i].Date.Before(events[j].Date)
		default:
			return false
		}
	})

	return events, nil
}

func (j *Journal) RemoveEvents(t friend.EventType, toRemove []friend.Event) {
	for _, act := range toRemove {
		if t == friend.EventTypeActivity {
			for i, a := range j.Activities {
				if a.ID == act.ID {
					j.Activities = append(j.Activities[:i], j.Activities[i+1:]...)
					j.SetDirty(true)

					break
				}
			}
		}

		if t == friend.EventTypeNote {
			for i, n := range j.Notes {
				if n.ID == act.ID {
					j.Notes = append(j.Notes[:i], j.Notes[i+1:]...)
					j.SetDirty(true)

					break
				}
			}
		}
	}
}

func (j *Journal) AddFriendDate(fID string, d friend.Date) (friend.Date, error) {
	f, err := j.GetFriend(fID)
	if err != nil {
		return friend.Date{}, fmt.Errorf("failed to get friend %s: %w", fID, err)
	}

	if d.ID == "" {
		d.ID = ksuid.New().String()
	}

	f.Dates = append(f.Dates, &d)

	j.SetDirty(true)

	return d, nil
}

func (j *Journal) UpdateFriendDate(o, n friend.Date) (friend.Date, error) {
	n.ID = o.ID

	for _, f := range j.Friends {
		for i, d := range f.Dates {
			if d.ID == o.ID {
				f.Dates[i] = &n

				j.SetDirty(true)

				return n, nil
			}
		}
	}

	return friend.Date{}, fmt.Errorf("date with ID %s not found", o.ID)
}

func (j *Journal) GetFriendDate(dID string) (friend.Date, error) {
	for _, f := range j.Friends {
		for _, d := range f.Dates {
			if d.ID == dID {
				d.Person = f.ID

				return *d, nil
			}
		}
	}

	return friend.Date{}, fmt.Errorf("date with ID %s not found", dID)
}

func (j *Journal) ListFriendDates(q friend.ListDateQuery) ([]friend.Date, error) { //nolint:cyclop
	dates := make([]friend.Date, 0, 10)

	frs := j.Friends

	if len(q.Friends) > 0 {
		frs = make([]*friend.Person, 0, len(q.Friends))

		for _, fID := range q.Friends {
			f, err := j.GetFriend(fID)
			if err != nil {
				return dates, fmt.Errorf("failed to get friend %s: %w", fID, err)
			}

			frs = append(frs, &f)
		}
	}

	for _, f := range frs {
		for _, d := range f.Dates {
			if q.Keyword != "" &&
				!strings.Contains(strings.ToLower(d.Desc), strings.ToLower(q.Keyword)) {
				continue
			}

			if len(q.Tags) > 0 && !tag.HasTags(d, q.Tags) {
				continue
			}

			d.Person = f.ID

			dates = append(dates, *d)
		}
	}

	if len(dates) == 0 {
		return dates, nil
	}

	// TODO: sorting by date?

	return dates, nil
}

func (j *Journal) RemoveFriendDates(toRemove []friend.Date) error {
	for _, dID := range toRemove {
		found := false

		for _, f := range j.Friends {
			for i, d := range f.Dates {
				if d.ID == dID.ID {
					f.Dates = append(f.Dates[:i], f.Dates[i+1:]...)
					found = true

					j.SetDirty(true)

					break
				}
			}

			if found {
				break
			}
		}

		if !found {
			return fmt.Errorf("date with ID %s not found", dID.ID)
		}
	}

	return nil
}

func (j *Journal) Stats() Stats {
	return Stats{
		Friends:    len(j.Friends),
		Locations:  len(j.Locations),
		Activities: len(j.Activities),
		Notes:      len(j.Notes),
	}
}

func (j *Journal) locMatcher() *matcher.Matcher[friend.Location] {
	return j.locationMatcher
}

func (j *Journal) frenMatcher() *matcher.Matcher[friend.Person] {
	return j.friendMatcher
}
