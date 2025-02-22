package location

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLocation_Validation(t *testing.T) {
	var l Location

	require.Error(t, l.Validate())

	l.Name = "Los Angeles"

	require.NoError(t, l.Validate())
}

func TestLocation_Match(t *testing.T) {
	var l Location

	l.Name = "Los Angeles"
	l.AddAlias("LA")

	require.True(t, l.Match("los"))
	require.True(t, l.Match("angeles"))
	require.True(t, l.Match("los angeles"))
	require.False(t, l.Match("new"))
	require.True(t, l.Match("la"))
}

func TestLocation_Alias(t *testing.T) {
	var l Location

	l.AddAlias("LA")
	l.AddAlias("Los Angeles")

	require.Len(t, l.Alias, 2)

	l.RemoveAlias("la")
	l.RemoveAlias("LA")

	require.Len(t, l.Alias, 1)
}
