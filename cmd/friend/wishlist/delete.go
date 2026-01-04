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
	"fmt"

	jctx "github.com/roma-glushko/frens/internal/context"
	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/tui"
	"github.com/roma-glushko/frens/internal/utils"
	"github.com/urfave/cli/v2"
)

var DeleteCommand = &cli.Command{
	Name:      "delete",
	Aliases:   []string{"del", "rm", "d"},
	Usage:     `Delete a wishlist item from your friend`,
	UsageText: `frens friend wishlist delete [OPTIONS] <WISHLIST_ITEM_ID> [, <WISHLIST_ITEM_ID>...]`,
	Description: `Delete friend's wishlist items from your journal by item IDs.
	Examples:
		frens friend wishlist delete 2zpWoEiUYn6vrSl9w03NAVkWxMn
		frens friend wishlist d -f 2zpWoEiUYn6vrSl9w03NAVkWxMn
	`,
	Args:      true,
	ArgsUsage: `<WISHLIST_ITEM_ID> [, <WISHLIST_ITEM_ID>...]`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Value:   false,
			Usage:   "Force delete without confirmation",
		},
	},
	Action: func(c *cli.Context) error {
		if len(c.Args().Slice()) == 0 {
			return cli.Exit("Please provide a wishlist item ID to delete.", 1)
		}

		ctx := c.Context
		jctx := jctx.FromCtx(ctx)
		s := jctx.Store

		return s.Tx(ctx, func(j *journal.Journal) error {
			items := make([]friend.WishlistItem, 0, len(c.Args().Slice()))

			for _, wID := range c.Args().Slice() {
				w, err := j.GetFriendWishlistItem(wID)
				if err != nil {
					return err
				}

				items = append(items, w)
			}

			itemWord := utils.P(len(items), "item", "items")
			fmt.Printf("Found %d wishlist %s:\n", len(items), itemWord)

			for _, item := range items {
				desc := item.Desc
				if desc == "" {
					desc = item.Link
				}
				fmt.Printf("   %s: %s\n", item.ID, desc)
			}

			fmt.Println("\n  You're about to permanently delete the wishlist " + itemWord + ".")
			if !c.Bool("force") && !tui.ConfirmAction("Are you sure?") {
				fmt.Println("\n  Deletion canceled.")
				return nil
			}

			if err := j.RemoveFriendWishlistItems(items); err != nil {
				return err
			}

			fmt.Printf("\n  Wishlist %s deleted.\n", itemWord)

			return nil
		})
	},
}
