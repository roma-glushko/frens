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
	"regexp"
	"strings"

	"github.com/roma-glushko/frens/internal/tag"

	"github.com/roma-glushko/frens/internal/utils"
)

var tagRe *regexp.Regexp

func init() {
	tagRe = regexp.MustCompile(`#([\p{L}\p{N}]+(?::[\p{L}\p{N}]+)?(?:-[\p{L}\p{N}]+)?)`)
}

// ExtractTags extracts tags from a string e.g. "#tag1 #tag2" and returns a slice of unique Tag objects.
func ExtractTags(s string) []tag.Tag {
	matches := tagRe.FindAllString(s, -1)
	tags := make([]tag.Tag, len(matches))

	for i, match := range matches {
		tags[i] = tag.NewTag(match)
	}

	return utils.Unique(tags)
}

func RemoveTags(s string) string {
	return tagRe.ReplaceAllString(s, "")
}

func RenderTags(ts []string) string {
	tags := make([]tag.Tag, 0, len(ts))

	for _, t := range ts {
		t = strings.TrimSpace(t)

		if t == "" {
			continue
		}

		tags = append(tags, tag.NewTag(t))
	}

	return tag.Tags(tags).String()
}
