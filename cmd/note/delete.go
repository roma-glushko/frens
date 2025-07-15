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

package note

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/log"

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
	Usage:     `Delete a note`,
	UsageText: `frens note delete [OPTIONS] [INFO]`,
	Description: `Delete notes from your journal by their ID.
	Examples:
		frens note delete 2zpWoEiUYn6vrSl9w03NAVkWxMn 2zu4V8MAQSvQv9IpAKNYJwaielS
		frens note d -f 2zu4V8MAQSvQv9IpAKNYJwaielS 
	`,
	Args:      true,
	ArgsUsage: `<NOTE_ID> [, <NOTE_ID>...]`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Value:   false,
			Usage:   "Force delete without confirmation",
		},
	},
	Action: func(ctx *cli.Context) error {
		journalDir, err := journaldir.DefaultDir()
		if err != nil {
			return err
		}

		jr, err := journaldir.Load(journalDir)
		if err != nil {
			return err
		}

		if len(ctx.Args().Slice()) == 0 {
			return cli.Exit("Please provide a note ID to delete.", 1)
		}

		events := make([]friend.Event, 0, len(ctx.Args().Slice()))

		for _, actID := range ctx.Args().Slice() {
			act, err := jr.GetEvent(friend.EventTypeNote, actID)
			if err != nil {
				if errors.Is(err, journal.ErrEventNotFound) {
					return cli.Exit("Note not found: "+actID, 1)
				}

				log.Error("Failed to get note", "err", err, "note_id", actID)
				return err
			}

			events = append(events, act)
		}

		actWord := utils.P(len(events), "note", "notes")
		fmt.Printf("🔍 Found %d %s:\n", len(events), actWord)

		for _, act := range events {
			fmt.Printf("   • %s\n", act.ID)
		}

		// TODO: check if interactive mode
		fmt.Println("\n⚠️  You're about to permanently delete the " + actWord + ".")
		if !ctx.Bool("force") && !tui.ConfirmAction("Are you sure?") {
			fmt.Println("\n↩️  Deletion canceled.")
			return nil
		}

		err = journaldir.Update(jr, func(j *journal.Journal) error {
			j.RemoveEvents(friend.EventTypeNote, events)
			return nil
		})
		if err != nil {
			return err
		}

		fmt.Printf("\n🗑️  %s deleted.\n", utils.TitleCaser.String(actWord))

		return nil
	},
}
