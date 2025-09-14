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
	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/tui"
	"github.com/roma-glushko/frens/internal/utils"
	"github.com/urfave/cli/v2"
)

var DeleteCommand = &cli.Command{
	Name:      "delete",
	Aliases:   []string{"del", "rm", "d"},
	Usage:     `Delete a date from your friend`,
	UsageText: `frens friend date delete [OPTIONS] [INFO]`,
	Description: `Delete friend's date from your journal by date IDs.
	Examples:
		frens friend date delete 2zpWoEiUYn6vrSl9w03NAVkWxMn 2zpWoEiUYn6vrSl9w03NAVkWxMx
		frens friend date d -f 2zpWoEiUYn6vrSl9w03NAVkWxMn 
	`,
	Args:      true,
	ArgsUsage: `<DATE_ID> [, <DATE_ID>...]`,
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
			return cli.Exit("Please provide a date ID to delete.", 1)
		}

		dates := make([]friend.Date, 0, len(c.Args().Slice()))

		ctx := c.Context
		jctx := jctx.FromCtx(ctx)
		s := jctx.Store

		return s.Tx(ctx, func(j *journal.Journal) error {
			for _, actID := range c.Args().Slice() {
				dt, err := j.GetFriendDate(actID)
				if err != nil {
					return err
				}

				dates = append(dates, dt)
			}

			dtWord := utils.P(len(dates), "date", "dates")
			fmt.Printf("🔍 Found %d %s:\n", len(dates), dtWord)

			for _, act := range dates {
				fmt.Printf("   • %s\n", act.ID)
			}

			// TODO: check if interactive mode
			fmt.Println("\n⚠️  You're about to permanently delete the " + dtWord + ".")
			if !c.Bool("force") && !tui.ConfirmAction("Are you sure?") {
				fmt.Println("\n↩️  Deletion canceled.")
				return nil
			}

			err := j.RemoveFriendDates(dates)
			if err != nil {
				return err
			}

			fmt.Printf("\n🗑️  %s deleted.\n", utils.TitleCaser.String(dtWord))

			return nil
		})
	},
}
