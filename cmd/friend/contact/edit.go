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
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	jctx "github.com/roma-glushko/frens/internal/context"
	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/lang"
	"github.com/roma-glushko/frens/internal/log"
	"github.com/roma-glushko/frens/internal/tui"
	"github.com/urfave/cli/v2"
)

var EditCommand = &cli.Command{
	Name:      "edit",
	Aliases:   []string{"e", "modify", "update"},
	Usage:     "Update a contact",
	Args:      true,
	ArgsUsage: `<CONTACT_ID>`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "type",
			Aliases: []string{"tp"},
			Usage:   "Contact type (email, phone, telegram, etc.)",
		},
		&cli.StringFlag{
			Name:    "value",
			Aliases: []string{"v"},
			Usage:   "Contact value",
		},
		&cli.StringSliceFlag{
			Name:    "tag",
			Aliases: []string{"t"},
			Usage:   "Set tags for the contact",
		},
	},
	Action: func(c *cli.Context) error {
		if c.NArg() < 1 {
			return cli.Exit(
				"You must provide a contact ID. Execute `frens friend contact ls` to find out.",
				1,
			)
		}

		cID := strings.Join(c.Args().Slice(), " ")

		ctx := c.Context
		appCtx := jctx.FromCtx(ctx)
		s := appCtx.Store

		return s.Tx(ctx, func(j *journal.Journal) error {
			cOld, err := j.GetFriendContact(cID)
			if err != nil {
				return err
			}

			inputForm := tui.NewEditorForm(tui.EditorOptions{
				Title:      "Edit contact (" + cOld.ID + "):",
				SyntaxHint: lang.FormatContactInfo,
			})

			inputForm.Textarea.SetValue(lang.RenderContact(cOld))

			teaUI := tea.NewProgram(inputForm, tea.WithMouseAllMotion())

			if _, err := teaUI.Run(); err != nil {
				log.Errorf("uh oh: %v", err)
				return err
			}

			infoTxt := inputForm.Textarea.Value()

			if infoTxt == "" {
				return errors.New("no contact info provided")
			}

			cNew, err := lang.ExtractContact(infoTxt)
			if err != nil && !errors.Is(err, lang.ErrNoInfo) {
				log.Errorf("failed to parse contact info: %v", err)
				return err
			}

			// Apply CLI overrides
			if tp := c.String("type"); tp != "" {
				cNew.Type = friend.ParseContactType(tp)
			}

			if v := c.String("value"); v != "" {
				cNew.Value = v
			}

			if tags := c.StringSlice("tag"); len(tags) > 0 {
				cNew.Tags = tags
			}

			if err := cNew.Validate(); err != nil {
				return err
			}

			cNew, err = j.UpdateFriendContact(cOld, cNew)
			if err != nil {
				return err
			}

			log.Info(" Contact updated")
			log.Info("==> Contact Information\n")

			return appCtx.Printer.Print(cNew)
		})
	},
}
