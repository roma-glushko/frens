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

	for _, tc := range testcases {
		t.Run(tc.useCase, func(t *testing.T) {
			got, err := ExtractPerson(tc.input)
			require.NoError(t, err)

			require.NotEmpty(t, got)
			require.Equal(t, tc.name, got.Name)
			require.Equal(t, tc.nicknames, got.Nicknames)
			require.Equal(t, tc.tags, got.Tags)
			require.Equal(t, tc.locations, got.Locations)
			require.Equal(t, tc.desc, got.Desc)
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

	for _, tc := range testcases {
		t.Run(tc.title, func(t *testing.T) {
			personInfo := RenderPerson(tc.person)

			require.Equal(t, tc.want, personInfo)
		})
	}
}
