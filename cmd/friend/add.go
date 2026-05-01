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

package friend

import (
	"errors"
	"strings"

	log "github.com/roma-glushko/frens/internal/log"

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
		interactive := c.NArg() == 0
		var info string

		if interactive {
			info = promptFriendInfo("")
		} else {
			info = strings.Join(c.Args().Slice(), " ")
		}

		f, err := parseFriend(c, info)
		if err != nil {
			return err
		}

		ctx := c.Context
		appCtx := jctx.FromCtx(ctx)

		err = appCtx.Store.Tx(ctx, func(j *journal.Journal) error {
			f, err = j.AddFriend(f)
			if err == nil {
				return nil
			}

			if !errors.Is(err, journal.ErrDuplicateFriend) || !interactive {
				return err
			}

			// Interactive mode: let user fix the input
			log.Warnf("Friend with ID %q already exists. Please change the name or set a different $id:", f.ID)

			for {
				info = promptFriendInfo(info)

				f, err = parseFriend(c, info)
				if err != nil {
					return err
				}

				f, err = j.AddFriend(f)
				if err == nil {
					return nil
				}

				if !errors.Is(err, journal.ErrDuplicateFriend) {
					return err
				}

				log.Warnf("Friend with ID %q already exists. Please change the name or set a different $id:", f.ID)
			}
		})
		if err != nil {
			return err
		}

		log.Success("Friend added")
		log.Header("Friend Information")

		return appCtx.Printer.Print(f)
	},
}

func promptFriendInfo(prefill string) string {
	opts := tui.EditorOptions{
		Title:      "Add a new friend information:",
		SyntaxHint: lang.FormatPersonInfo,
	}

	inputForm := tui.NewEditorForm(opts)

	if prefill != "" {
		inputForm.Textarea.SetValue(prefill)
	}

	teaUI := tea.NewProgram(inputForm, tea.WithMouseAllMotion())

	if _, err := teaUI.Run(); err != nil {
		log.Errorf("uh oh: %v", err)
		return ""
	}

	return inputForm.Textarea.Value()
}

func parseFriend(c *cli.Context, info string) (friend.Person, error) {
	var f friend.Person
	var err error

	if info != "" {
		f, err = lang.ExtractPerson(info)

		if err != nil && !errors.Is(err, lang.ErrNoInfo) {
			return friend.Person{}, err
		}
	}

	if id := c.String("id"); id != "" {
		f.ID = id
	}

	if name := c.String("name"); name != "" {
		f.Name = name
	}

	if desc := c.String("desc"); desc != "" {
		f.Desc = desc
	}

	if nicknames := c.StringSlice("nickname"); len(nicknames) > 0 {
		f.Nicknames = nicknames
	}

	if tags := c.StringSlice("tag"); len(tags) > 0 {
		f.Tags = tags
	}

	if locs := c.StringSlice("location"); len(locs) > 0 {
		f.Locations = locs
	}

	if err := f.Validate(); err != nil {
		return friend.Person{}, err
	}

	return f, nil
}
