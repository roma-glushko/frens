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

func TestExtractURLs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single URL",
			input:    "Check this out: https://example.com",
			expected: []string{"https://example.com"},
		},
		{
			name:     "multiple URLs",
			input:    "Visit http://example.com and https://another-example.com for more info.",
			expected: []string{"http://example.com", "https://another-example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractURLs(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractWishlistItem(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected friend.WishlistItem
	}{
		{
			name:  "basic item",
			input: "https://example.com/keychron-a90 $price:100USD #tag1#tag2",
			expected: friend.WishlistItem{
				Link:  "https://example.com/keychron-a90",
				Price: "100USD",
				Tags:  []string{"tag1", "tag2"},
			},
		},
		{
			name:  "item with description",
			input: "Keychron A90 mechanical keyboard https://example.com/keychron-a90 $price:100USD #tech #keyboard",
			expected: friend.WishlistItem{
				Desc:  "Keychron A90 mechanical keyboard",
				Link:  "https://example.com/keychron-a90",
				Price: "100USD",
				Tags:  []string{"tech", "keyboard"},
			},
		},
		{
			name:  "item with desc only",
			input: "Keychron A90 mechanical keyboard",
			expected: friend.WishlistItem{
				Desc: "Keychron A90 mechanical keyboard",
			},
		},
		{
			name:  "URL only",
			input: "https://example.com/keychron-a90",
			expected: friend.WishlistItem{
				Link: "https://example.com/keychron-a90",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExtractWishlistItem(tt.input)
			require.NoError(t, err)

			require.Equal(t, tt.expected.Link, result.Link)
			require.Equal(t, tt.expected.Desc, result.Desc)
			require.Equal(t, tt.expected.Price, result.Price)
			require.ElementsMatch(t, tt.expected.Tags, result.Tags)
		})
	}
}
