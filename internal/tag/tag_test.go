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

	require.False(t, HasTags(&e, []string{"corporate"}))

	AddStr(&e, []string{"sales"})
	AddStr(&e, []string{"sales"})
	AddStr(&e, []string{"accounting"})

	require.Len(t, e.Tags, 2)

	require.True(t, HasTags(&e, []string{"Sales"}))
	require.False(t, HasTags(&e, []string{"warehouse"}))

	Remove(&e, "SALES")

	require.Len(t, e.Tags, 1)
}
