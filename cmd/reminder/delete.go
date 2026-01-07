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

package reminder

import (
	"fmt"

	jctx "github.com/roma-glushko/frens/internal/context"
	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/log"
	"github.com/urfave/cli/v2"
)

var DeleteCommand = &cli.Command{
	Name:      "delete",
	Aliases:   []string{"del", "rm"},
	Usage:     "Delete a reminder",
	Args:      true,
	ArgsUsage: "<ID>",
	Action: func(c *cli.Context) error {
		if c.NArg() == 0 {
			return cli.Exit("You must provide a reminder ID", 1)
		}

		reminderID := c.Args().First()

		ctx := c.Context
		appCtx := jctx.FromCtx(ctx)

		return appCtx.Store.Tx(ctx, func(j *journal.Journal) error {
			r, err := j.GetReminder(reminderID)
			if err != nil {
				return fmt.Errorf("reminder not found: %v", err)
			}

			if err := j.RemoveReminders([]friend.Reminder{r}); err != nil {
				return fmt.Errorf("failed to delete reminder: %v", err)
			}

			log.Infof(" ✔ Reminder deleted: %s", reminderID)

			return nil
		})
	},
}
