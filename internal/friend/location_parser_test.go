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

func TestLocationParser(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		input string
		want  []string
	}{
		{"@online", []string{"online"}},
		{"@NYC @scranton", []string{"NYC", "scranton"}},
		{"@ohio@georgia", []string{"ohio", "georgia"}},
		{"@東京", []string{"東京"}},
	}

	for _, tc := range testcases {
		t.Run(tc.input, func(t *testing.T) {
			got := ParseLocMarkers(tc.input)

			for _, want := range tc.want {
				require.Contains(t, got, want, "Expected to find location %v in %v", want, got)
			}
		})
	}
}
