package add

import "github.com/urfave/cli/v2"

var AddNoteCommand = &cli.Command{
	Name:    "note",
	Aliases: []string{"n"},
	Usage:   "Add a new note",
	Action: func(context *cli.Context) error {
		// TODO: implement

		return nil
	},
}
