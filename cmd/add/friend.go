package add

import "github.com/urfave/cli/v2"

var AddFriendCommand = &cli.Command{
	Name:    "friend",
	Aliases: []string{"f"},
	Usage:   "Add a new friend",
	Action: func(context *cli.Context) error {
		// TODO: implement
	},
}
