// Copyright 2026 Roma Hlushko
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
	"testing"

	"github.com/roma-glushko/frens/internal/friend"
	"github.com/stretchr/testify/require"
)

func TestParseFrenTXT(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name       string
		input      string
		friends    []friend.Person
		locations  []friend.Location
		notes      int
		activities int
	}{
		{
			name: "single friend",
			input: `/f
Alice (aka Al, Ally) :: Close friend from college.
#college, #bestie @Boston $id:ALICE`,
			friends: []friend.Person{
				{
					ID:        "ALICE",
					Name:      "Alice",
					Nicknames: []string{"Al", "Ally"},
					Desc:      "Close friend from college.",
					Tags:      []string{"college", "bestie"},
					Locations: []string{"Boston"},
				},
			},
		},
		{
			name: "single location",
			input: `/l
Paris, France (aka City of Light) :: A city I visited in 2022.
#travel, #art $id:PARIS`,
			locations: []friend.Location{
				{
					ID:      "PARIS",
					Name:    "Paris",
					Country: "France",
					Aliases: []string{"City of Light"},
					Desc:    "A city I visited in 2022.",
					Tags:    []string{"travel", "art"},
				},
			},
		},
		{
			name: "single note",
			input: `/n
2023-08-15 :: Had dinner with Alice in Paris.
#catchup, #dinner @PARIS`,
			notes: 1,
		},
		{
			name: "single activity",
			input: `/act
2023-08-16 :: Jogged with Alice around the Seine.
#exercise, #friends @Paris`,
			activities: 1,
		},
		{
			name: "mixed content",
			input: `/f
Alice (aka Al, Ally) :: Close friend from college.
We met during freshman orientation and have stayed close since.
#college, #bestie @Boston, @Cambridge $id:ALICE

/l
Paris, France (aka City of Light) :: A city I visited in 2022.
Absolutely fell in love with the art and food culture.
#travel, #art $id:PARIS

/n
2023-08-15 :: Had dinner with Alice in Paris. Talked for hours about life and work.
Really missed spending time with her.
#catchup, #dinner @PARIS

/act
2023-08-16 :: Jogged with Alice around the Seine.
Great weather, lots of tourists.
#exercise, #friends @Paris`,
			friends: []friend.Person{
				{
					ID:        "ALICE",
					Name:      "Alice",
					Nicknames: []string{"Al", "Ally"},
					Tags:      []string{"college", "bestie"},
					Locations: []string{"Boston", "Cambridge"},
				},
			},
			locations: []friend.Location{
				{
					ID:      "PARIS",
					Name:    "Paris",
					Country: "France",
					Aliases: []string{"City of Light"},
					Tags:    []string{"travel", "art"},
				},
			},
			notes:      1,
			activities: 1,
		},
		{
			name: "multiline description",
			input: `/f
Bob :: First line of description.
Second line of description.
Third line with #tag @location`,
			friends: []friend.Person{
				{
					Name:      "Bob",
					Tags:      []string{"tag"},
					Locations: []string{"location"},
				},
			},
		},
		{
			name:  "empty input",
			input: "",
		},
		{
			name: "multiple friends",
			input: `/f
Alice :: Friend one
#friend1

/f
Bob :: Friend two
#friend2

/f
Charlie :: Friend three
#friend3`,
			friends: []friend.Person{
				{Name: "Alice", Desc: "Friend one", Tags: []string{"friend1"}},
				{Name: "Bob", Desc: "Friend two", Tags: []string{"friend2"}},
				{Name: "Charlie", Desc: "Friend three", Tags: []string{"friend3"}},
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseFrenTXT(tt.input)
			require.NoError(t, err)
			require.NotNil(t, result)

			if tt.friends != nil {
				require.Len(t, result.Friends, len(tt.friends))

				for i, expected := range tt.friends {
					actual := result.Friends[i]
					require.Equal(t, expected.Name, actual.Name)

					if expected.ID != "" {
						require.Equal(t, expected.ID, actual.ID)
					}

					if expected.Nicknames != nil {
						require.ElementsMatch(t, expected.Nicknames, actual.Nicknames)
					}

					if expected.Tags != nil {
						require.ElementsMatch(t, expected.Tags, actual.Tags)
					}

					if expected.Locations != nil {
						require.ElementsMatch(t, expected.Locations, actual.Locations)
					}
				}
			}

			if tt.locations != nil {
				require.Len(t, result.Locations, len(tt.locations))

				for i, expected := range tt.locations {
					actual := result.Locations[i]
					require.Equal(t, expected.Name, actual.Name)

					if expected.ID != "" {
						require.Equal(t, expected.ID, actual.ID)
					}

					if expected.Country != "" {
						require.Equal(t, expected.Country, actual.Country)
					}

					if expected.Aliases != nil {
						require.ElementsMatch(t, expected.Aliases, actual.Aliases)
					}

					if expected.Tags != nil {
						require.ElementsMatch(t, expected.Tags, actual.Tags)
					}
				}
			}

			require.Len(t, result.Notes, tt.notes)
			require.Len(t, result.Activities, tt.activities)
		})
	}
}

