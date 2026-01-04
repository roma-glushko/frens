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

package activity

import (
	"fmt"
	"strings"

	jctx "github.com/roma-glushko/frens/internal/context"

	"github.com/roma-glushko/frens/internal/friend"

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
	Usage:     "Update an activity log",
	Args:      true,
	ArgsUsage: `<ACTIVITY_ID> [<DESCRIPTION>]`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "date",
			Aliases: []string{"d"},
			Usage:   "Set the date of the activity (format: YYYY/MM/DD or relative like 'yesterday')",
		},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context
		if c.NArg() < 1 {
			return cli.Exit(
				"You must provide an activity ID to edit. Execute `frens activity ls` to find out.",
				1,
			)
		}

		actID := c.Args().First()
		desc := strings.TrimSpace(strings.Join(c.Args().Slice()[1:], " "))

		appCtx := jctx.FromCtx(ctx)

		return appCtx.Store.Tx(ctx, func(j *journal.Journal) error {
			actOld, err := j.GetEvent(friend.EventTypeActivity, actID)
			if err != nil {
				return cli.Exit("Activity not found: "+actID, 1)
			}

			inputForm := tui.NewEditorForm(tui.EditorOptions{
				Title:      fmt.Sprintf("Edit activity log (%s):", actOld.ID),
				SyntaxHint: lang.FormatEventInfo,
			})
			inputForm.Textarea.SetValue(lang.RenderEvent(actOld))

			// TODO: check if interactive mode is enabled
			teaUI := tea.NewProgram(inputForm, tea.WithMouseAllMotion())

			if _, err := teaUI.Run(); err != nil {
				log.Errorf("uh oh: %v", err)
				return err
			}

			infoTxt := inputForm.Textarea.Value()

			if desc != "" {
				infoTxt = desc
			}

			actNew, err := lang.ExtractEvent(friend.EventTypeActivity, infoTxt)
			if err != nil {
				return cli.Exit("Failed to parse activity description: "+err.Error(), 1)
			}

			if err := actNew.Validate(); err != nil {
				return err
			}

			actNew, err = j.UpdateEvent(actOld, actNew)
			if err != nil {
				return fmt.Errorf("failed to update activity: %w", err)
			}

			log.Info(" ✔ Activity updated")
			log.Info("==> Activity Information\n")

			return appCtx.Printer.Print(actNew)
		})
	},
}
