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

package wishlist

import (
	jctx "github.com/roma-glushko/frens/internal/context"
	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/log"

	"github.com/urfave/cli/v2"
)

var ListCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"l", "ls"},
	Usage:   "List wishlist items for all friends",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "search",
			Aliases: []string{"q"},
			Usage:   "Search by description or link",
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
		s := appCtx.Store

		return s.Tx(ctx, func(j *journal.Journal) error {
			items, err := j.ListFriendWishlistItems(friend.ListWishlistQuery{
				Keyword: c.String("search"),
				Friends: c.StringSlice("with"),
				Tags:    c.StringSlice("tag"),
			})
			if err != nil {
				return err
			}

			if len(items) == 0 {
				log.Info("No wishlist items found for given query.")
				return nil
			}

			return appCtx.Printer.PrintList(items)
		})
	},
}
