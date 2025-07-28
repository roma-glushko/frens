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

	jctx "github.com/roma-glushko/frens/internal/context"

	"github.com/roma-glushko/frens/internal/utils"

	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/journaldir"
	"github.com/roma-glushko/frens/internal/tui"
	"github.com/urfave/cli/v2"
)

var DeleteCommand = &cli.Command{
	Name:      "delete",
	Aliases:   []string{"del", "rm", "d"},
	Usage:     `Delete an activity`,
	UsageText: `frens activity delete [OPTIONS] [INFO]`,
	Description: `Delete activity logs from your journal by their ID.
	Examples:
		frens activity delete 2zpWoEiUYn6vrSl9w03NAVkWxMn 2zpWoEiUYn6vrSl9w03NAVkWxMx
		frens activity d -f 2zpWoEiUYn6vrSl9w03NAVkWxMn 
	`,
	Args:      true,
	ArgsUsage: `<ACTIVITY_ID> [, <ACTIVITY_ID>...]`,
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
			return cli.Exit("Please provide a activity ID to delete.", 1)
		}

		activities := make([]*friend.Event, 0, len(c.Args().Slice()))

		jctx := jctx.FromCtx(c.Context)
		jr := jctx.Journal

		for _, actID := range c.Args().Slice() {
			act, err := jr.GetEvent(friend.EventTypeActivity, actID)
			if err != nil {
				return err
			}

			activities = append(activities, act)
		}

		actWord := utils.P(len(activities), "activity", "activities")
		fmt.Printf("🔍 Found %d %s:\n", len(activities), actWord)

		for _, act := range activities {
			fmt.Printf("   • %s\n", act.ID)
		}

		// TODO: check if interactive mode
		fmt.Println("\n⚠️  You're about to permanently delete the " + actWord + ".")
		if !c.Bool("force") && !tui.ConfirmAction("Are you sure?") {
			fmt.Println("\n↩️  Deletion canceled.")
			return nil
		}

		err := journaldir.Update(jr, func(j *journal.Journal) error {
			j.RemoveEvents(friend.EventTypeActivity, activities)
			return nil
		})
		if err != nil {
			return err
		}

		fmt.Printf("\n🗑️  %s deleted.\n", utils.TitleCaser.String(actWord))

		return nil
	},
}
