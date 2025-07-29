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

package friend

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/roma-glushko/frens/internal/matcher"

	"github.com/roma-glushko/frens/internal/tag"

	"github.com/roma-glushko/frens/internal/utils"
)

var ErrFriendNameEmpty = errors.New("friend name must be provided")

type Persons []Person

var _ sort.Interface = (*Persons)(nil)

func (p Persons) Len() int           { return len(p) }
func (p Persons) Less(i, j int) bool { return p[i].Name < p[j].Name }
func (p Persons) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type Person struct {
	ID        string   `toml:"id"`
	Name      string   `toml:"name"`
	Desc      string   `toml:"desc,omitempty"`
	Nicknames []string `toml:"nicknames,omitempty"`
	Tags      []string `toml:"tags,omitempty"`
	Locations []string `toml:"locations,omitempty"`
	Dates     []*Date  `toml:"dates,omitempty"`
	// Cached information
	Activities         int       `toml:"activities,omitempty"`
	Notes              int       `toml:"notes,omitempty"`
	MostRecentActivity time.Time `toml:"most_recent_activity,omitzero"`
	// internal use only
	Score int `toml:"-"`
}

var (
	_ tag.Tagged        = (*Person)(nil)
	_ matcher.Matchable = (*Person)(nil)
)

func (p *Person) Validate() error {
	if p.Name == "" {
		return ErrFriendNameEmpty
	}

	return nil
}

func (p Person) Refs() []string {
	names := make([]string, 0, 3+len(p.Nicknames))

	names = append(names, p.Name)

	parts := strings.Split(p.Name, " ")

	if len(p.Nicknames) > 0 {
		names = append(names, p.Nicknames...)
	}

	if len(parts) >= 2 {
		// Add first name and last name as separate references
		names = append(names, parts[0], parts[len(parts)-1])
	}

	if p.ID != "" {
		names = append(names, p.ID)
	}

	for i, c := 0, len(names); i < c; i++ {
		names[i] = strings.ToLower(names[i])
	}

	return utils.Unique(names)
}

func (p *Person) AddNickname(n string) {
	p.Nicknames = utils.Unique(append(p.Nicknames, n))
}

func (p *Person) RemoveNickname(n string) {
	var nicks []string

	for _, nick := range p.Nicknames {
		if !strings.EqualFold(nick, n) {
			nicks = append(nicks, nick)
		}
	}

	p.Nicknames = nicks
}

func (p *Person) SetTags(tags []string) {
	p.Tags = tags
}

func (p *Person) GetTags() []string {
	return p.Tags
}

func (p *Person) HasLocations(ls []string) bool {
	for _, loc := range p.Locations {
		for _, l := range ls {
			if strings.EqualFold(loc, l) {
				return true
			}
		}
	}

	return false
}

func (p *Person) AddLocation(l string) {
	p.Locations = utils.Unique(append(p.Locations, l))
}

func (p *Person) RemoveLocation(l string) {
	var locs []string

	for _, loc := range p.Locations {
		if !strings.EqualFold(loc, l) {
			locs = append(locs, loc)
		}
	}

	p.Locations = locs
}

func (p *Person) String() string {
	var sb strings.Builder

	sb.WriteString(p.Name)

	if len(p.Nicknames) > 0 {
		sb.WriteString(fmt.Sprintf(" (a.k.a %s)", strings.Join(p.Nicknames, ", ")))
	}

	return sb.String()
}
