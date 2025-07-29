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

package location

import (
	"fmt"
	"strings"

	jctx "github.com/roma-glushko/frens/internal/context"
	"github.com/roma-glushko/frens/internal/log/formatter"

	"github.com/roma-glushko/frens/internal/friend"

	"github.com/urfave/cli/v2"
)

var ListCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"l", "ls"},
	Usage:   "List all locations",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "search",
			Aliases: []string{"q"},
			Usage:   "Search by name or description",
		},
		&cli.StringSliceFlag{
			Name:    "country",
			Aliases: []string{"c"},
			Usage:   "Filter locations by country",
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
				return friend.ValidateSortOption(s)
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
		jctx := jctx.FromCtx(ctx)
		jr := jctx.Journal

		sortOrder := friend.SortOrderDirect

		if c.Bool("reverse") {
			sortOrder = friend.SortOrderReverse
		}

		locations := jr.ListLocations(friend.ListLocationQuery{
			Keyword:   strings.TrimSpace(c.String("search")),
			Countries: c.StringSlice("country"),
			Tags:      c.StringSlice("tag"),
			SortBy:    friend.SortOption(c.String("sort")),
			SortOrder: sortOrder,
		})

		if len(locations) == 0 {
			fmt.Println("No locations found")
			return nil
		}

		fmtr := formatter.LocationTextFormatter{}

		o, _ := fmtr.FormatList(locations)
		fmt.Println(o)

		return nil
	},
}
