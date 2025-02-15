package cmd

import (
	"github.com/roma-glushko/frens/cmd/friend"
	"github.com/urfave/cli/v2"
)

var ListCommands = &cli.Command{
	Name:    "list",
	Aliases: []string{"l"},
	Usage:   "List your friends, activities, locations",
	Subcommands: []*cli.Command{
		friend.ListCommand,
	},
}
