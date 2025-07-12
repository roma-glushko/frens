// Copyright 2025 Roma Hlushko
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"os"
	"time"

	"github.com/roma-glushko/frens/cmd/journal"
	"github.com/roma-glushko/frens/cmd/location"
	"github.com/roma-glushko/frens/cmd/note"

	"github.com/roma-glushko/frens/cmd/activity"
	"github.com/roma-glushko/frens/cmd/friend"
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
		Usage:                "A friendship management & journaling. Build friendship that lasts.",
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
			journal.Commands,
			friend.Commands,
			location.Commands,
			note.Commands,
			activity.Commands,
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
