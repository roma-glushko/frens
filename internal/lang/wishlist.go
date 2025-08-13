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

package lang

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/tag"
)

var FormatWishlistItem = fmt.Sprintf(
	"[DESCRIPTION or URL] [%s] [%s] $price:PRICE",
	FormatTags,
	FormatLocationMarkers,
)

var urlRegex *regexp.Regexp

func init() {
	urlRegex = regexp.MustCompile(`https?://[^\s]+`)
}

type itemProps struct {
	Price string `frentxt:"price"`
}

func ExtractURLs(s string) []string {
	return urlRegex.FindAllString(s, -1)
}

func RemoveURLs(s string) string {
	return urlRegex.ReplaceAllString(s, "")
}

func ExtractWishlistItem(s string) (friend.WishlistItem, error) {
	if s == "" {
		return friend.WishlistItem{}, ErrNoInfo
	}

	props, err := ExtractProps[itemProps](s)
	if err != nil {
		return friend.WishlistItem{}, fmt.Errorf(
			"failed to parse wishlist item properties: %w",
			err,
		)
	}

	tags := tag.Tags(ExtractTags(s)).ToNames()
	urls := ExtractURLs(s)

	s = RemoveTags(s)
	s = RemoveProps(s)
	s = RemoveURLs(s)

	desc := strings.TrimSpace(s)

	var url string

	if len(urls) == 1 {
		url = urls[0]
	}

	if len(urls) > 1 {
		return friend.WishlistItem{}, fmt.Errorf(
			"wishlist item cannot contain multiple URLs: %s",
			strings.Join(urls, ", "),
		)
	}

	return friend.WishlistItem{
		CreatedAt: time.Now(),
		Price:     props.Price,
		Desc:      desc,
		Link:      url,
		Tags:      tags,
	}, nil
}

func RenderWishlistItem(item friend.WishlistItem) string {
	var sb strings.Builder

	if item.Desc != "" {
		sb.WriteString(item.Desc)
	}

	if item.Link != "" {
		if sb.Len() > 0 {
			sb.WriteString(" ")
		}

		sb.WriteString(item.Link)
	}

	if item.Price != "" {
		props := itemProps{Price: item.Price}
		sb.WriteString(RenderProps[itemProps](props))
	}

	if len(item.Tags) > 0 {
		sb.WriteString(RenderTags(item.Tags))
	}

	return sb.String()
}
