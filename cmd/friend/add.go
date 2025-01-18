package friend

import (
	"strings"

	"github.com/charmbracelet/log"
	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/lifedir"
	"github.com/urfave/cli/v2"
)

var AddCommand = &cli.Command{
	Name:      "friend",
	Aliases:   []string{"f"},
	Usage:     "Add a new friend",
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
