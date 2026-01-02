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

package location

import (
	"fmt"

	jctx "github.com/roma-glushko/frens/internal/context"

	"github.com/roma-glushko/frens/internal/utils"

	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/tui"
	"github.com/urfave/cli/v2"
)

var DeleteCommand = &cli.Command{
	Name:      "delete",
	Aliases:   []string{"del", "rm", "d"},
	Usage:     "Delete a location",
	Args:      true,
	ArgsUsage: `<LOCATION_NAME, LOCATION_NICKNAME, LOCATION_ID> [...]`,
	Description: `Delete locations from your journal by their name, alias, or ID.
	Examples:
		frens friend delete "Nashua"
		frens friend d -f "Utica"
	`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Value:   false,
			Usage:   "Force delete without confirmation",
		},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		if c.NArg() == 0 {
			return cli.Exit("Please provide a location name, alias, or ID to delete.", 1)
		}

		locations := make([]friend.Location, 0, c.NArg())

		appCtx := jctx.FromCtx(ctx)

		return appCtx.Store.Tx(ctx, func(j *journal.Journal) error {
			for _, lID := range c.Args().Slice() {
				l, err := j.GetLocation(lID)
				if err != nil {
					return err
				}

				locations = append(locations, l)
			}

			locWord := utils.P(len(locations), "location", "locations")
			fmt.Printf("\n Found %d %s:\n\n", len(locations), locWord)

			for _, l := range locations {
				fmt.Printf(" • %s \n", l.String())
			}

			// TODO: check if interactive mode
			fmt.Println("\n You're about to permanently delete the " + locWord + ".")
			if !c.Bool("force") && !tui.ConfirmAction(" Are you sure?") {
				fmt.Println("\n ↩ Deletion canceled.")
				return nil
			}

			j.RemoveLocations(locations)

			fmt.Printf("\n ✔ %s deleted.\n", utils.TitleCaser.String(locWord))

			return nil
		})
	},
}
