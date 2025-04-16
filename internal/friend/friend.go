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

	"github.com/roma-glushko/frens/internal/utils"
)

var ErrFriendNameEmpty = errors.New("friend name must be provided")

type Friends []Friend

var _ sort.Interface = (*Friends)(nil)

func (f Friends) Len() int           { return len(f) }
func (f Friends) Less(i, j int) bool { return f[i].Name < f[j].Name }
func (f Friends) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

type Friend struct {
	Name       string    `toml:"name"`
	Birthday   time.Time `toml:"birthday,omitempty"`
	Nicknames  []string  `toml:"nicknames,omitempty"`
	Tags       []string  `toml:"tags,omitempty"`
	Notes      []Event   `toml:"notes,omitempty"`
	Locations  []string  `toml:"locations,omitempty"`
	Reminders  []string  `toml:"reminders,omitempty"`
	Activities int       `toml:"activities,omitempty"`
}

var _ Tagged = (*Friend)(nil)

func (f *Friend) Validate() error {
	if f.Name == "" {
		return ErrFriendNameEmpty
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
		if !strings.EqualFold(nick, n) {
			nicks = append(nicks, nick)
		}
	}

	f.Nicknames = nicks
}

func (f *Friend) SetTags(tags []string) {
	f.Tags = tags
}

func (f *Friend) GetTags() []string {
	return f.Tags
}

func (f *Friend) HasLocation(l string) bool {
	for _, loc := range f.Locations {
		if strings.EqualFold(loc, l) {
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
		if !strings.EqualFold(loc, l) {
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
