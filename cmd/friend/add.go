package friend

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/lifedir"
	"github.com/roma-glushko/frens/internal/tui"
	"github.com/urfave/cli/v2"
	"strings"
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
			Usage:   "Add friend's nickname",
		},
	},
	Action: func(context *cli.Context) error {
		lifeDir, err := lifedir.DefaultDir()
		if err != nil {
			return err
		}

		life, err := lifedir.Load(lifeDir)
		if err != nil {
			return err
		}

		nicknames := context.StringSlice("nickname")
		tags := context.StringSlice("tag")
		locs := context.StringSlice("location")

		var friend friend.Friend

		friend.Nicknames = nicknames
		friend.Tags = tags
		friend.Locations = locs

		if context.NArg() == 0 {
			// return cli.ShowCommandHelp(context, context.Command.Name)
			if _, err := tea.NewProgram(tui.NewFriendForm(friend)).Run(); err != nil {
				log.Error("uh oh", "err", err)
				return err
			}
		} else {
			name := strings.Join(context.Args().Slice(), " ")

			friend.Name = name
		}

		if err := friend.Validate(); err != nil {
			return err
		}

		life.AddFriend(friend)

		if err = lifedir.Save(lifeDir, life); err != nil {
			return err
		}

		log.Info("New friend added")

		return nil
	},
}
