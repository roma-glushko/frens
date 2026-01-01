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

package date

import (
	"fmt"

	jctx "github.com/roma-glushko/frens/internal/context"
	"github.com/roma-glushko/frens/internal/log/formatter"

	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/log"

	"github.com/urfave/cli/v2"
)

var ListCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"l", "ls"},
	Usage:   "List all dates for all friends",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "search",
			Aliases: []string{"q"},
			Usage:   "Search by description",
		},
		&cli.StringSliceFlag{
			Name:    "with",
			Aliases: []string{"w"},
			Usage:   "Filter by friend(s)",
		},
		&cli.StringSliceFlag{
			Name:    "tag",
			Aliases: []string{"t"},
			Usage:   "Filter by tag(s)",
		},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context
		appCtx := jctx.FromCtx(ctx)
		jr := appCtx.Repository.Journal()

		dates, err := jr.ListFriendDates(friend.ListDateQuery{
			Keyword: c.String("search"),
			Friends: c.StringSlice("with"),
			Tags:    c.StringSlice("tag"),
		})
		if err != nil {
			return err
		}

		if len(dates) == 0 {
			log.Info("No dates found for given query.")
			return nil
		}

		fmtr := formatter.DateTextFormatter{}

		o, _ := fmtr.FormatList(dates)
		fmt.Println(o)

		return nil
	},
}
