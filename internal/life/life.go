package life

import (
	"github.com/roma-glushko/frens/internal/activity"
	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/location"
	"github.com/roma-glushko/frens/internal/tag"
)

type Data struct {
	Tags       []tag.Tag
	Friends    []friend.Friend
	Locations  []location.Location
	Activities []activity.Activity
}

func (d *Data) AddFriend(f friend.Friend) {
	d.Friends = append(d.Friends, f)
}
