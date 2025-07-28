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

func TestPersonParser(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		useCase   string
		input     string
		id        string
		name      string
		nicknames []string
		locations []string
		desc      string
		tags      []string
	}{
		{
			useCase:   "full person info",
			input:     "Michael Harry Scott (a.k.a. The World's Best Boss, Mike) :: my Dunder Mifflin boss #office @Scranton $id:mscott",
			id:        "mscott",
			name:      "Michael Harry Scott",
			nicknames: []string{"The World's Best Boss", "Mike"},
			locations: []string{"Scranton"},
			desc:      "my Dunder Mifflin boss",
			tags:      []string{"office"},
		},
		{
			useCase:   "cyrillic person info",
			input:     "Тарас Шевченко (a.k.a. Тарас Григорович) :: український поет #укрліт @kyiv $id:shevchenko",
			id:        "shevchenko",
			name:      "Тарас Шевченко",
			nicknames: []string{"Тарас Григорович"},
			locations: []string{"kyiv"},
			desc:      "український поет",
			tags:      []string{"укрліт"},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.useCase, func(t *testing.T) {
			got, err := ExtractPerson(tt.input)
			require.NoError(t, err)

			require.NotEmpty(t, got)
			require.Equal(t, tt.name, got.Name)
			require.Equal(t, tt.nicknames, got.Nicknames)
			require.Equal(t, tt.tags, got.Tags)
			require.Equal(t, tt.locations, got.Locations)
			require.Equal(t, tt.desc, got.Desc)
		})
	}
}

func TestPersonFormatter(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		title  string
		person friend.Person
		want   string
	}{
		{
			title: "Person with all information",
			person: friend.Person{
				Name:      "Michael Harry Scott",
				Nicknames: []string{"The World's Best Boss", "Mike"},
				Desc:      "my Dunder Mifflin boss",
				Locations: []string{"Scranton"},
				Tags:      []string{"office"},
			},
			want: "Michael Harry Scott (a.k.a. The World's Best Boss, Mike) :: my Dunder Mifflin boss @Scranton #office",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.title, func(t *testing.T) {
			personInfo := RenderPerson(tt.person)

			require.Equal(t, tt.want, personInfo)
		})
	}
}

func TestExtractPersonQuery(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		title string
		input string
		query friend.ListFriendQuery
	}{
		{
			title: "keyword search",
			input: "michael",
			query: friend.ListFriendQuery{
				Keyword: "michael",
			},
		},
		{
			title: "keyword search & locations",
			input: "michael @scranton @nyc",
			query: friend.ListFriendQuery{
				Keyword:   "michael",
				Locations: []string{"scranton", "nyc"},
			},
		},
		{
			title: "tags only",
			input: "#office #dunderm",
			query: friend.ListFriendQuery{
				Tags: []string{"office", "dunderm"},
			},
		},
		{
			title: "locations only",
			input: "@scranton @utica",
			query: friend.ListFriendQuery{
				Locations: []string{"scranton", "utica"},
			},
		},
		{
			title: "sort",
			input: "$sort:alpha $order:reverse",
			query: friend.ListFriendQuery{
				SortBy:    friend.SortAlpha,
				SortOrder: friend.SortOrderReverse,
			},
		},
		{
			title: "all query information",
			input: "pam #art @nyc $sort:recency $order:direct",
			query: friend.ListFriendQuery{
				Keyword:   "pam",
				Locations: []string{"nyc"},
				Tags:      []string{"art"},
				SortBy:    friend.SortRecency,
				SortOrder: friend.SortOrderDirect,
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.title, func(t *testing.T) {
			q, err := ExtractPersonQuery(tt.input)
			require.NoError(t, err)

			require.Equal(t, tt.query.Keyword, q.Keyword)
			require.ElementsMatch(t, tt.query.Locations, q.Locations)
			require.ElementsMatch(t, tt.query.Tags, q.Tags)
			require.Equal(t, tt.query.SortBy, q.SortBy)
			require.Equal(t, tt.query.SortOrder, q.SortOrder)
		})
	}
}
