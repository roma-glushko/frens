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

package cmd

import (
	"fmt"
	"os"

	"github.com/roma-glushko/frens/internal/config"
	"github.com/roma-glushko/frens/internal/store/file"

	"github.com/roma-glushko/frens/cmd/telegram"

	"github.com/mattn/go-isatty"
	"github.com/roma-glushko/frens/cmd/activity"
	"github.com/roma-glushko/frens/cmd/friend"
	"github.com/roma-glushko/frens/cmd/journal"
	"github.com/roma-glushko/frens/cmd/location"
	"github.com/roma-glushko/frens/cmd/note"
	jctx "github.com/roma-glushko/frens/internal/context"
	"github.com/roma-glushko/frens/internal/log"
	"github.com/roma-glushko/frens/internal/version"
	"github.com/urfave/cli/v2"
)

func InitLogging(verbose bool, quiet bool) {
	level := log.LogLevelStandard

	if verbose {
		level = log.LogLevelVerbose
	}

	if quiet {
		level = log.LogLevelQuiet
	}

	if verbose && quiet {
		level = log.LogLevelStandard

		log.Warn("cannot set both verbose and quiet modes at the same time, ignoring both flags")
	}

	log.SetLevel(level)
}

const Copyright = `2025-Present, Roma Hlushko & Friends (c)`

func NewApp() cli.App {
	return cli.App{
		Name:                 "frens",
		Usage:                "A friendship management & journaling app. Build friendships that last.",
		Version:              version.FullVersion,
		Copyright:            Copyright,
		Suggest:              true,
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "enable verbose output",
			},
			&cli.BoolFlag{
				Name:    "quiet",
				Aliases: []string{"q"},
				Usage:   "enable quiet output",
			},
			&cli.BoolFlag{
				Name:    "interactive",
				Aliases: []string{"i"},
				Value:   isatty.IsTerminal(os.Stdin.Fd()),
				Usage:   "Enable interactive questions and prompts",
			},
			&cli.StringFlag{
				Name:    "journal",
				Aliases: []string{"j"},
				Usage:   "path to the journal directory (default: ~/.config/frens/)",
			},
		},
		Before: func(c *cli.Context) error {
			ctx := c.Context
			debugLevel := c.Bool("debug")
			quietLevel := c.Bool("quiet")

			InitLogging(debugLevel, quietLevel)

			jDir, err := config.Dir(c.String("journal"))
			if err != nil {
				return fmt.Errorf("could not load journal directory from %s: %v", jDir, err)
			}

			log.Debugf(" Using journal directory: %s", jDir)

			appCtx := jctx.AppContext{
				JournalDir: jDir,
				Store:      file.NewTOMLFileStore(jDir),
			}

			c.Context = jctx.WithCtx(ctx, &appCtx)

			return nil
		},
		Commands: []*cli.Command{
			journal.Commands,
			friend.Commands,
			location.Commands,
			note.Commands,
			activity.Commands,
			telegram.Commands,
			ZenCommand,
		},
	}
}
