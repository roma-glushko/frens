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

package location

import (
	"errors"
	"fmt"
	"strings"

	"github.com/roma-glushko/frens/internal/log/formatter"

	jctx "github.com/roma-glushko/frens/internal/context"

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
	Usage:     "Update main location information",
	Args:      true,
	ArgsUsage: `<LOCATION_NAME, LOCATION_NICKNAME, LOCATION_ID>`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "id",
			Usage: "Location's unique identifier (used for linking with other data, editing, etc.)",
		},
		&cli.StringFlag{
			Name:    "name",
			Aliases: []string{"n"},
			Usage:   "Set location's name",
		},
		&cli.StringFlag{
			Name:    "desc",
			Aliases: []string{"d"},
			Usage:   "Set description of the location",
		},
		&cli.StringSliceFlag{
			Name:    "alias",
			Aliases: []string{"a", "aka", "nick"},
			Usage:   "Set location's aliases (override existing ones)",
		},
	},
	Action: func(ctx *cli.Context) error {
		if ctx.NArg() < 1 {
			return cli.Exit(
				"You must provide a location name, nickname, or ID to edit. Execute `frens location ls` to find out.",
				1,
			)
		}

		lID := strings.Join(ctx.Args().Slice(), " ")

		jctx := jctx.FromCtx(ctx.Context)
		jr := jctx.Journal

		lOld, err := jr.GetLocation(lID)
		if err != nil {
			return err
		}

		inputForm := tui.NewEditorForm(tui.EditorOptions{
			Title:      "Edit " + lOld.Name + " information:",
			SyntaxHint: lang.FormatLocationInfo,
		})
		inputForm.Textarea.SetValue(lang.RenderLocation(lOld))

		// TODO: check if interactive mode is enabled
		teaUI := tea.NewProgram(inputForm, tea.WithMouseAllMotion())

		if _, err := teaUI.Run(); err != nil {
			log.Error("uh oh", "err", err)
			return err
		}

		infoTxt := inputForm.Textarea.Value()

		if infoTxt == "" {
			return errors.New("no location info found")
		}

		lNew, err := lang.ExtractLocation(infoTxt)
		if err != nil {
			return err
		}

		id := ctx.String("id")
		name := ctx.String("name")
		desc := ctx.String("desc")
		aliases := ctx.StringSlice("alias")

		if lNew.ID == "" {
			lNew.ID = lOld.ID
		}

		if id != "" {
			lNew.ID = id
		}

		if name != "" {
			lNew.Name = name
		}

		if desc != "" {
			lNew.Desc = desc
		}

		if len(aliases) > 0 {
			lNew.Aliases = aliases
		}

		if err != nil && !errors.Is(err, lang.ErrNoInfo) {
			log.Error("failed to parse friend info", "err", err)
			return err
		}

		if err := lNew.Validate(); err != nil {
			return err
		}

		err = journaldir.Update(jr, func(j *journal.Journal) error {
			j.UpdateLocation(lOld, lNew)
			return nil
		})
		if err != nil {
			return err
		}

		fmt.Println(" ✔ Location updated")
		log.Info("==> Location Information\n")

		fmtr := formatter.LocationTextFormatter{}

		o, _ := fmtr.FormatSingle(lNew)
		fmt.Println(o)

		return nil
	},
}
