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

package date

import (
	"errors"
	"fmt"
	"strings"

	"github.com/roma-glushko/frens/internal/log/formatter"

	jctx "github.com/roma-glushko/frens/internal/context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/journaldir"
	"github.com/roma-glushko/frens/internal/lang"
	"github.com/roma-glushko/frens/internal/log"
	"github.com/roma-glushko/frens/internal/tui"
	"github.com/urfave/cli/v2"
)

var EditCommand = &cli.Command{
	Name:      "edit",
	Aliases:   []string{"e", "modify", "update"},
	Usage:     "Update date information",
	Args:      true,
	ArgsUsage: `<DATE_ID>`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "desc",
			Usage: "Description of the date",
		},
		&cli.StringFlag{
			Name:    "date",
			Aliases: []string{"d"},
			Usage:   "Date in a free-form format (e.g., 'May 13th', '1996-7-30', '1985')",
		},
		&cli.StringFlag{
			Name:    "calendar",
			Aliases: []string{"cal"},
			Usage:   "Calendar type to use for the date (e.g., gregorian, hebrew)",
		},
		&cli.StringSliceFlag{
			Name:    "tag",
			Aliases: []string{"t"},
			Usage:   "Add tags to the date (e.g., 'birthday', 'anniversary')",
		},
	},
	Action: func(ctx *cli.Context) error {
		if ctx.NArg() < 1 {
			return cli.Exit(
				"You must provide a date ID. Execute `frens friend dt ls` to find out.",
				1,
			)
		}

		dtID := strings.Join(ctx.Args().Slice(), " ")

		jctx := jctx.FromCtx(ctx.Context)
		jr := jctx.Journal

		dtOld, err := jr.GetFriendDate(dtID)
		if err != nil {
			return err
		}

		inputForm := tui.NewEditorForm(tui.EditorOptions{
			Title:      "Edit " + dtOld.ID + " information:",
			SyntaxHint: lang.FormatDateInfo,
		})

		dateInfo := lang.RenderDateInfo(dtOld)
		inputForm.Textarea.SetValue(dateInfo)

		// TODO: check if interactive mode is enabled
		teaUI := tea.NewProgram(inputForm, tea.WithMouseAllMotion())

		if _, err := teaUI.Run(); err != nil {
			log.Errorf("uh oh: %v", err)
			return err
		}

		infoTxt := inputForm.Textarea.Value()

		if infoTxt == "" {
			return errors.New("no date info provided")
		}

		dtNew, err := lang.ExtractDateInfo(infoTxt)

		date := ctx.String("date")
		desc := ctx.String("desc")
		cal := ctx.String("calendar")
		tags := ctx.StringSlice("tag")

		if date != "" {
			dtNew.DateExpr = date
		}

		if desc != "" {
			dtNew.Desc = desc
		}

		if cal != "" {
			dtNew.Calendar = cal
		}

		if len(tags) > 0 {
			dtNew.Tags = tags
		}

		if err != nil && !errors.Is(err, lang.ErrNoInfo) {
			log.Errorf(" ✖ failed to parse friend info: %v", err)
			return err
		}

		if err := dtNew.Validate(); err != nil {
			return err
		}

		err = journaldir.Update(jr, func(j *journal.Journal) error {
			dtNew, err = j.UpdateFriendDate(dtOld, dtNew)

			return err
		})
		if err != nil {
			return err
		}

		log.Info(" ✔ Date updated")
		log.Info("==> Date Information\n")

		fmtr := formatter.DateTextFormatter{}

		o, _ := fmtr.FormatSingle(dtNew)
		fmt.Println(o)

		return nil
	},
}
