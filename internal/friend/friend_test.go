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

package friend

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFriend_Validation(t *testing.T) {
	var noName Friend

	require.Error(t, noName.Validate())

	withName := Friend{Name: "Jim Halpert"}
	require.NoError(t, withName.Validate())
}

func TestFriend_Location(t *testing.T) {
	loc := "Scranton"
	f := Friend{Name: "Jim Halpert"}

	f.AddLocation(loc)

	require.True(t, f.HasLocation(loc))

	f.AddLocation(loc)

	require.Len(t, f.Locations, 1)
}

func TestFriend_Nickname(t *testing.T) {
	nickname := "Big Tuna"
	f := Friend{Name: "Jim Halpert"}

	f.AddNickname(nickname)
	f.AddNickname(nickname)

	require.Len(t, f.Nicknames, 1)

	f.RemoveNickname(nickname)
	f.RemoveNickname(nickname)

	require.Empty(t, f.Nicknames)
}

func TestFriend_String(t *testing.T) {
	name := "Jim Halpert"
	nick1 := "Big Tuna"
	nick2 := "Jimbo"

	f := Friend{Name: name}

	require.Equal(t, name, f.String())

	f.AddNickname(nick1)
	f.AddNickname(nick2)

	require.Contains(t, f.String(), name)
	require.Contains(t, f.String(), nick1)
	require.Contains(t, f.String(), nick2)
}
