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
	"strings"
	"time"

	"github.com/roma-glushko/frens/internal/log/formatter"

	jctx "github.com/roma-glushko/frens/internal/context"

	"github.com/roma-glushko/frens/internal/friend"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/lang"
	"github.com/roma-glushko/frens/internal/log"
	"github.com/roma-glushko/frens/internal/tui"
	"github.com/urfave/cli/v2"
)

var AddCommand = &cli.Command{
	Name:    "add",
	Aliases: []string{"a", "new", "create"},
	Usage:   "Add a new activity",
	Args:    true,
	ArgsUsage: `<DESCR>

	<DESCR> is a description of the activity to record.
	
	Examples:
		"Michael wrote a book 'Somehow I managed'" - no date, will be recorded as today
		"yesterday :: Jim Halpert put my stuff in jello #pranks" - relative date & description
		"2009/09/08 :: "Jim and Pam got married at Niagara Falls #theoffice" - absolute date & description
`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "date",
			Aliases: []string{"d"},
			Usage:   "Set the date of the activity (format: YYYY/MM/DD or relative like 'yesterday')",
		},
	},
	Action: func(ctx *cli.Context) error {
		var info string

		if ctx.NArg() == 0 {
			// TODO: also check if we are in the interactive mode
			inputForm := tui.NewEditorForm(tui.EditorOptions{
				Title:      "Add a new activity:",
				SyntaxHint: lang.FormatEventInfo,
			})
			teaUI := tea.NewProgram(inputForm, tea.WithMouseAllMotion())

			if _, err := teaUI.Run(); err != nil {
				log.Errorf("uh oh: %v", err)
				return err
			}

			info = inputForm.Textarea.Value()
		} else {
			info = strings.Join(ctx.Args().Slice(), " ")
		}

		if info == "" {
			return cli.Exit("You must provide a description for the activity.", 1)
		}

		e, err := lang.ExtractEvent(friend.EventTypeActivity, info)
		if err != nil {
			return cli.Exit("Failed to parse activity description: "+err.Error(), 1)
		}

		date := ctx.String("date")

		if date != "" {
			e.Date = lang.ExtractDate(date, time.Now().UTC())
		}

		if err := e.Validate(); err != nil {
			return err
		}

		appCtx := jctx.FromCtx(ctx.Context)

		err = appCtx.Repository.Update(func(j *journal.Journal) error {
			e, err = j.AddEvent(e)
			return err
		})
		if err != nil {
			return err
		}

		log.Infof(" ✔ Activity added")
		log.Info("==> Activity Information\n")

		fmtr := formatter.EventTextFormatter{}

		o, _ := fmtr.FormatSingle(e)
		fmt.Println(o)

		return nil
	},
}
