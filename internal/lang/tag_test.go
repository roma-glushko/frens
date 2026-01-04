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

	"github.com/roma-glushko/frens/internal/tag"

	"github.com/stretchr/testify/require"
)

func TestExtractTags(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		useCase string
		input   string
		want    []tag.Tag
	}{
		{
			useCase: "space-separated tags",
			input:   "#tag1 #tag2",
			want:    []tag.Tag{{Name: "tag1"}, {Name: "tag2"}},
		},
		{
			useCase: "joined tags",
			input:   "#tag3#tag4",
			want:    []tag.Tag{{Name: "tag3"}, {Name: "tag4"}},
		},
		{
			useCase: "semi-colon subgrouped tags",
			input:   "#school:biology #school:math",
			want:    []tag.Tag{{Name: "school:biology"}, {Name: "school:math"}},
		},
		{
			useCase: "cyrillic tags",
			input:   "#укрліт:поезія #школа #сімя-батьки",
			want:    []tag.Tag{{Name: "укрліт:поезія"}, {Name: "школа"}, {Name: "сімя-батьки"}},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.input, func(t *testing.T) {
			got := ExtractTags(tc.input)

			for _, want := range tc.want {
				require.Contains(t, got, want, "Expected to find tag %v in %v", want, got)
			}
		})
	}
}
