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

	require.Len(t, f.Nicknames, 0)
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
