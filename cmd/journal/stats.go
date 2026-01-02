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

package journal

import (
	"fmt"

	"github.com/roma-glushko/frens/internal/journal"

	jctx "github.com/roma-glushko/frens/internal/context"

	"github.com/urfave/cli/v2"
)

var StatsCommand = &cli.Command{
	Name:  "stats",
	Usage: "Show journal statistics",
	Action: func(c *cli.Context) error {
		ctx := c.Context
		appCtx := jctx.FromCtx(ctx)

		return appCtx.Store.Tx(ctx, func(j *journal.Journal) error {
			stats := j.Stats()

			fmt.Println("Journal Statistics:")
			fmt.Printf("  • Friends: %d\n", stats.Friends)
			fmt.Printf("  • Locations: %d\n", stats.Locations)
			fmt.Printf("  • Activities: %d\n", stats.Activities)
			fmt.Printf("  • Notes: %d\n", stats.Notes)

			return nil
		})
	},
}
