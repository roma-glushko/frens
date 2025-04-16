package cmd

import (
	"fmt"
)

var HelloCommand = &cli.Command{
	Name:    "hello",
	Aliases: []string{"h"},
	Usage:   "A sample CLI command",
	Flags:   []cli.Flag{},
	Action: func(_ *cli.Context) error {
		fmt.Println("Hey")

		return nil
	},
}
