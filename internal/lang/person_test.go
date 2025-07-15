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
		input     string
		Name      string
		Nicknames []string
		Locations []string
		Desc      string
		Tags      []string
		ID        string
	}{
		{
			"Michael Harry Scott (a.k.a. The World's Best Boss, Mike) :: my Dunder Mifflin boss #office @Scranton",
			"Michael Harry Scott",
			[]string{"The World's Best Boss", "Mike"},
			[]string{"Scranton"},
			"my Dunder Mifflin boss",
			[]string{"office"},
			"mscott",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ExtractPerson(tc.input)
			require.NoError(t, err)

			require.NotEmpty(t, got)
			require.Equal(t, tc.Name, got.Name)
			require.Equal(t, tc.Nicknames, got.Nicknames)
			require.Equal(t, tc.Tags, got.Tags)
			require.Equal(t, tc.Locations, got.Locations)
			require.Equal(t, tc.Desc, got.Desc)
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
