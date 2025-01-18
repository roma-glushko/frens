package main

import (
	"github.com/roma-glushko/frens/cmd"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	cliApp := cli.App{
		Name:  "frens",
		Usage: "A friendship management & journaling application. Build friendship that lasts",
		Commands: []*cli.Command{
			cmd.InitCommand,
			cmd.AddCommands,
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
