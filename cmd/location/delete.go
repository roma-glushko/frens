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
	"github.com/roma-glushko/frens/internal/journaldir"
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
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args().Slice()) == 0 {
			return cli.Exit("Please provide a location name, alias, or ID to delete.", 1)
		}

		locations := make([]friend.Location, 0, len(ctx.Args().Slice()))

		jctx := jctx.FromCtx(ctx.Context)
		jr := jctx.Journal

		for _, lID := range ctx.Args().Slice() {
			l, err := jr.GetLocation(lID)
			if err != nil {
				return err
			}

			locations = append(locations, *l)
		}

		locWord := utils.P(len(locations), "location", "locations")
		fmt.Printf("🔍 Found %d %s:\n", len(locations), locWord)

		for _, l := range locations {
			fmt.Printf("   • %s \n", l.String())
		}

		// TODO: check if interactive mode
		fmt.Println("\n⚠️  You're about to permanently delete the " + locWord + ".")
		if !ctx.Bool("force") && !tui.ConfirmAction("Are you sure?") {
			fmt.Println("\n↩️  Deletion canceled.")
			return nil
		}

		err := journaldir.Update(jr, func(j *journal.Journal) error {
			j.RemoveLocations(locations)
			return nil
		})
		if err != nil {
			return err
		}

		fmt.Printf("\n🗑️  %s deleted.\n", utils.TitleCaser.String(locWord))

		return nil
	},
}