func TestParseFrenTXT_NoteContent(t *testing.T) {
	t.Parallel()

	input := `/n
2023-08-15 :: Had dinner with Alice.
#catchup @PARIS`

	result, err := ParseFrenTXT(input)
	require.NoError(t, err)
	require.Len(t, result.Notes, 1)

	note := result.Notes[0]
	require.Equal(t, friend.EventTypeNote, note.Type)
	require.Contains(t, note.Desc, "Had dinner with Alice")
	require.ElementsMatch(t, []string{"catchup"}, note.Tags)
	require.ElementsMatch(t, []string{"PARIS"}, note.LocationIDs)
}

func TestParseFrenTXT_ActivityContent(t *testing.T) {
	t.Parallel()

	input := `/act
2023-08-16 :: Jogged around the Seine.
#exercise @Paris`

	result, err := ParseFrenTXT(input)
	require.NoError(t, err)
	require.Len(t, result.Activities, 1)

	activity := result.Activities[0]
	require.Equal(t, friend.EventTypeActivity, activity.Type)
	require.Contains(t, activity.Desc, "Jogged around the Seine")
	require.ElementsMatch(t, []string{"exercise"}, activity.Tags)
	require.ElementsMatch(t, []string{"Paris"}, activity.LocationIDs)
}

func TestRenderFrenTXT(t *testing.T) {
	t.Parallel()

	data := &FrenTXTImport{
		Friends: []friend.Person{
			{
				ID:        "alice",
				Name:      "Alice",
				Nicknames: []string{"Al"},
				Desc:      "My best friend",
				Tags:      []string{"bestie"},
				Locations: []string{"Boston"},
			},
		},
		Locations: []friend.Location{
			{
				ID:      "paris",
				Name:    "Paris",
				Country: "France",
				Desc:    "City of Light",
				Tags:    []string{"travel"},
			},
		},
	}

	output := RenderFrenTXT(data)

	require.Contains(t, output, "/f")
	require.Contains(t, output, "Alice")
	require.Contains(t, output, "Al")
	require.Contains(t, output, "#bestie")
	require.Contains(t, output, "@Boston")
	require.Contains(t, output, "$id:alice")

	require.Contains(t, output, "/l")
	require.Contains(t, output, "Paris")
	require.Contains(t, output, "France")
	require.Contains(t, output, "#travel")
	require.Contains(t, output, "$id:paris")
}

func TestRenderFrenTXT_RoundTrip(t *testing.T) {
	t.Parallel()

	original := &FrenTXTImport{
		Friends: []friend.Person{
			{
				Name:      "Bob",
				Nicknames: []string{"Bobby"},
				Desc:      "College friend",
				Tags:      []string{"college"},
			},
		},
		Locations: []friend.Location{
			{
				Name:    "New York",
				Country: "USA",
				Desc:    "The Big Apple",
			},
		},
	}

	// Render to FrenTXT
	output := RenderFrenTXT(original)

	// Parse back
	parsed, err := ParseFrenTXT(output)
	require.NoError(t, err)

	// Verify round-trip
	require.Len(t, parsed.Friends, 1)
	require.Equal(t, "Bob", parsed.Friends[0].Name)
	require.ElementsMatch(t, []string{"Bobby"}, parsed.Friends[0].Nicknames)
	require.ElementsMatch(t, []string{"college"}, parsed.Friends[0].Tags)

	require.Len(t, parsed.Locations, 1)
	require.Equal(t, "New York", parsed.Locations[0].Name)
	require.Equal(t, "USA", parsed.Locations[0].Country)
}
