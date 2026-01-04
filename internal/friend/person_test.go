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

package friend

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPerson_Validation(t *testing.T) {
	var noName Person

	require.Error(t, noName.Validate())

	withName := Person{Name: "Jim Halpert"}
	require.NoError(t, withName.Validate())
}

func TestPerson_Location(t *testing.T) {
	loc := "Scranton"
	p := Person{Name: "Jim Halpert"}

	p.AddLocation(loc)

	require.True(t, p.HasLocations([]string{loc}))

	p.AddLocation(loc)

	require.Len(t, p.Locations, 1)
}

func TestPerson_Nickname(t *testing.T) {
	nickname := "Big Tuna"
	p := Person{Name: "Jim Halpert"}

	p.AddNickname(nickname)
	p.AddNickname(nickname)

	require.Len(t, p.Nicknames, 1)

	p.RemoveNickname(nickname)
	p.RemoveNickname(nickname)

	require.Empty(t, p.Nicknames)
}

func TestPerson_String(t *testing.T) {
	name := "Jim Halpert"
	nick1 := "Big Tuna"
	nick2 := "Jimbo"

	p := Person{Name: name}

	require.Equal(t, name, p.String())

	p.AddNickname(nick1)
	p.AddNickname(nick2)

	require.Contains(t, p.String(), name)
	require.Contains(t, p.String(), nick1)
	require.Contains(t, p.String(), nick2)
}
