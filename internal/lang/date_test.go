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
	"time"

	"github.com/roma-glushko/frens/internal/friend"

	"github.com/stretchr/testify/require"
)

func TestExtractDate(t *testing.T) {
	t.Parallel()

	n := time.Now().UTC()

	tests := []struct {
		dateStr string
		date    time.Time
	}{
		{
			dateStr: "a min ago",
			date:    n.Add(-1 * time.Minute),
		},
		{
			dateStr: "yesterday",
			date:    n.Add(-24 * time.Hour),
		},
		{
			dateStr: "2 days ago",
			date:    n.Add(-2 * 24 * time.Hour),
		},
		{
			dateStr: "3 weeks ago",
			date:    n.Add(-3 * 7 * 24 * time.Hour),
		},
		{
			dateStr: "a year ago",
			date:    n.Add(-1 * 365 * 24 * time.Hour),
		},
		{
			dateStr: "Apr 1st",
			date:    time.Date(n.Year(), time.April, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			dateStr: "9/11",
			date:    time.Date(n.Year(), time.September, 11, 0, 0, 0, 0, time.UTC),
		},
		{
			dateStr: "1967-07-30",
			date:    time.Date(1967, time.July, 30, 0, 0, 0, 0, time.UTC),
		},
		{
			dateStr: "March 21st",
			date:    time.Date(n.Year(), time.March, 21, 0, 0, 0, 0, time.UTC),
		},
		{
			dateStr: "tomorrow",
			date:    n.Add(24 * time.Hour),
		},
		{
			dateStr: "today 5pm",
			date: time.Date(
				n.Year(),
				n.Month(),
				n.Day(),
				17,
				0,
				0,
				0,
				time.UTC,
			),
		},
		{
			dateStr: "in 3 days",
			date:    time.Now().Add(3 * 24 * time.Hour),
		},
		{
			dateStr: "1991",
			date:    time.Date(1991, n.Month(), n.Day(), 0, 0, 0, 0, time.UTC),
		},
	}

	for _, test := range tests {
		t.Run(test.dateStr, func(t *testing.T) {
			gotDate := ExtractDate(test.dateStr)

			require.WithinDuration(t, test.date, gotDate, 1*time.Second)
		})
	}
}

func TestExtractDateInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected friend.Date
	}{
		{
			name:  "date & desc",
			input: "May 13th :: birthday",
			expected: friend.Date{
				DateExpr: "May 13th",
				Desc:     "birthday",
				Calendar: friend.CalendarGregorian,
			},
		},
		{
			name:  "just date",
			input: "Jul 30 1996",
			expected: friend.Date{
				DateExpr: "Jul 30 1996",
				Calendar: friend.CalendarGregorian,
			},
		},
		{
			name:  "date & desc & cal",
			input: "Av 16 5784 :: birthday   $cal:hebrew",
			expected: friend.Date{
				Calendar: friend.CalendarHebrew,
				DateExpr: "Av 16 5784",
				Desc:     "birthday",
			},
		},
		{
			name:  "date & desc & tags",
			input: "Jul 30 :: birthday #bday",
			expected: friend.Date{
				Calendar: friend.CalendarGregorian,
				DateExpr: "Jul 30",
				Desc:     "birthday",
				Tags:     []string{"bday"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ExtractDateInfo(tt.input)
			require.NoError(t, err)

			require.Equal(t, tt.expected.Calendar, result.Calendar)
			require.Equal(t, tt.expected.DateExpr, result.DateExpr)
			require.Equal(t, tt.expected.Desc, result.Desc)
			require.ElementsMatch(t, tt.expected.Tags, result.Tags)
		})
	}
}

func TestRenderDateInfo(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		date     friend.Date
		expected string
	}{
		{
			name: "date with description",
			date: friend.Date{
				DateExpr: "May 13th",
				Desc:     "birthday",
				Calendar: friend.CalendarGregorian,
			},
			expected: "May 13th :: birthday",
		},
		{
			name: "date without description",
			date: friend.Date{
				DateExpr: "Jul 30 1996",
				Calendar: friend.CalendarGregorian,
			},
			expected: "Jul 30 1996",
		},
		{
			name: "date with description and tags",
			date: friend.Date{
				DateExpr: "Jul 30",
				Desc:     "birthday",
				Tags:     []string{"bday"},
			},
			expected: "Jul 30 :: birthday #bday",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderDateInfo(tt.date)
			require.Equal(t, tt.expected, result)
		})
	}
}
