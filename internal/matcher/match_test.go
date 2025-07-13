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
	"testing"

	"github.com/stretchr/testify/require"
)

type testEntity struct {
	ID         int
	References []string
}

func (e testEntity) Refs() []string {
	return e.References
}

func TestMatcher(t *testing.T) {
	t.Parallel()

	matcher := NewMatcher[testEntity]()

	matcher.Add(testEntity{ID: 1, References: []string{"Philly", "Philadelphia"}})
	matcher.Add(testEntity{ID: 2, References: []string{"Scranton", "Electric City"}})
	matcher.Add(testEntity{ID: 3, References: []string{"New York", "NY", "NYC"}})

	tests := []struct {
		input    string
		matches  int
		wantRefs []string
	}{
		{"It's a big ride from Philly to Scranton", 2, []string{"Philly", "Scranton"}},
		{"NY, the city of love", 1, []string{"NY"}},
		{"Nychthemeron is a full period of a night and a day", 0, []string{}},
		{"Kevin has his own band - Scrantonicity 2", 0, []string{}},
		{"_New York_ pizza is *the best*", 0, []string{}},
		{"\\Scranton, Pennsylvania", 0, []string{}},
		{"nyc nY neWYork", 2, []string{"NYC", "NY"}},
		{"Phil drives corvet", 0, []string{}},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			matches := matcher.Match(test.input)

			require.Len(t, matches, test.matches)

			foundRefs := make([]string, 0, len(matches))

			for _, match := range matches {
				foundRefs = append(foundRefs, match.MatchedRef)
			}

			for _, ref := range test.wantRefs {
				require.Contains(t, foundRefs, ref, "Expected to find '%s' in matches", ref)
			}
		})
	}
}
