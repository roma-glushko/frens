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

func TestLocation_Validation(t *testing.T) {
	var l Location

	require.Error(t, l.Validate())

	l.Name = "Los Angeles"

	require.NoError(t, l.Validate())
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
