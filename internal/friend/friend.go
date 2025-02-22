package friend

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/roma-glushko/frens/internal/utils"
)

var ErrNameEmpty = errors.New("friend name must be provided")

type Friend struct {
	Name      string    `toml:"name"`
	Birthday  time.Time `toml:"birthday,omitempty"`
	Nicknames []string  `toml:"nicknames,omitempty"`
	Tags      []string  `toml:"tags,omitempty"`
	Locations []string  `toml:"locations,omitempty"`
	Reminders []string  `toml:"reminders,omitempty"`
	// activities int
}

func (f *Friend) Validate() error {
	if f.Name == "" {
		return ErrNameEmpty
	}

	return nil
}

func (f *Friend) Match(q string) bool {
	q = strings.ToLower(q)

	if strings.Contains(strings.ToLower(f.Name), q) {
		return true
	}

	for _, n := range f.Nicknames {
		if strings.Contains(strings.ToLower(n), q) {
			return true
		}
	}

	return false
}

func (f *Friend) AddNickname(n string) {
	f.Nicknames = utils.Unique(append(f.Nicknames, n))
}

func (f *Friend) RemoveNickname(n string) {
	var nicks []string

	for _, nick := range f.Nicknames {
		if strings.ToLower(nick) != strings.ToLower(n) {
			nicks = append(nicks, nick)
		}
	}

	f.Nicknames = nicks
}

func (f *Friend) AddTag(t string) {
	f.Tags = utils.Unique(append(f.Tags, t))
}

func (f *Friend) HasTag(t string) bool {
	for _, tag := range f.Tags {
		if strings.ToLower(tag) == strings.ToLower(t) {
			return true
		}
	}

	return false
}

func (f *Friend) RemoveTag(t string) {
	var tags []string

	for _, tag := range f.Tags {
		if strings.ToLower(tag) != strings.ToLower(t) {
			tags = append(tags, tag)
		}
	}

	f.Tags = tags
}

func (f *Friend) HasLocation(l string) bool {
	for _, loc := range f.Locations {
		if strings.ToLower(loc) == strings.ToLower(l) {
			return true
		}
	}

	return false
}

func (f *Friend) AddLocation(l string) {
	f.Locations = utils.Unique(append(f.Locations, l))
}

func (f *Friend) RemoveLocation(l string) {
	var locs []string

	for _, loc := range f.Locations {
		if strings.ToLower(loc) != strings.ToLower(l) {
			locs = append(locs, loc)
		}
	}

	f.Locations = locs
}

func (f *Friend) String() string {
	var sb strings.Builder

	sb.WriteString(f.Name)

	if len(f.Nicknames) > 0 {
		sb.WriteString(fmt.Sprintf(" (a.k.a %s)", strings.Join(f.Nicknames, ", ")))
	}

	return sb.String()
}
