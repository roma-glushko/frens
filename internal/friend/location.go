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
	"sort"
	"strings"

	"github.com/roma-glushko/frens/internal/tag"

	"github.com/roma-glushko/frens/internal/utils"
)

var ErrLocNameEmpty = errors.New("location name must be provided")

type Location struct {
	Name    string   `toml:"name"`
	Country string   `toml:"country,omitempty"`
	Desc    string   `toml:"desc"`
	Aliases []string `toml:"aliases,omitempty"`
	Tags    []string `toml:"tags,omitempty"`

	Activities int `toml:"activities,omitempty"`
}

var (
	_ tag.Tagged      = (*Location)(nil)
	_ utils.Matchable = (*Location)(nil)
)

func (l *Location) Validate() error {
	if l.Name == "" {
		return ErrLocNameEmpty
	}

	return nil
}

func (l Location) Refs() []string {
	names := make([]string, 0, 1+len(l.Aliases))

	names = append(names, l.Name)

	if len(l.Aliases) > 0 {
		names = append(names, l.Aliases...)
	}

	return names
}

func (l *Location) AddAlias(a string) {
	l.Aliases = utils.Unique(append(l.Aliases, a))
}

func (l *Location) RemoveAlias(a string) {
	var aliases []string

	for _, alias := range l.Aliases {
		if !strings.EqualFold(alias, a) {
			aliases = append(aliases, alias)
		}
	}

	l.Aliases = aliases
}

func (l *Location) SetTags(tags []string) {
	l.Tags = tags
}

func (l *Location) GetTags() []string {
	return l.Tags
}

func (l *Location) String() string {
	var sb strings.Builder

	sb.WriteString(l.Name)

	if len(l.Country) > 0 {
		sb.WriteString(", " + l.Country)
	}

	return sb.String()
}

type Locations []Location

var _ sort.Interface = (*Locations)(nil)

func (l Locations) Len() int           { return len(l) }
func (l Locations) Less(i, j int) bool { return l[i].Name < l[j].Name }
func (l Locations) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
