package utils

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
