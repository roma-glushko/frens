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

package tag

import (
	"slices"
	"sort"
	"strings"

	"github.com/roma-glushko/frens/internal/utils"
)

type Tag struct {
	Name string
}

func NewTag(t string) Tag {
	t = strings.TrimLeft(t, "#")
	t = strings.TrimSpace(t)
	t = strings.ToLower(t)

	return Tag{
		Name: t,
	}
}

func (t *Tag) Equal(other Tag) bool {
	return strings.EqualFold(t.Name, other.Name)
}

func (t *Tag) Match(s string) bool {
	return strings.Contains(strings.ToLower(s), t.Name)
}

func (t *Tag) String() string {
	return "#" + t.Name
}

type Tags []Tag

var _ sort.Interface = (*Tags)(nil)

func (t Tags) Len() int           { return len(t) }
func (t Tags) Less(i, j int) bool { return t[i].Name < t[j].Name }
func (t Tags) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

func (t Tags) ToNames() []string {
	names := make([]string, len(t))

	for i, tag := range t {
		names[i] = tag.Name
	}

	return names
}

func (t Tags) String() string {
	tags := make([]string, 0, len(t))

	for _, tag := range t {
		if tag.Name == "" {
			continue
		}

		tags = append(tags, tag.String())
	}

	slices.Sort(tags)

	return strings.Join(tags, " ")
}

func (t Tags) Unique() Tags {
	return utils.Unique(t)
}

type Tagged interface {
	SetTags(tags []string)
	GetTags() []string
}

func AddStr(e Tagged, tags []string) {
	ts := e.GetTags()

	ts = utils.Unique(append(ts, tags...))

	e.SetTags(tags)
}

func Add(e Tagged, t []Tag) {
	tags := e.GetTags()

	for _, tag := range t {
		tags = append(tags, tag.Name)
	}
}

func Remove(e Tagged, t string) {
	var tags []string

	for _, tag := range e.GetTags() {
		if !strings.EqualFold(tag, t) {
			tags = append(tags, tag)
		}
	}

	e.SetTags(tags)
}

func HasTag(e Tagged, t string) bool {
	for _, tag := range e.GetTags() {
		if strings.EqualFold(tag, t) {
			return true
		}
	}

	return false
}
