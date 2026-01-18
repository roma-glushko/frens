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
	"os"

	jctx "github.com/roma-glushko/frens/internal/context"
	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/lang"
	"github.com/roma-glushko/frens/internal/log"
	"github.com/urfave/cli/v2"
)

var ExportCommand = &cli.Command{
	Name:      "export",
	Aliases:   []string{"exp"},
	Usage:     "Export journal data to FrenTXT format",
	UsageText: "frens journal export [OPTIONS] [FILE_PATH]",
	Args:      true,
	ArgsUsage: `[FILE_PATH]
		Path to write the FrenTXT file. If not provided, outputs to stdout.
	`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "friends",
			Aliases: []string{"f"},
			Usage:   "Export friends",
		},
		&cli.BoolFlag{
			Name:    "locations",
			Aliases: []string{"l"},
			Usage:   "Export locations",
		},
		&cli.BoolFlag{
			Name:    "notes",
			Aliases: []string{"n"},
			Usage:   "Export notes",
		},
		&cli.BoolFlag{
			Name:    "activities",
			Aliases: []string{"a"},
			Usage:   "Export activities",
		},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context
		appCtx := jctx.FromCtx(ctx)

		// Determine what to export
		exportFriends := c.Bool("friends")
		exportLocations := c.Bool("locations")
		exportNotes := c.Bool("notes")
		exportActivities := c.Bool("activities")

		// If no filters specified, export everything
		exportAll := !exportFriends && !exportLocations && !exportNotes && !exportActivities
		if exportAll {
			exportFriends = true
			exportLocations = true
			exportNotes = true
			exportActivities = true
		}

		filePath := c.Args().First()

		var output string

		err := appCtx.Store.Tx(ctx, func(j *journal.Journal) error {
			data := &lang.FrenTXTImport{}

			if exportFriends {
				data.Friends = make([]friend.Person, 0, len(j.Friends))
				for _, f := range j.Friends {
					data.Friends = append(data.Friends, *f)
				}
			}

			if exportLocations {
				data.Locations = make([]friend.Location, 0, len(j.Locations))
				for _, l := range j.Locations {
					data.Locations = append(data.Locations, *l)
				}
			}

			if exportNotes {
				data.Notes = make([]friend.Event, 0, len(j.Notes))
				for _, n := range j.Notes {
					data.Notes = append(data.Notes, *n)
				}
			}

			if exportActivities {
				data.Activities = make([]friend.Event, 0, len(j.Activities))
				for _, a := range j.Activities {
					data.Activities = append(data.Activities, *a)
				}
			}

			output = lang.RenderFrenTXT(data)

			return nil
		})

		if err != nil {
			return fmt.Errorf("failed to read journal: %w", err)
		}

		if filePath == "" || filePath == "-" {
			fmt.Print(output)
		} else {
			if err := os.WriteFile(filePath, []byte(output), 0o644); err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}

			log.Successf("Exported to %s", filePath)
		}

		return nil
	},
}
