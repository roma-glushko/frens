package friend

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/roma-glushko/frens/cmd/friend/tui"
	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/lifedir"
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
			Aliases: []string{"a", "alias"},
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

		if context.NArg() == 0 {
			// return cli.ShowCommandHelp(context, context.Command.Name)
			if _, err := tea.NewProgram(tui.NewFriendModel()).Run(); err != nil {
				log.Error("could not start tui", "err", err)
				return err
			}
		}

		name := strings.Join(context.Args().Slice(), " ")
		nicknames := context.StringSlice("nickname")
		tags := context.StringSlice("tag")
		locs := context.StringSlice("location")

		life.Friends = append(life.Friends, friend.Friend{
			Name:      name,
			Nicknames: nicknames,
			Tags:      tags,
			Locations: locs,
		})

		err = lifedir.Save(lifeDir, life)
		if err != nil {
			return err
		}

		log.Info("New friend added")

		return nil
	},
}
