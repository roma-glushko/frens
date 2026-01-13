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

package note

import (
	"fmt"
	"strings"

	"github.com/roma-glushko/frens/internal/journal"

	"github.com/roma-glushko/frens/internal/log"

	jctx "github.com/roma-glushko/frens/internal/context"

	"github.com/roma-glushko/frens/internal/friend"

	"github.com/roma-glushko/frens/internal/lang"

	"github.com/urfave/cli/v2"
)

var ListCommand = &cli.Command{
	Name:      "list",
	Aliases:   []string{"l", "ls"},
	Usage:     "List all notes",
	UsageText: "frens note list [OPTIONS]",
	Description: `List and filter notes from your journal.

Examples:
  frens note list                            # list all notes
  frens note ls -q "allergic"                # search by keyword
  frens note ls -t health -t travel          # filter by tags
  frens note ls --from 2024/01/01            # notes since a date
  frens note ls --since "last week"          # notes from last week
  frens note ls -s alpha                     # sort alphabetically
`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "search",
			Aliases: []string{"q"},
			Usage:   "Search by keyword",
		},
		&cli.StringSliceFlag{
			Name:    "tag",
			Aliases: []string{"t"},
			Usage:   "Filter by tag(s)",
		},
		&cli.StringFlag{
			Name:    "from",
			Aliases: []string{"since"},
			Usage:   "Filter notes since a specific date",
		},
		&cli.StringFlag{
			Name:    "to",
			Aliases: []string{"until"},
			Usage:   "Filter notes until a specific date",
		},
		&cli.StringFlag{
			Name:    "sort",
			Aliases: []string{"s"},
			Value:   "recency",
			Usage:   "Sort by one of: recency, alpha",
			Action: func(c *cli.Context, s string) error {
				return friend.ValidateEventSortOption(s)
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
			notes, err := j.ListEvents(friend.ListEventQuery{
				Type:      friend.EventTypeNote,
				Keyword:   strings.TrimSpace(c.String("search")),
				Tags:      c.StringSlice("tag"),
				Since:     lang.ExtractDate(c.String("from")),
				Until:     lang.ExtractDate(c.String("to")),
				SortBy:    friend.SortOption(c.String("sort")),
				SortOrder: sortOrder,
			})
			if err != nil {
				return fmt.Errorf("failed to list notes: %w", err)
			}

			if len(notes) == 0 {
				log.Empty("notes")
				return nil
			}

			return appCtx.Printer.PrintList(notes)
		})
	},
}
