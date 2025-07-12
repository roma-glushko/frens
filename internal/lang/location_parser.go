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

	"github.com/roma-glushko/frens/internal/utils"
)

var regex *regexp.Regexp

func init() {
	regex = regexp.MustCompile(`@([\p{L}0-9_-]+)`)
}

func ParseLocMarkers(s string) []string {
	matches := regex.FindAllString(s, -1)
	locationIDs := make([]string, len(matches))

	for i, match := range matches {
		locationIDs[i] = strings.TrimLeft(match, "@")
	}

	return utils.Unique(locationIDs)
}

func RemoveLocMarkers(s string) string {
	return regex.ReplaceAllString(s, "")
}

func RenderLocMarkers(locations []string) string {
	if len(locations) == 0 {
		return ""
	}

	markers := make([]string, len(locations))

	for i, loc := range locations {
		loc = strings.TrimSpace(loc)

		if loc == "" {
			continue
		}

		markers[i] = "@" + loc
	}

	return strings.Join(markers, " ")
}
