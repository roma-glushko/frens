package utils

import (
	"maps"
	"slices"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Pattern[T any] struct {
	Ref      string
	Entities []T
}

type Matcher[T any] struct {
	EntityPatterns map[string]Pattern[T]
}

type Match[T any] struct {
	Entities   []T
	MatchedRef string
}

func NewMatcher[T any]() *Matcher[T] {
	return &Matcher[T]{
		EntityPatterns: make(map[string]Pattern[T]),
	}
}

func (m *Matcher[T]) Add(entity T, refs []string) {
	for _, ref := range refs {
		if pattern, exists := m.EntityPatterns[ref]; exists {
			pattern.Entities = append(pattern.Entities, entity)

			m.EntityPatterns[ref] = pattern

			continue
		}

		m.EntityPatterns[ref] = Pattern[T]{Ref: ref, Entities: []T{entity}}
	}
}

func (m *Matcher[T]) Match(input string) []Match[T] {
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
