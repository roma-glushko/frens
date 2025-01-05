package add

import (
	"github.com/urfave/cli/v2"
)

var Commands = &cli.Command{
	Name:    "add",
	Aliases: []string{"a"},
	Usage:   "Add a new friend, location, activity, etc.",
	Subcommands: []*cli.Command{
		AddFriendCommand,
		AddNoteCommand,
	},
}
