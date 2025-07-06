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
	"regexp"
	"strings"
)

var FormatPersonInfo = "NAME [(aka NICK1[, NICK2])] :: description [#tag1, #tag2] [@location1, @location2] [$id:FRIEND_ID]"
var ErrNoInfo = errors.New("no information provided")
var regexPerson *regexp.Regexp

func init() {
	regexPerson = regexp.MustCompile(`(?m)^(?P<name>[^\(\$:\n]+?)\s*(?:\(\s*a\.?k\.?a\.?\s+(?P<nicknames>[^)]*)\))?\s*(?:\$id:(?P<id>[^\s:]+))?\s*(?:::\s*(?P<desc>.+))?$`)
}

func parseNicknames(raw string) []string {
	raw = strings.ReplaceAll(raw, `"`, "")
	raw = strings.ReplaceAll(raw, `'`, "")
	parts := strings.Split(raw, ",")

	var cleaned []string

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}

	return cleaned
}

func ParsePerson(s string) (Person, error) {
	if s == "" {
		return Person{}, ErrNoInfo
	}

	matches := regexPerson.FindStringSubmatch(s)

	if matches == nil {
		return Person{}, ErrNoInfo
	}

	name := strings.TrimSpace(matches[1])
	nicknames := parseNicknames(matches[2])
	desc := strings.TrimSpace(matches[3])

	return Person{
		Name:      name,
		Nicknames: nicknames,
		Desc:      desc,
	}, nil
}
