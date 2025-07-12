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

package friend

import (
	"errors"
	"strings"

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
	Name:    "add",
	Aliases: []string{"a", "new", "create"},
	Usage:   "Add a new friend",
	Args:    true,
	ArgsUsage: `<INFO>
		If no arguments are provided, a textarea will be shown to fill in the details interactively.
		Otherwise, the information will be parsed from the command options.
		
		<INFO> format:
			` + lang.FormatPersonInfo + `

		For example:
			Michael Harry Scott (a.k.a. The World's Best Boss, Mike) :: my Dunder Mifflin boss #office @Scranton $id:mscott
	`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "id",
			Usage: "Friend's unique identifier (used for linking with other data, editing, etc.)",
		},
		&cli.StringFlag{
			Name:    "name",
			Aliases: []string{"n"},
			Usage:   "Friend's name (required if no arguments are provided)",
		},
		&cli.StringFlag{
			Name:    "desc",
			Aliases: []string{"d"},
			Usage:   "Description of the friend (optional, used for additional information)",
		},
		&cli.StringSliceFlag{
			Name:    "tag",
			Aliases: []string{"t"},
			Usage:   "Add tags to the friend",
		},
		&cli.StringSliceFlag{
			Name:    "location",
			Aliases: []string{"l", "loc"},
			Usage:   "Add locations to the friend",
		},
		&cli.StringSliceFlag{
			Name:    "nickname",
			Aliases: []string{"a", "aka", "alias", "nick"},
			Usage:   "Add friend's nicknames (used in search and matching the friend in activities)",
		},
	},
	Action: func(ctx *cli.Context) error {
		journalDir, err := journaldir.DefaultDir()
		if err != nil {
			return err
		}

		j, err := journaldir.Load(journalDir)
		if err != nil {
			return err
		}

		var info string

		if ctx.NArg() == 0 {
			// TODO: also check if we are in the interactive mode
			inputForm := tui.NewInputForm(tui.FormOptions{
				Title:      "Add a new friend information:",
				SyntaxHint: lang.FormatPersonInfo,
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

		if info != "" {
			f, err = lang.ParsePerson(info)

			if err != nil && !errors.Is(err, lang.ErrNoInfo) {
				log.Error("failed to parse friend info", "err", err)
				return err
			}
		}

		// apply CLI flags
		// id := ctx.String("id")
		name := ctx.String("name")
		desc := ctx.String("desc")
		nicknames := ctx.StringSlice("nickname")
		tags := ctx.StringSlice("tag")
		locs := ctx.StringSlice("location")

		// TODO: add support for ID
		//if id != "" {
		//	friend.ID = id
		//}

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

		err = journaldir.Update(j, func(l *journal.Data) error {
			l.AddFriend(f)
			return nil
		})
		if err != nil {
			return err
		}

		log.Info("✅Added new friend: " + f.Name)
		if len(f.Locations) > 0 {
			log.Info("📍 Locations: " + strings.Join(f.Locations, ", "))
		}
		if len(f.Tags) > 0 {
			log.Info("🏷️ Tags: " + strings.Join(f.Tags, ", "))
		}

		return nil
	},
}
