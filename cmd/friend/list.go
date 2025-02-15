package friend

import "github.com/urfave/cli/v2"

var ListCommand = &cli.Command{
	Name:    "friend",
	Aliases: []string{"f"},
	Usage:   "List all friends",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "location",
			Aliases: []string{"l", "loc", "in"},
		},
		&cli.StringFlag{
			Name:    "tag",
			Aliases: []string{"t"},
			Usage:   "Filter by tag",
		},
		&cli.StringFlag{
			Name: "sort",
		},
	},
	Action: func(context *cli.Context) error {
		return nil
	},
}
