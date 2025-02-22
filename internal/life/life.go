package life

import (
	"fmt"
	"strings"

	"github.com/roma-glushko/frens/internal/event"
	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/location"
	"github.com/roma-glushko/frens/internal/tag"
)

type ListFriendQuery struct {
	Location string
	Tag      string
}

type Data struct {
	dirty      bool
	Tags       []tag.Tag
	Friends    friend.Friends
	Locations  location.Locations
	Activities []event.Activity
}

func (d *Data) Init() {
	// TODO: implement
}

func (d *Data) Dirty() bool {
	return d.dirty
}

func (d *Data) AddFriend(f friend.Friend) {
	d.Friends = append(d.Friends, f)

	d.dirty = true
}

func (d *Data) GetFriend(q string) (*friend.Friend, error) {
	var found []friend.Friend

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

func (d *Data) AddLocation(l location.Location) {
	d.Locations = append(d.Locations, l)

	d.dirty = true
}

func (d *Data) GetLocation(q string) (*location.Location, error) {
	var found []location.Location

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

func (d *Data) AddTag(t tag.Tag) {
	d.Tags = append(d.Tags, t)

	d.dirty = true
}

func (d *Data) AddActivity(a event.Activity) {
	d.Activities = append(d.Activities, a)

	d.dirty = true
}

func (d *Data) ListFriends(q ListFriendQuery) []friend.Friend {
	var v []friend.Friend

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
