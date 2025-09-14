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
	"fmt"
	"strings"

	log "github.com/roma-glushko/frens/internal/log"
	"github.com/roma-glushko/frens/internal/log/formatter"

	jctx "github.com/roma-glushko/frens/internal/context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/lang"
	"github.com/roma-glushko/frens/internal/tui"

	"github.com/roma-glushko/frens/internal/friend"
	"github.com/urfave/cli/v2"
)

var AddCommand = &cli.Command{
	Name:      "add",
	Aliases:   []string{"a", "new", "create"},
	Usage:     "Add a new friend",
	UsageText: "frens friend add [OPTIONS] [INFO]",
	Args:      true,
	ArgsUsage: `<INFO>
		If no arguments are provided, a textarea will be shown to fill in the details interactively.
		Otherwise, the information will be parsed from the command options.
		
		<INFO> format:
			` + lang.FormatPersonInfo + `

		For example:
			Michael Harry Scott (a.k.a. The World's Best Boss, Mike) :: my Dunder Mifflin boss #office @Scranton $id:mscott
			Jim Halpert (aka Jimbo, Big Tuna) :: my best friend and prankster extraordinaire #office @Scranton
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
			Aliases: []string{"a", "aka", "nick"},
			Usage:   "Add friend's nicknames (used in search and matching the friend in activities)",
		},
	},
	Action: func(c *cli.Context) error {
		var info string

		if c.NArg() == 0 {
			// TODO: also check if we are in the interactive mode
			inputForm := tui.NewEditorForm(tui.EditorOptions{
				Title:      "Add a new friend information:",
				SyntaxHint: lang.FormatPersonInfo,
			})
			teaUI := tea.NewProgram(inputForm, tea.WithMouseAllMotion())

			if _, err := teaUI.Run(); err != nil {
				log.Errorf("uh oh: %v", err)
				return err
			}

			info = inputForm.Textarea.Value()
		} else {
			info = strings.Join(c.Args().Slice(), " ")
		}

		var f friend.Person
		var err error

		if info != "" {
			f, err = lang.ExtractPerson(info)

			if err != nil && !errors.Is(err, lang.ErrNoInfo) {
				log.Errorf("failed to parse friend info: %v", err)
				return err
			}
		}

		// apply CLI flags
		id := c.String("id")
		name := c.String("name")
		desc := c.String("desc")
		nicknames := c.StringSlice("nickname")
		tags := c.StringSlice("tag")
		locs := c.StringSlice("location")

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

		ctx := c.Context
		jctx := jctx.FromCtx(ctx)
		s := jctx.Store

		err = s.Tx(ctx, func(j *journal.Journal) error {
			j.AddFriend(f)
			return nil
		})
		if err != nil {
			return err
		}

		log.Info(" ✔ Friend added\n")
		log.Info("==> Friend Information\n")
		// log.PrintEntity(f)

		fmtr := formatter.PersonTextFormatter{}

		o, _ := fmtr.FormatSingle(f)
		fmt.Println(o)

		return nil
	},
}
