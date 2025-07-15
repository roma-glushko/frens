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
	"fmt"
	"testing"
	"time"

	"github.com/roma-glushko/frens/internal/friend"

	"github.com/stretchr/testify/require"
)

func TestExtractActivity_EmptyDate(t *testing.T) {
	t.Parallel()

	wantDesc := "Angela has a little secret with Dwight"

	tests := []struct {
		desc string
	}{
		{desc: wantDesc},
		{desc: ":: " + wantDesc},
		{desc: " :: " + wantDesc},
		{desc: " :: " + wantDesc},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			e, err := ExtractEvent(friend.EventTypeActivity, test.desc)
			require.NoError(t, err)

			require.WithinDuration(t, time.Now(), e.Date, 1*time.Second)
			require.Equal(t, friend.EventTypeActivity, e.Type)
			require.Equal(t, wantDesc, e.Desc)
		})
	}
}

func TestExtractActivity_DescTrimmed(t *testing.T) {
	t.Parallel()

	wantDesc := "I've just met Bob Vance, Vance Refrigeration"

	tests := []struct {
		desc string
	}{
		{desc: fmt.Sprintf(" %s ", wantDesc)},
		{desc: fmt.Sprintf(" 	%s		", wantDesc)},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			e, err := ExtractEvent(friend.EventTypeActivity, test.desc)
			require.NoError(t, err)

			require.Equal(t, wantDesc, e.Desc)
		})
	}
}
