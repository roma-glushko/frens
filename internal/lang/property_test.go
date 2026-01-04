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

package lang

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type props struct {
	ID        string `frentxt:"id"`
	SortBy    string `frentxt:"sort"`
	SortOrder string `frentxt:"order"`
}

func TestExtractProperty(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected props
	}{
		{
			name:  "Extract ID",
			input: "My test entity $id:12345 hello $world",
			expected: props{
				ID: "12345",
			},
		},
		{
			name:  "Extract SortBy and SortOrder",
			input: "Sort by name $sort:name $order:asc",
			expected: props{
				SortBy:    "name",
				SortOrder: "asc",
			},
		},
		{
			name:  "Ignore unknown strings in the props format",
			input: "$id:12345 $name:jim $wife:pam",
			expected: props{
				ID: "12345",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExtractProps[props](tt.input)
			require.NoError(t, err)

			require.Equal(t, tt.expected.ID, result.ID)
		})
	}
}

func TestRenderProps(t *testing.T) {
	tests := []struct {
		name     string
		input    props
		expected string
	}{
		{
			name: "Render ID",
			input: props{
				ID: "12345",
			},
			expected: "$id:12345",
		},
		{
			name: "Render SortBy and SortOrder",
			input: props{
				SortBy:    "name",
				SortOrder: "asc",
			},
			expected: "$sort:name $order:asc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderProps(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}
