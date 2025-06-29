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
	"strings"

	"github.com/roma-glushko/frens/internal/life"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/lifedir"
	"github.com/roma-glushko/frens/internal/tui"
	"github.com/urfave/cli/v2"
)

var AddCommand = &cli.Command{
	Name:      "friend",
	Aliases:   []string{"f"},
	Usage:     "Add a new friend",
	Args:      true,
	ArgsUsage: "<NAME>",
	Flags: []cli.Flag{
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
			Aliases: []string{"a", "alias", "nick", "n"},
			Usage:   "Add friend's nicknames",
		},
	},
	Action: func(ctx *cli.Context) error {
		lifeDir, err := lifedir.DefaultDir()
		if err != nil {
			return err
		}

		l, err := lifedir.Load(lifeDir)
		if err != nil {
			return err
		}

		nicknames := ctx.StringSlice("nickname")
		tags := ctx.StringSlice("tag")
		locs := ctx.StringSlice("location")

		var friend friend.Person

		friend.Nicknames = nicknames
		friend.Tags = tags
		friend.Locations = locs

		if ctx.NArg() == 0 {
			// return cli.ShowCommandHelp(context, context.Command.Name)
			teaUI := tea.NewProgram(tui.NewFriendForm(&friend), tea.WithMouseAllMotion())

			if _, err := teaUI.Run(); err != nil {
				log.Error("uh oh", "err", err)
				return err
			}
		} else {
			name := strings.Join(ctx.Args().Slice(), " ")

			friend.Name = name
		}

		if err := friend.Validate(); err != nil {
			return err
		}

		err = lifedir.Update(l, func(l *life.Data) error {
			l.AddFriend(friend)

			return nil
		})
		if err != nil {
			return err
		}

		log.Info(friend.Name + " has been added")

		return nil
	},
}
