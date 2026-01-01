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
	"strings"
	"time"
)

type WishlistItem struct {
	ID        string    `toml:"id"`
	CreatedAt time.Time `toml:"created_at"`
	Desc      string    `toml:"desc,omitempty"`
	Link      string    `toml:"link,omitempty"`
	Price     string    `toml:"price,omitempty"`
	Tags      []string  `toml:"tags,omitempty"`
	Location  []string  `toml:"locations,omitempty"`
	Person    string    `toml:"-"`
}

func (w *WishlistItem) SetTags(tags []string) {
	w.Tags = tags
}

func (w *WishlistItem) GetTags() []string {
	return w.Tags
}

func (i *WishlistItem) Validate() error {
	if i.Desc == "" && i.Link == "" {
		return errors.New("wishlist item must have either a description or a link")
	}

	if i.Link != "" && !strings.HasPrefix(i.Link, "http") {
		return errors.New("wishlist item link must start with http or https")
	}

	return nil
}
