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

package dates

import (
	"errors"
	"fmt"
	"strings"

	jctx "github.com/roma-glushko/frens/internal/context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/lang"
	"github.com/roma-glushko/frens/internal/tui"

	"github.com/charmbracelet/log"
	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/journaldir"
	"github.com/urfave/cli/v2"
)

var AddCommand = &cli.Command{
	Name:      "add",
	Aliases:   []string{"a", "new", "create"},
	Usage:     "Add a new date to a friend",
	UsageText: "frens friend dates add [OPTIONS] [INFO]",
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
			Name:    "label",
			Aliases: []string{"l"},
			Usage:   "Label for the date (e.g., birthday, anniversary)",
		},
		&cli.StringFlag{
			Name:    "date",
			Aliases: []string{"d"},
			Usage:   "Date in a free-form format (e.g., 'May 13th', '2009-9-09')",
		},
		&cli.StringFlag{
			Name:    "calendar",
			Aliases: []string{"cal"},
			Usage:   "Calendar type to use for the date (e.g., gregorian, hebrew)",
		},
	},
	Action: func(ctx *cli.Context) error {
		var info string

		if ctx.NArg() == 0 {
			// TODO: also check if we are in the interactive mode
			inputForm := tui.NewEditorForm(tui.EditorOptions{
				Title:      "Add a new friend date information:",
				SyntaxHint: lang.FormatDateInfo,
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

		var f friend.Person
		var err error

		if info != "" {
			f, err = lang.ExtractPerson(info)

			if err != nil && !errors.Is(err, lang.ErrNoInfo) {
				log.Error("failed to parse friend info", "err", err)
				return err
			}
		}

		// apply CLI flags
		id := ctx.String("id")
		name := ctx.String("name")
		desc := ctx.String("desc")
		nicknames := ctx.StringSlice("nickname")
		tags := ctx.StringSlice("tag")
		locs := ctx.StringSlice("location")

		if id != "" {
			f.ID = id
		}

		if name != "" {
			f.Name = name
		}

		if desc != "" {
			f.Desc = desc
		}

		if len(nicknames) > 0 {
			f.Nicknames = nicknames
		}

		if len(tags) > 0 {
			f.Tags = tags
		}

		if len(locs) > 0 {
			f.Locations = locs
		}

		if err := f.Validate(); err != nil {
			return err
		}

		jctx := jctx.FromCtx(ctx.Context)
		jr := jctx.Journal

		err = journaldir.Update(jr, func(l *journal.Journal) error {
			l.AddFriend(f)
			return nil
		})
		if err != nil {
			return err
		}

		fmt.Println("✅ Added new friend: " + f.String())
		if len(f.Locations) > 0 {
			fmt.Println("📍 Locations: " + strings.Join(f.Locations, ", "))
		}
		if len(f.Tags) > 0 {
			fmt.Println("🏷️ Tags: " + strings.Join(f.Tags, ", "))
		}

		return nil
	},
}
