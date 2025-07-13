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
	"testing"

	"github.com/roma-glushko/frens/internal/friend"

	"github.com/stretchr/testify/require"
)

func TestExtractLocMarker(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		input string
		want  []string
	}{
		{"@online", []string{"online"}},
		{"@NYC @scranton", []string{"NYC", "scranton"}},
		{"@ohio@georgia", []string{"ohio", "georgia"}},
		{"@東京", []string{"東京"}},
	}

	for _, tc := range testcases {
		t.Run(tc.input, func(t *testing.T) {
			got := ExtractLocMarkers(tc.input)

			for _, want := range tc.want {
				require.Contains(t, got, want, "Expected to find location %v in %v", want, got)
			}
		})
	}
}

func TestExtractLocation(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name  string
		input string
		want  friend.Location
	}{
		{
			name:  "full info",
			input: "Scranton, USA (a.k.a. The Electric City, Scranton) :: Located a branch of Dunder Mifflin #theoffice",
			want: friend.Location{
				Name:    "Scranton",
				Country: "USA",
				Desc:    "Located a branch of Dunder Mifflin",
				Aliases: []string{"The Electric City", "Scranton"},
				Tags:    []string{"theoffice"},
			},
		},
		{
			name:  "no country",
			input: "New York City (aka NYC, The Big Apple) :: A bustling metropolis known for its skyscrapers and culture",
			want: friend.Location{
				Name:    "New York City",
				Desc:    "A bustling metropolis known for its skyscrapers and culture",
				Aliases: []string{"NYC", "The Big Apple"},
			},
		},
		{
			name:  "no country, no aliases",
			input: "Nashua :: A city in New Hampshire known for its beautiful parks and vibrant community",
			want: friend.Location{
				Name: "Nashua",
				Desc: "A city in New Hampshire known for its beautiful parks and vibrant community",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ExtractLocation(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want.Name, got.Name, "Expected location name to match")
		})
	}
}
