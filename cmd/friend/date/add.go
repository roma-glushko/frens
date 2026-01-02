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
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/roma-glushko/frens/internal/tui"

	jctx "github.com/roma-glushko/frens/internal/context"

	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/lang"
	"github.com/roma-glushko/frens/internal/log"
	"github.com/urfave/cli/v2"
)

var AddCommand = &cli.Command{
	Name:      "add",
	Aliases:   []string{"a", "new", "create"},
	Usage:     "Add a new date to a friend",
	UsageText: "frens friend date add [OPTIONS] <FRIEND_NAME, FRIEND_NICKNAME, FRIEND_ID> [INFO]",
	Args:      true,
	ArgsUsage: `<INFO>
		If no arguments are provided, a textarea will be shown to fill in the details interactively.
		Otherwise, the information will be parsed from the command options.
		
		<INFO> format:
			` + lang.FormatDateInfo + `

		For example:
			"birthday :: May 13th"
			"anniversary :: 2009-9-09"
	`,
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
	Action: func(c *cli.Context) error {
		var info string

		if c.NArg() == 0 {
			return cli.Exit("You must provide a friend name, nickname, or ID to add a date.", 1)
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
				// TODO: also check if we are in the interactive mode
				inputForm := tui.NewEditorForm(tui.EditorOptions{
					Title:      "Add a new friend date information:",
					SyntaxHint: lang.FormatDateInfo,
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

			var d friend.Date

			if info != "" {
				d, err = lang.ExtractDateInfo(info)

				if err != nil && !errors.Is(err, lang.ErrNoInfo) {
					log.Errorf("failed to parse date info: %v", err)
					return err
				}
			}

			desc := c.String("desc")
			dateExpr := c.String("date")
			calendar := c.String("calendar") // TODO: parse and validate calendar
			tags := c.StringSlice("tag")

			if desc != "" {
				d.Desc = desc
			}

			if dateExpr != "" {
				d.DateExpr = dateExpr
			}

			if calendar != "" {
				d.Calendar = calendar
			}

			if len(tags) > 0 {
				d.Tags = tags
			}

			if err := d.Validate(); err != nil {
				return err
			}

			d, err = j.AddFriendDate(p.ID, d)
			if err != nil {
				return err
			}

			log.Info(" ✔ Date added")
			log.Infof("  %s: %s", d.DateExpr, d.Desc) // TODO: improve this output

			return nil
		})
	},
}
