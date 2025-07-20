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
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/roma-glushko/frens/internal/lang"
	"os"
	"text/tabwriter"

	"github.com/roma-glushko/frens/internal/journal"
	"github.com/urfave/cli/v2"
)

var boldNameStyle = lipgloss.NewStyle().Bold(true)

var ListCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"l", "ls"},
	Usage:   "List all friends",
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
				return journal.ValidateSortOption(s)
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
		jr := journal.FromCtx(ctx)

		orderBy := journal.OrderDirect

		if c.Bool("reverse") {
			orderBy = journal.OrderReverse
		}

		friends := jr.ListFriends(journal.ListFriendQuery{
			Search:    c.String("search"),
			Locations: c.StringSlice("location"),
			Tags:      c.StringSlice("tag"),
			SortBy:    journal.SortOption(c.String("sort")),
			OrderBy:   orderBy,
		})

		if len(friends) == 0 {
			fmt.Println("No friends found")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", "", "", "", "")
		_, _ = fmt.Fprintln(w, "\t👤  Name\t🏷️  Tags\t📍  Location")
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", "", "", "", "")

		for _, f := range friends {
			// TODO: improve output formatting
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", f.ID, boldNameStyle.Render(f.String()), lang.RenderTags(f.Tags), lang.RenderLocMarkers(f.Locations))
		}

		_ = w.Flush()

		return nil
	},
}
