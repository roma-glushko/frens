package cmd

import (
	"fmt"

	"github.com/roma-glushko/frens/internal/lifedir"
	"github.com/urfave/cli/v2"
)

var InitCommand = &cli.Command{
	Name:    "init",
	Aliases: []string{"i"},
	Usage:   "Init a new life space",
	Flags:   []cli.Flag{},
	Action: func(context *cli.Context) error {
		lifeDir, err := lifedir.DefaultDir()
		if err != nil {
			return err
		}

		err = lifedir.Init(lifeDir)
		if err != nil {
			return err
		}

		fmt.Println("Life space initialized at", lifeDir)

		return nil
	},
}
