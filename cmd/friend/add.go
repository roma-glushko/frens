package friend

import (
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"
)

var AddCommand = &cli.Command{
	Name:    "friend",
	Aliases: []string{"f"},
	Usage:   "Add a new friend",
	Action: func(context *cli.Context) error {
		log.Info("Adding a new friend..")

		return nil
	},
}
