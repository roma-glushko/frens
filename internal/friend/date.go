// Copyright 2026 Roma Hlushko
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

	"github.com/roma-glushko/frens/internal/tag"
)

type Calendar = string

var (
	CalendarGregorian Calendar = "gregorian"
	CalendarHebrew    Calendar = "hebrew"
)

type Date struct {
	ID       string   `toml:"id"`
	Calendar Calendar `toml:"calendar"`
	DateExpr string   `toml:"date_expr"`
	Desc     string   `toml:"desc"`
	Tags     []string `toml:"tags"`
	Person   string   `toml:"-"`
}

func (d *Date) SetTags(tags []string) {
	d.Tags = tags
}

func (d *Date) GetTags() []string {
	return d.Tags
}

var _ tag.Tagged = (*Date)(nil)

func (d *Date) Validate() error {
	if d.DateExpr == "" {
		return errors.New("date expression cannot be empty")
	}

	return nil
}
