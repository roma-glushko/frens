package main

import (
	"github.com/charmbracelet/log"
	"github.com/muesli/termenv"
	"github.com/roma-glushko/frens/cmd"
	"github.com/roma-glushko/frens/internal/version"
	"github.com/urfave/cli/v2"
	"os"
	"time"
)

func InitLogging(debugLevel bool) {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	log.SetReportTimestamp(false)
	log.SetColorProfile(termenv.TrueColor)

	if debugLevel {
		log.SetLevel(log.DebugLevel)
		log.SetReportTimestamp(true)
		log.SetTimeFormat(time.Kitchen)
		log.SetReportCaller(true)
		log.SetPrefix(version.AppName)
	}
}

func main() {
	cliApp := cli.App{
		Name:                 "frens",
		Usage:                "A friendship management & journaling application. Build friendship that lasts",
		Version:              version.FullVersion,
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
			cmd.AddCommands,
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
