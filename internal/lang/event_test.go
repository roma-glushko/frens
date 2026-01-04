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

func TestExtractEventQuery(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		title string
		input string
		query friend.ListEventQuery
	}{
		{
			title: "keyword search",
			input: "electric",
			query: friend.ListEventQuery{
				Keyword: "electric",
			},
		},
		{
			title: "tags only",
			input: "#office #dunderm",
			query: friend.ListEventQuery{
				Tags: []string{"office", "dunderm"},
			},
		},
		{
			title: "locations only",
			input: "@scranton @utica",
			query: friend.ListEventQuery{
				Locations: []string{"scranton", "utica"},
			},
		},
		{
			title: "sort & order",
			input: "$sort:alpha $order:reverse",
			query: friend.ListEventQuery{
				SortBy:    friend.SortAlpha,
				SortOrder: friend.SortOrderReverse,
			},
		},
		{
			title: "date range",
			input: "$since:2023-01-01 $until:2023-12-31",
			query: friend.ListEventQuery{
				Since: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
				Until: time.Date(2023, time.December, 31, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			title: "all query information",
			input: "new #corporate $since:2023-01-01 $until:2023-12-31 $sort:recency $order:direct",
			query: friend.ListEventQuery{
				Keyword:   "new",
				Tags:      []string{"corporate"},
				Since:     time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
				Until:     time.Date(2023, time.December, 31, 0, 0, 0, 0, time.UTC),
				SortBy:    friend.SortRecency,
				SortOrder: friend.SortOrderDirect,
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.title, func(t *testing.T) {
			q, err := ExtractEventQuery(tt.input)
			require.NoError(t, err)

			require.Equal(t, tt.query.Keyword, q.Keyword)
			require.ElementsMatch(t, tt.query.Tags, q.Tags)
			require.ElementsMatch(t, tt.query.Locations, q.Locations)
			require.WithinDuration(t, tt.query.Since, q.Since, 1*time.Second)
			require.WithinDuration(t, tt.query.Until, q.Until, 1*time.Second)
			require.Equal(t, tt.query.SortBy, q.SortBy)
			require.Equal(t, tt.query.SortOrder, q.SortOrder)
		})
	}
}
