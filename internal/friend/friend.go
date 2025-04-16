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
	"slices"
	"strings"
	"time"
)

var ErrNameEmpty = errors.New("friend name must be provided")

type Friend struct {
	Name      string    `toml:"name"`
	Birthday  time.Time `toml:"birthday,omitempty"`
	Nicknames []string  `toml:"nicknames,omitempty"`
	Tags      []string  `toml:"tags,omitempty"`
	Locations []string  `toml:"locations,omitempty"`
	Reminders []string  `toml:"reminders,omitempty"`
}

func (f *Friend) Validate() error {
	if f.Name == "" {
		return ErrNameEmpty
	}

	return nil
}

func (f *Friend) AddNickname(n string) {
	f.Nicknames = append(f.Nicknames, n)
}

func (f *Friend) AddTag(t string) {
	f.Tags = append(f.Tags, t)
}

func (f *Friend) HasTag(t string) bool {
	return slices.Contains(f.Tags, t)
}

func (f *Friend) HasLocation(l string) bool {
	return slices.Contains(f.Locations, l)
}

func (f *Friend) AddLocation(l string) {
	f.Locations = append(f.Locations, l)
}

func (f *Friend) String() string {
	var sb strings.Builder

	sb.WriteString(f.Name)

	if len(f.Nicknames) > 0 {
		sb.WriteString(fmt.Sprintf(" (a.k.a %s)", strings.Join(f.Nicknames, ", ")))
	}

	return sb.String()
}
