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

	jctx "github.com/roma-glushko/frens/internal/context"

	"github.com/roma-glushko/frens/internal/utils"

	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/journaldir"
	"github.com/roma-glushko/frens/internal/tui"
	"github.com/urfave/cli/v2"
)

var (
	labelStyle = lipgloss.NewStyle().Bold(true)
	warnStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("3")) // red
)

var DeleteCommand = &cli.Command{
	Name:      "delete",
	Aliases:   []string{"del", "rm", "d"},
	Usage:     `Delete a friend`,
	UsageText: `frens friend delete [OPTIONS] [INFO]`,
	Description: `Delete friends from your journal by their name, nickname, or ID.
	Examples:
		frens friend delete "Toby Flenderson"
		frens friend d -f "Toby Flenderson"
	`,
	Args:      true,
	ArgsUsage: `<FRIEND_NAME, FRIEND_NICKNAME, FRIEND_ID> [...]`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Value:   false,
			Usage:   "Force delete without confirmation",
		},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		if len(c.Args().Slice()) == 0 {
			return cli.Exit("Please provide a friend name, nickname, or ID to delete.", 1)
		}

		friends := make([]friend.Person, 0, len(c.Args().Slice()))

		jctx := jctx.FromCtx(ctx)
		jr := jctx.Journal

		for _, fID := range c.Args().Slice() {
			f, err := jr.GetFriend(fID)
			if err != nil {
				return err
			}

			friends = append(friends, f)
		}

		frenWord := utils.P(len(friends), "friend", "friends")
		fmt.Printf("\n Found %d %s:\n\n", len(friends), frenWord)

		for _, f := range friends {
			fmt.Printf(" • %s [%s]\n", labelStyle.Render(f.String()), f.ID)
		}

		// TODO: check if interactive mode
		fmt.Println(warnStyle.Render("\n You're about to permanently delete the " + frenWord + "."))
		if !c.Bool("force") && !tui.ConfirmAction(warnStyle.Render(" Are you sure?")) {
			fmt.Println("\n ↩ Deletion canceled.")
			return nil
		}

		err := journaldir.Update(jr, func(j *journal.Journal) error {
			j.RemoveFriends(friends)
			return nil
		})
		if err != nil {
			return err
		}

		fmt.Printf("\n ✔ %s deleted.\n", utils.TitleCaser.String(frenWord))

		return nil
	},
}
