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
	"time"

	"github.com/roma-glushko/frens/cmd/telegram"

	"github.com/roma-glushko/frens/internal/journaldir"

	"github.com/charmbracelet/log"
	"github.com/mattn/go-isatty"
	"github.com/muesli/termenv"
	"github.com/roma-glushko/frens/cmd/activity"
	"github.com/roma-glushko/frens/cmd/friend"
	"github.com/roma-glushko/frens/cmd/journal"
	"github.com/roma-glushko/frens/cmd/location"
	"github.com/roma-glushko/frens/cmd/note"
	jctx "github.com/roma-glushko/frens/internal/context"
	jrnal "github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/version"
	"github.com/urfave/cli/v2"
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
				Usage:   "set verbose level",
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
		Before: func(ctx *cli.Context) error {
			debugLevel := ctx.Bool("debug")

			InitLogging(debugLevel)

			jDir, err := journaldir.Dir(ctx.String("journal"))
			if err != nil {
				return fmt.Errorf("could not load journal directory from %s: %v", jDir, err)
			}

			jCtx := jctx.AppContext{
				JournalDir: jDir,
			}

			ctx.Context = jctx.WithCtx(ctx.Context, &jCtx)

			if journaldir.Exists(jDir) {
				// load only if the journal directory exists (it may not if this is the first run or a new journal path)
				jr, err := journaldir.Load(jDir)
				if err != nil {
					return fmt.Errorf("failed to load journal from %s: %w", jDir, err)
				}

				ctx.Context = jrnal.WithCtx(ctx.Context, jr)
			}

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
