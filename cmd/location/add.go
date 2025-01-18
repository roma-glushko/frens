package location

import (
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"
)

var AddCommand = &cli.Command{
	Name:    "location",
	Aliases: []string{"l"},
	Usage:   "Add a new location",
	Action: func(context *cli.Context) error {
		log.Info("Adding a new location..")

		return nil
	},
}
