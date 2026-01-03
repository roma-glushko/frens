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

package contact

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
	Usage:     `Delete a contact`,
	UsageText: `frens friend contact delete [OPTIONS] <CONTACT_ID> [, <CONTACT_ID>...]`,
	Description: `Delete friend's contacts from your journal by contact IDs.
	Examples:
		frens friend contact delete 2zpWoEiUYn6vrSl9w03NAVkWxMn
		frens friend contact d -f 2zpWoEiUYn6vrSl9w03NAVkWxMn
	`,
	Args:      true,
	ArgsUsage: `<CONTACT_ID> [, <CONTACT_ID>...]`,
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
			return cli.Exit("Please provide a contact ID to delete.", 1)
		}

		ctx := c.Context
		jctx := jctx.FromCtx(ctx)
		s := jctx.Store

		return s.Tx(ctx, func(j *journal.Journal) error {
			contacts := make([]friend.Contact, 0, len(c.Args().Slice()))

			for _, cID := range c.Args().Slice() {
				contact, err := j.GetFriendContact(cID)
				if err != nil {
					return err
				}

				contacts = append(contacts, contact)
			}

			contactWord := utils.P(len(contacts), "contact", "contacts")
			fmt.Printf("Found %d %s:\n", len(contacts), contactWord)

			for _, contact := range contacts {
				fmt.Printf("   %s (%s): %s\n", contact.ID, contact.Type, contact.Value)
			}

			fmt.Println("\n  You're about to permanently delete the " + contactWord + ".")
			if !c.Bool("force") && !tui.ConfirmAction("Are you sure?") {
				fmt.Println("\n  Deletion canceled.")
				return nil
			}

			if err := j.RemoveFriendContacts(contacts); err != nil {
				return err
			}

			fmt.Printf("\n  %s deleted.\n", utils.TitleCaser.String(contactWord))

			return nil
		})
	},
}
