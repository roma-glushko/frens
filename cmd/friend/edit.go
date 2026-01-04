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

package friend

import (
	"errors"
	"fmt"
	"strings"

	"github.com/roma-glushko/frens/internal/log/formatter"

	jctx "github.com/roma-glushko/frens/internal/context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/lang"
	"github.com/roma-glushko/frens/internal/log"
	"github.com/roma-glushko/frens/internal/tui"
	"github.com/urfave/cli/v2"
)

var EditCommand = &cli.Command{
	Name:      "edit",
	Aliases:   []string{"e", "modify", "update"},
	Usage:     "Update main friend information",
	Args:      true,
	ArgsUsage: `<FRIEND_NAME, FRIEND_NICKNAME, FRIEND_ID>`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "id",
			Usage: "Set friend's unique identifier (used for linking with other data, editing, etc.)",
		},
		&cli.StringFlag{
			Name:    "name",
			Aliases: []string{"n"},
			Usage:   "Set friend's name",
		},
		&cli.StringFlag{
			Name:    "desc",
			Aliases: []string{"d"},
			Usage:   "Set description of the friend",
		},
		&cli.StringSliceFlag{
			Name:    "nickname",
			Aliases: []string{"a", "aka", "alias", "nick"},
			Usage:   "Set friend's nicknames (override existing ones)",
		},
	},
	Action: func(c *cli.Context) error {
		if c.NArg() < 1 {
			return cli.Exit(
				"You must provide a friend name, nickname, or ID to edit. Execute `frens friend ls` to find out.",
				1,
			)
		}

		pID := strings.Join(c.Args().Slice(), " ")

		ctx := c.Context
		appCtx := jctx.FromCtx(c.Context)

		return appCtx.Store.Tx(ctx, func(j *journal.Journal) error {
			pOld, err := j.GetFriend(pID)
			if err != nil {
				return err
			}

			inputForm := tui.NewEditorForm(tui.EditorOptions{
				Title:      "Edit " + pOld.Name + " information:",
				SyntaxHint: lang.FormatPersonInfo,
			})
			inputForm.Textarea.SetValue(lang.RenderPerson(pOld))

			// TODO: check if interactive mode is enabled
			teaUI := tea.NewProgram(inputForm, tea.WithMouseAllMotion())

			if _, err := teaUI.Run(); err != nil {
				log.Errorf("uh oh: %v", err)
				return err
			}

			infoTxt := inputForm.Textarea.Value()

			if infoTxt == "" {
				return errors.New("no friend info provided")
			}

			pNew, err := lang.ExtractPerson(infoTxt)

			id := c.String("id")
			name := c.String("name")
			desc := c.String("desc")
			nicknames := c.StringSlice("nickname")

			if pNew.ID == "" {
				pNew.ID = pOld.ID
			}

			if id != "" {
				pNew.ID = id
			}

			if name != "" {
				pNew.Name = name
			}

			if desc != "" {
				pNew.Desc = desc
			}

			if len(nicknames) > 0 {
				pNew.Nicknames = nicknames
			}

			if err != nil && !errors.Is(err, lang.ErrNoInfo) {
				log.Errorf(" ✖ failed to parse friend info: %v", err)
				return err
			}

			if err := pNew.Validate(); err != nil {
				return err
			}

			j.UpdateFriend(pOld, pNew)

			log.Info(" ✔ Friend updated")
			log.Info("==> Friend Information\n")

			fmtr := formatter.PersonTextFormatter{}

			o, _ := fmtr.FormatSingle(pNew)
			fmt.Println(o)

			return nil
		})
	},
}
