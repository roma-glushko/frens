package main

import (
	"os"
	"time"

	"github.com/roma-glushko/frens/cmd"
	"github.com/urfave/cli/v2"

	"github.com/charmbracelet/log"
	"github.com/muesli/termenv"
	"github.com/roma-glushko/frens/internal/version"
)

func InitLogging(debugLevel bool) {
	log.SetOutput(os.Stdout)
	log.SetPrefix(version.AppName)
	log.SetLevel(log.InfoLevel)
	log.SetReportTimestamp(false)
	log.SetColorProfile(termenv.TrueColor)

	if debugLevel {
		log.SetLevel(log.DebugLevel)
		log.SetReportTimestamp(true)
		log.SetTimeFormat(time.Kitchen)
		log.SetReportCaller(true)
	}
}

const Copyright = `2025-Present, Roma Hlushko & Friends (c)`

func main() {
	cliApp := cli.App{
		Name:                 "frens",
		Usage:                "A friendship management & journaling app. Build friendship that lasts.",
		Version:              version.FullVersion,
		Copyright:            Copyright,
		Suggest:              true,
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "set verbose level",
			},
		},
		Before: func(c *cli.Context) error {
			debugLevel := c.Bool("debug")

			InitLogging(debugLevel)

			return nil
		},
		Commands: []*cli.Command{
			cmd.InitCommand,
			cmd.EditCommand,
			cmd.AddCommands,
			cmd.ListCommands,
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
