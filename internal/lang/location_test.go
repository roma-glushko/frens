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
		{
			name:  "cyrillic location",
			input: "Київ, Україна (aka Kyiv)::Столиця України $id:kyiv",
			want: friend.Location{
				ID:      "kyiv",
				Name:    "Київ",
				Country: "Україна",
				Aliases: []string{"Kyiv"},
				Desc:    "Столиця України",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ExtractLocation(tc.input)
			require.NoError(t, err)

			require.Equal(t, tc.want.Name, got.Name)
			require.Equal(t, tc.want.Country, got.Country)
			require.ElementsMatch(t, tc.want.Aliases, got.Aliases)
			require.Equal(t, tc.want.Desc, got.Desc)
			require.ElementsMatch(t, tc.want.Tags, got.Tags)
		})
	}
}

func TestRenderLocation(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		title string
		loc   friend.Location
		want  string
	}{
		{
			title: "Location with all information",
			loc: friend.Location{
				Name:    "Scranton",
				Country: "USA",
				Aliases: []string{"The Electric City"},
				Desc:    "Branch of Dunder Mifflin",
				Tags:    []string{"office", "dunderm"},
			},
			want: "Scranton, USA (a.k.a. The Electric City) :: Branch of Dunder Mifflin #dunderm #office",
		},
		{
			title: "Location with ID",
			loc: friend.Location{
				ID:      "nyc",
				Name:    "New York City",
				Country: "USA",
				Aliases: []string{"NYC", "The City of Love"},
				Desc:    "HQ of Dunder Mifflin",
				Tags:    []string{"office", "dunderm"},
			},
			want: "New York City, USA (a.k.a. NYC, The City of Love) :: HQ of Dunder Mifflin #dunderm #office $id:nyc",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.title, func(t *testing.T) {
			locInfo := RenderLocation(tt.loc)

			require.Equal(t, tt.want, locInfo)
		})
	}
}

func TestExtractLocationQuery(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		title string
		input string
		query friend.ListLocationQuery
	}{
		{
			title: "keyword search",
			input: "electric",
			query: friend.ListLocationQuery{
				Keyword: "electric",
			},
		},
		{
			title: "tags only",
			input: "#office #dunderm",
			query: friend.ListLocationQuery{
				Tags: []string{"office", "dunderm"},
			},
		},
		{
			title: "sort & order",
			input: "$sort:alpha $order:reverse",
			query: friend.ListLocationQuery{
				SortBy:    friend.SortAlpha,
				SortOrder: friend.SortOrderReverse,
			},
		},
		{
			title: "all query information",
			input: "new #corporate $sort:recency $order:direct",
			query: friend.ListLocationQuery{
				Keyword:   "new",
				Tags:      []string{"corporate"},
				SortBy:    friend.SortRecency,
				SortOrder: friend.SortOrderDirect,
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.title, func(t *testing.T) {
			q, err := ExtractLocationQuery(tt.input)
			require.NoError(t, err)

			require.Equal(t, tt.query.Keyword, q.Keyword)
			require.ElementsMatch(t, tt.query.Tags, q.Tags)
			require.Equal(t, tt.query.SortBy, q.SortBy)
			require.Equal(t, tt.query.SortOrder, q.SortOrder)
		})
	}
}
