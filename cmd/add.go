package cmd

import (
	"github.com/roma-glushko/frens/cmd/activity"
	"github.com/roma-glushko/frens/cmd/friend"
	"github.com/roma-glushko/frens/cmd/location"
	"github.com/roma-glushko/frens/cmd/note"
	"github.com/urfave/cli/v2"
)

var AddCommands = &cli.Command{
	Name:    "add",
	Aliases: []string{"a"},
	Usage:   "Add a new friend, location, activity, etc.",
	Subcommands: []*cli.Command{
		friend.AddCommand,
		location.AddCommand,
		note.AddCommand,
		activity.AddCommand,
	},
}
