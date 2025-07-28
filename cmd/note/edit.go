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
	"fmt"
	jctx "github.com/roma-glushko/frens/internal/context"
	"strings"

	"github.com/roma-glushko/frens/internal/friend"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/journaldir"
	"github.com/roma-glushko/frens/internal/lang"
	"github.com/roma-glushko/frens/internal/tui"
	"github.com/urfave/cli/v2"
)

var EditCommand = &cli.Command{
	Name:      "edit",
	Aliases:   []string{"e", "modify", "update"},
	Usage:     "Update a note",
	Args:      true,
	ArgsUsage: `<NOTE_ID> [<DESCRIPTION>]`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "date",
			Aliases: []string{"d"},
			Usage:   "Set the date of the note (format: YYYY/MM/DD or relative like 'yesterday')",
		},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		if c.NArg() < 1 {
			return cli.Exit("You must provide an note ID to edit.", 1)
		}

		actID := c.Args().First()
		desc := strings.TrimSpace(strings.Join(c.Args().Slice()[1:], " "))

		jctx := jctx.FromCtx(ctx)
		jr := jctx.Journal

		actOld, err := jr.GetEvent(friend.EventTypeNote, actID)
		if err != nil {
			return cli.Exit("Note not found: "+actID, 1)
		}

		inputForm := tui.NewEditorForm(tui.EditorOptions{
			Title:      "Edit note (" + actOld.ID + "):",
			SyntaxHint: lang.FormatEventInfo,
		})
		inputForm.Textarea.SetValue(lang.RenderEvent(actOld))

		// TODO: check if interactive mode is enabled
		teaUI := tea.NewProgram(inputForm, tea.WithMouseAllMotion())

		if _, err := teaUI.Run(); err != nil {
			log.Error("uh oh", "err", err)
			return err
		}

		infoTxt := inputForm.Textarea.Value()

		if desc != "" {
			infoTxt = desc
		}

		actNew, err := lang.ExtractEvent(friend.EventTypeNote, infoTxt)
		if err != nil {
			return cli.Exit("Failed to parse note description: "+err.Error(), 1)
		}

		if err := actNew.Validate(); err != nil {
			return err
		}

		err = journaldir.Update(jr, func(j *journal.Journal) error {
			actNew, err = j.UpdateEvent(*actOld, actNew)
			return err
		})
		if err != nil {
			return err
		}

		fmt.Println("🔄 Updated note: " + actNew.ID)

		return nil
	},
}
