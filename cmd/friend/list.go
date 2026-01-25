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

package friend

import (
	"strings"

	"github.com/roma-glushko/frens/internal/journal"

	jctx "github.com/roma-glushko/frens/internal/context"

	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/log"

	"github.com/urfave/cli/v2"
)

var ListCommand = &cli.Command{
	Name:      "list",
	Aliases:   []string{"l", "ls"},
	Usage:     "List all friends",
	UsageText: "frens friend list [OPTIONS]",
	Description: `List and filter friends from your journal.

Examples:
  frens friend list                          # list all friends
  frens friend ls -q "Jim"                   # search by name or description
  frens friend ls -t work -t college         # filter by multiple tags
  frens friend ls -l NYC -l Scranton         # filter by locations
  frens friend ls -s recency                 # sort by most recent activity
  frens friend ls -s activities -r           # sort by activity count, reversed
  frens friend ls -t family -s alpha         # combine filters and sorting
`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "search",
			Aliases: []string{"q"},
			Usage:   "Search by name or description",
		},
		&cli.StringSliceFlag{
			Name:    "location",
			Aliases: []string{"l", "loc", "in"},
			Usage:   "List friends by location(s)",
		},
		&cli.StringSliceFlag{
			Name:    "tag",
			Aliases: []string{"t"},
			Usage:   "Filter by tag(s)",
		},
		&cli.StringFlag{
			Name:    "sort",
			Aliases: []string{"s"},
			Value:   "alpha",
			Usage:   "Sort by one of alpha, activities, recency",
			Action: func(c *cli.Context, s string) error {
				return friend.ValidateEntitySortOption(s)
			},
		},
		&cli.BoolFlag{
			Name:    "reverse",
			Aliases: []string{"r"},
			Value:   false,
			Usage:   "Reverse sort order",
		},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context
		appCtx := jctx.FromCtx(ctx)
		s := appCtx.Store

		sortOrder := friend.SortOrderDirect

		if c.Bool("reverse") {
			sortOrder = friend.SortOrderReverse
		}

		return s.Tx(ctx, func(j *journal.Journal) error {
			friends := j.ListFriends(friend.ListFriendQuery{
				Keyword:   strings.TrimSpace(c.String("search")),
				Locations: c.StringSlice("location"),
				Tags:      c.StringSlice("tag"),
				SortBy:    friend.SortOption(c.String("sort")),
				SortOrder: sortOrder,
			})

			if len(friends) == 0 {
				log.Empty("friends")
				return nil
			}

			return appCtx.Printer.PrintList(friends)
		})
	},
}
