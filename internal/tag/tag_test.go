package tag

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testEntity struct {
	Tags []string
}

var _ Tagged = (*testEntity)(nil)

func (e *testEntity) SetTags(tags []string) {
	e.Tags = tags
}

func (e *testEntity) GetTags() []string {
	return e.Tags
}

func TestTags(t *testing.T) {
	var e testEntity

	require.False(t, HasTag(&e, "corporate"))

	Add(&e, "sales")
	Add(&e, "sales")
	Add(&e, "accounting")

	require.Len(t, e.Tags, 2)

	require.True(t, HasTag(&e, "Sales"))
	require.False(t, HasTag(&e, "warehouse"))

	Remove(&e, "SALES")

	require.Len(t, e.Tags, 1)
}
