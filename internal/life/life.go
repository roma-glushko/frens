package life

import (
	"github.com/roma-glushko/frens/internal/activity"
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
	Friends    []friend.Friend
	Locations  []location.Location
	Activities []activity.Activity
}

func (d *Data) Dirty() bool {
	return d.dirty
}

func (d *Data) AddFriend(f friend.Friend) {
	d.Friends = append(d.Friends, f)

	d.dirty = true
}

func (d *Data) AddLocation(l location.Location) {
	d.Locations = append(d.Locations, l)

	d.dirty = true
}

func (d *Data) AddTag(t tag.Tag) {
	d.Tags = append(d.Tags, t)

	d.dirty = true
}

func (d *Data) AddActivity(a activity.Activity) {
	d.Activities = append(d.Activities, a)

	d.dirty = true
}

func (d *Data) ListFriends(q ListFriendQuery) []friend.Friend {
	view := make([]friend.Friend, 5)

	for _, f := range d.Friends {
		if q.Location != "" {
			if !f.HasLocation(q.Location) {
				continue
			}
		}

		if q.Tag != "" {
			if !f.HasTag(q.Tag) {
				continue
			}
		}

		view = append(view, f)
	}

	return view
}
