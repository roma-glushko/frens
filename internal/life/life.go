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

package life

import (
	"fmt"
	"strings"

	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/tag"
)

type ListFriendQuery struct {
	Location string
	Tag      string
}

type Data struct {
	dirty      bool
	DirPath    string
	Tags       tag.Tags
	Friends    []friend.Person
	Locations  friend.Locations
	Activities []friend.Event
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
	var found []friend.Person

	for _, f := range d.Friends {
		if f.Match(q) {
			found = append(found, f)
		}
	}

	if len(found) == 0 {
		return nil, fmt.Errorf("no friends found for '%s'", q)
	}

	if len(found) > 1 {
		names := make([]string, len(found))

		for i, f := range found {
			names[i] = f.Name
		}

		return nil, fmt.Errorf("multiple friends found for '%s': %s", q, strings.Join(names, ", "))
	}

	f := found[0]

	return &f, nil
}

func (d *Data) AddLocation(l friend.Location) {
	d.Locations = append(d.Locations, l)

	d.dirty = true
}

func (d *Data) GetLocation(q string) (*friend.Location, error) {
	var found []friend.Location

	for _, l := range d.Locations {
		if l.Match(q) {
			found = append(found, l)
		}
	}

	if len(found) == 0 {
		return nil, fmt.Errorf("no locations found for '%s'", q)
	}

	if len(found) > 1 {
		names := make([]string, len(found))

		for i, f := range found {
			names[i] = f.Name
		}

		return nil, fmt.Errorf("multiple locations found for '%s': %s", q, strings.Join(names, ", "))
	}

	l := found[0]

	return &l, nil
}

func (d *Data) AddTags(t []tag.Tag) {
	d.Tags = append(d.Tags, t...).Unique()

	d.dirty = true
}

func (d *Data) AddActivity(e friend.Event) {
	//friendMatcher := utils.NewMatcher[friend.Location]()
	//
	//for _, l := range d.Locations {
	//	friendMatcher.Add(l, l.Refs())
	//}
	//
	//friends := friendMatcher.Match(e.Desc)

	//locMatcher := utils.NewMatcher[friend.Location]()
	//
	//for _, l := range d.Locations {
	//	locMatcher.Add(l, l.Refs())
	//}

	//locations := locMatcher.Match(e.Desc)

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
