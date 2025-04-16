package note

import "github.com/urfave/cli/v2"

var AddCommand = &cli.Command{
	Name:    "note",
	Aliases: []string{"n"},
	Usage:   "Add a new note",
	Action: func(_ *cli.Context) error {
		// TODO: implement

		return nil
	},
}
