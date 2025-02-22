package toml

import (
	"github.com/roma-glushko/frens/internal/event"
	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/location"
	"github.com/roma-glushko/frens/internal/tag"
)

const (
	FileNameFriends    = "friends.toml"
	FileNameActivities = "activities.toml"
)

type FriendsFile struct {
	Tags      []tag.Tag           `toml:"tags"`
	Friends   []friend.Friend     `toml:"friends"`
	Locations []location.Location `toml:"locations"`
}

type ActivitiesFile struct {
	Activities []event.Activity `toml:"activities"`
}
