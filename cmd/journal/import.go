// Copyright 2026 Roma Hlushko
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

package journal

import (
	"fmt"
	"io"
	"os"

	jctx "github.com/roma-glushko/frens/internal/context"
	"github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/lang"
	"github.com/roma-glushko/frens/internal/log"
	"github.com/urfave/cli/v2"
)

var ImportCommand = &cli.Command{
	Name:      "import",
	Aliases:   []string{"imp"},
	Usage:     "Import journal data from a FrenTXT file or stdin",
	UsageText: "frens journal import [OPTIONS] [FILE_PATH | -]",
	Args:      true,
	ArgsUsage: `[FILE_PATH | -]
		Path to the FrenTXT file to import, or "-" to read from stdin.
		If no argument is provided, reads from stdin.

		FrenTXT format uses section markers to denote different entity types:
		- /f  : Friend
		- /l  : Location
		- /n  : Note
		- /act: Activity

		Example:
			/f
			Alice (aka Al, Ally) :: Close friend from college.
			#college, #bestie @Boston $id:ALICE

			/l
			Paris, France (aka City of Light) :: A city I visited in 2022.
			#travel $id:PARIS

			/n
			2023-08-15 :: Had dinner with Alice in Paris.
			#catchup @PARIS

			/act
			Yesterday :: Jogged around the Seine.
			#exercise @Paris
	`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "dry-run",
			Aliases: []string{"n"},
			Usage:   "Preview what would be imported without making changes",
		},
	},
	Action: func(c *cli.Context) error {
		dryRun := c.Bool("dry-run")

		var content []byte
		var err error

		filePath := c.Args().First()

		if filePath == "" || filePath == "-" {
			content, err = io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("failed to read from stdin: %w", err)
			}
		} else {
			content, err = os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}
		}

		importData, err := lang.ParseFrenTXT(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse FrenTXT: %w", err)
		}

		friendCount := len(importData.Friends)
		locationCount := len(importData.Locations)
		noteCount := len(importData.Notes)
		activityCount := len(importData.Activities)

		totalCount := friendCount + locationCount + noteCount + activityCount

		if totalCount == 0 {
			log.Info("No data found to import.\n")
			return nil
		}

		if dryRun {
			log.Header("Dry Run - Would Import")
		} else {
			log.Header("Importing")
		}

		if friendCount > 0 {
			log.Bulletf("%d friend(s)", friendCount)
			for _, f := range importData.Friends {
				log.Infof("    - %s\n", f.Name)
			}
		}

		if locationCount > 0 {
			log.Bulletf("%d location(s)", locationCount)
			for _, l := range importData.Locations {
				log.Infof("    - %s\n", l.Name)
			}
		}

		if noteCount > 0 {
			log.Bulletf("%d note(s)", noteCount)
		}

		if activityCount > 0 {
			log.Bulletf("%d activity(ies)", activityCount)
		}

		if dryRun {
			log.Info("\nNo changes made (dry run).\n")
			return nil
		}

		ctx := c.Context
		appCtx := jctx.FromCtx(ctx)

		err = appCtx.Store.Tx(ctx, func(j *journal.Journal) error {
			for _, f := range importData.Friends {
				j.AddFriend(f)
			}

			for _, l := range importData.Locations {
				j.AddLocation(l)
			}

			for _, n := range importData.Notes {
				if _, err := j.AddEvent(n); err != nil {
					log.Warnf("Failed to add note: %v\n", err)
				}
			}

			for _, a := range importData.Activities {
				if _, err := j.AddEvent(a); err != nil {
					log.Warnf("Failed to add activity: %v\n", err)
				}
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("failed to save import: %w", err)
		}

		log.Successf("Imported %d item(s)", totalCount)

		return nil
	},
}
