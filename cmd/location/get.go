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

package location

import (
	"strings"

	jctx "github.com/roma-glushko/frens/internal/context"
	"github.com/roma-glushko/frens/internal/journal"

	"github.com/urfave/cli/v2"
)

var GetCommand = &cli.Command{
	Name:      "get",
	Aliases:   []string{"g", "view", "show"},
	Usage:     "Get and display location information",
	Args:      true,
	ArgsUsage: `<LOCATION_NAME, LOCATION_ALIAS, LOCATION_ID>`,
	Action: func(c *cli.Context) error {
		if c.NArg() < 1 {
			return cli.Exit(
				"You must provide a location name, alias, or ID. Execute `frens location ls` to find out.",
				1,
			)
		}

		lID := strings.Join(c.Args().Slice(), " ")

		ctx := c.Context
		appCtx := jctx.FromCtx(ctx)
		s := appCtx.Store

		return s.Tx(ctx, func(j *journal.Journal) error {
			l, err := j.GetLocation(lID)
			if err != nil {
				return err
			}

			return appCtx.Printer.Print(l)
		})
	},
}
