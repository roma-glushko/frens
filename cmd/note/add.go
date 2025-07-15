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
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/journaldir"
	"github.com/roma-glushko/frens/internal/lang"
	"github.com/roma-glushko/frens/internal/tui"
	"github.com/urfave/cli/v2"
)

var AddCommand = &cli.Command{
	Name:      "add",
	Aliases:   []string{"a", "create", "new"},
	Usage:     "Add a new note",
	UsageText: `frens note add [options]`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "date",
			Aliases: []string{"d"},
			Usage:   "Set the date of the note (format: YYYY/MM/DD or relative like 'yesterday')",
		},
	},
	Action: func(ctx *cli.Context) error {
		var info string

		if ctx.NArg() == 0 {
			// TODO: also check if we are in the interactive mode
			inputForm := tui.NewEditorForm(tui.EditorOptions{
				Title:      "Add a new note:",
				SyntaxHint: lang.FormatEventInfo,
			})
			teaUI := tea.NewProgram(inputForm, tea.WithMouseAllMotion())

			if _, err := teaUI.Run(); err != nil {
				log.Error("uh oh", "err", err)
				return err
			}

			info = inputForm.Textarea.Value()
		} else {
			info = strings.Join(ctx.Args().Slice(), " ")
		}

		if info == "" {
			return cli.Exit("You must provide a description for the note.", 1)
		}

		e, err := lang.ExtractEvent(friend.EventTypeNote, info)
		if err != nil {
			return cli.Exit("Failed to parse note description: "+err.Error(), 1)
		}

		date := ctx.String("date")

		if date != "" {
			e.Date = lang.ExtractDate(date)
		}

		if err := e.Validate(); err != nil {
			return err
		}

		jr := journal.FromCtx(ctx.Context)

		err = journaldir.Update(jr, func(j *journal.Journal) error {
			e, err = j.AddEvent(e)

			return err
		})
		if err != nil {
			return err
		}

		fmt.Println("✅ Added new note: " + e.ID)

		return nil
	},
}
