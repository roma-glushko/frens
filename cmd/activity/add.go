package activity

import (
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"
)

var AddCommand = &cli.Command{
	Name:    "activity",
	Aliases: []string{"a"},
	Usage:   "Add a new activity",
	Action: func(_ *cli.Context) error {
		log.Info("Adding a new activity..")

		return nil
	},
}
