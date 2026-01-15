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

package contact

import (
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	jctx "github.com/roma-glushko/frens/internal/context"
	"github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/lang"
	"github.com/roma-glushko/frens/internal/log"
	"github.com/roma-glushko/frens/internal/tui"
	"github.com/urfave/cli/v2"
)

var AddCommand = &cli.Command{
	Name:      "add",
	Aliases:   []string{"a", "new", "create"},
	Usage:     "Add contact information for a friend",
	UsageText: "frens friend contact add [OPTIONS] <FRIEND_NAME, FRIEND_NICKNAME, FRIEND_ID> [CONTACTS]",
	Args:      true,
	ArgsUsage: `<CONTACTS>
		Add one or more contacts in a flexible format.

		Format: ` + lang.FormatContactInfo + `

		Examples:
			frens friend contact add "John Doe" ig:@johndoe +1234567890 john@example.com
			frens friend contact add john x:@johndoe tg:@john_doe #social
	`,
	Flags: []cli.Flag{
		&cli.StringSliceFlag{
			Name:    "tag",
			Aliases: []string{"t"},
			Usage:   "Add tags to all contacts (e.g., 'work', 'personal')",
		},
	},
	Action: func(c *cli.Context) error {
		var info string

		if c.NArg() == 0 {
			return cli.Exit("You must provide a friend name, nickname, or ID to add contacts.", 1)
		}

		ctx := c.Context
		appCtx := jctx.FromCtx(ctx)

		pID := c.Args().First()

		return appCtx.Store.Tx(ctx, func(j *journal.Journal) error {
			p, err := j.GetFriend(pID)
			if err != nil {
				return err
			}

			if c.NArg() == 1 {
				inputForm := tui.NewEditorForm(tui.EditorOptions{
					Title:      "Add contact information for " + p.Name + ":",
					SyntaxHint: lang.FormatContactInfo,
				})
				teaUI := tea.NewProgram(inputForm, tea.WithMouseAllMotion())

				if _, err := teaUI.Run(); err != nil {
					log.Errorf("uh oh: %v", err)
					return err
				}

				info = inputForm.Textarea.Value()
			} else {
				info = strings.Join(c.Args().Slice()[1:], " ")
			}

			if info == "" {
				return errors.New("no contact information provided")
			}

			contacts, err := lang.ExtractContacts(info)
			if err != nil && !errors.Is(err, lang.ErrNoInfo) {
				log.Errorf("failed to parse contact info: %v", err)
				return err
			}

			// Apply CLI tags to all contacts
			cliTags := c.StringSlice("tag")
			if len(cliTags) > 0 {
				for i := range contacts {
					contacts[i].Tags = append(contacts[i].Tags, cliTags...)
				}
			}

			// Validate and add each contact
			addedContacts := make([]string, 0, len(contacts))

			for _, contact := range contacts {
				if err := contact.Validate(); err != nil {
					return fmt.Errorf("invalid contact %s: %w", contact.Value, err)
				}

				added, err := j.AddFriendContact(p.ID, contact)
				if err != nil {
					return err
				}

				addedContacts = append(
					addedContacts,
					fmt.Sprintf("%s: %s", added.Type, added.Value),
				)
			}

			log.Infof(" Added %d contact(s) for %s", len(addedContacts), p.Name)

			for _, contact := range contacts {
				if err := appCtx.Printer.Print(contact); err != nil {
					return err
				}
			}

			return nil
		})
	},
}
