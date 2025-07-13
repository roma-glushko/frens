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

package matcher

import (
	"maps"
	"slices"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Matchable interface {
	Refs() []string
}

type Pattern[T Matchable] struct {
	Ref      string
	Entities []T
}

type Matcher[T Matchable] struct {
	EntityPatterns map[string]Pattern[T]
}

type Match[T any] struct {
	Entities   []T
	MatchedRef string
}

func NewMatcher[T Matchable]() *Matcher[T] {
	return &Matcher[T]{
		EntityPatterns: make(map[string]Pattern[T]),
	}
}

func (m *Matcher[T]) Add(entity T) {
	for _, ref := range entity.Refs() {
		if pattern, exists := m.EntityPatterns[ref]; exists {
			pattern.Entities = append(pattern.Entities, entity)

			m.EntityPatterns[ref] = pattern

			continue
		}

		m.EntityPatterns[ref] = Pattern[T]{Ref: ref, Entities: []T{entity}}
	}
}

func (m *Matcher[T]) Match(input string) []Match[T] { //nolint:cyclop
	searchKeys := slices.Collect(maps.Keys(m.EntityPatterns))

	sort.SliceStable(searchKeys, func(i, j int) bool {
		return len(searchKeys[i]) > len(searchKeys[j])
	})

	var found []Match[T]

	lowerInput := strings.ToLower(input)

	for _, searchKey := range searchKeys {
		idx := 0
		pattern := m.EntityPatterns[searchKey]
		lowerKey := strings.ToLower(pattern.Ref)

		for {
			start := strings.Index(lowerInput[idx:], lowerKey)

			if start == -1 {
				break
			}

			start += idx
			end := start + len(lowerKey)

			// Check forbidden leading characters
			if start > 0 {
				r, _ := utf8.DecodeLastRuneInString(input[:start])
				if r == '\\' || r == '_' || r == '*' || unicode.IsLetter(r) {
					idx = start + 1
					continue
				}
			}

			// Check forbidden trailing characters
			if end < len(input) {
				r, _ := utf8.DecodeRuneInString(input[end:])

				if r == '_' || r == '*' || unicode.IsLetter(r) {
					idx = start + 1
					continue
				}
			}

			found = append(found, Match[T]{
				Entities:   pattern.Entities,
				MatchedRef: pattern.Ref,
			})

			idx = end
		}
	}

	return found
}
