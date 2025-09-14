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

package activity

import (
	"fmt"
	"strings"

	"github.com/roma-glushko/frens/internal/journal"

	jctx "github.com/roma-glushko/frens/internal/context"
	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/graph"
	"github.com/roma-glushko/frens/internal/lang"
	"github.com/urfave/cli/v2"
)

var GraphCommand = &cli.Command{
	Name:  "graph",
	Usage: "Print the zen of friendship",
	Flags: []cli.Flag{
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
		&cli.BoolFlag{
			Name:    "unscaled",
			Aliases: []string{"s"},
			Value:   false,
			Usage:   "Disable scaling of the graph bars",
		},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context
		jctx := jctx.FromCtx(ctx)
		s := jctx.Store

		return s.Tx(ctx, func(j *journal.Journal) error {
			activities, err := j.ListEvents(friend.ListEventQuery{
				Type:      friend.EventTypeActivity,
				Keyword:   strings.TrimSpace(c.String("search")),
				Tags:      c.StringSlice("tag"),
				Since:     lang.ExtractDate(c.String("from")),
				Until:     lang.ExtractDate(c.String("to")),
				SortBy:    friend.SortRecency,
				SortOrder: friend.SortOrderDirect,
			})
			if err != nil {
				return fmt.Errorf("failed to list activities: %w", err)
			}

			if len(activities) == 0 {
				fmt.Println("No activities found for given query.")
				return nil
			}

			graph := graph.NewActivityGraph(
				activities,
				j.Activities,
				!c.Bool("unscaled"),
			)

			for _, o := range graph.Output() {
				fmt.Println(o)
			}

			return nil
		})
	},
}
