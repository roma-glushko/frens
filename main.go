package main

import (
	"github.com/roma-glushko/frens/cmd"
	"github.com/roma-glushko/frens/cmd/add"
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
			add.Commands,
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
