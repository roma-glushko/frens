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
	"regexp"

	"github.com/roma-glushko/frens/internal/utils"
)

var regex *regexp.Regexp

func init() {
	regex = regexp.MustCompile(`#([a-zA-Z0-9]+(?::[a-zA-Z0-9]+)?(?:-[a-zA-Z0-9]+)?)`)
}

// Parse extracts tags from a string e.g. "#tag1 #tag2" and returns a slice of unique Tag objects.
func Parse(s string) []Tag {
	matches := regex.FindAllString(s, -1)
	tags := make([]Tag, len(matches))

	for i, match := range matches {
		tags[i] = NewTag(match)
	}

	return utils.Unique(tags)
}
