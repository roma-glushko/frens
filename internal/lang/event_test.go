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

func TestExtractDate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		dateStr string
		date    time.Time
	}{
		{
			dateStr: "a min ago",
			date:    time.Now().Add(-1 * time.Minute),
		},
		{
			dateStr: "yesterday",
			date:    time.Now().Add(-24 * time.Hour),
		},
		{
			dateStr: "2 days ago",
			date:    time.Now().Add(-2 * 24 * time.Hour),
		},
		{
			dateStr: "3 weeks ago",
			date:    time.Now().Add(-3 * 7 * 24 * time.Hour),
		},
		{
			dateStr: "a year ago",
			date:    time.Now().Add(-1 * 365 * 24 * time.Hour),
		},
		{
			dateStr: "Apr 1st",
			date:    time.Date(time.Now().Year(), time.April, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			dateStr: "9/11",
			date:    time.Date(time.Now().Year(), time.September, 11, 0, 0, 0, 0, time.UTC),
		},
		{
			dateStr: "1967-07-30",
			date:    time.Date(1967, time.July, 30, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, test := range tests {
		t.Run(test.dateStr, func(t *testing.T) {
			gotDate := ExtractDate(test.dateStr)

			require.WithinDuration(t, test.date, gotDate, 1*time.Second)
		})
	}
}

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
