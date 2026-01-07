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

var ListCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls", "l"},
	Usage:   "List reminders",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "state",
			Usage: "Filter by state: pending, fired, acknowledged",
		},
		&cli.StringFlag{
			Name:  "type",
			Usage: "Filter by linked entity type: date, wishlist, activity, note",
		},
		&cli.StringSliceFlag{
			Name:    "tag",
			Aliases: []string{"t"},
			Usage:   "Filter by tags",
		},
	},
	Action: func(c *cli.Context) error {
		state := c.String("state")
		linkedType := c.String("type")
		tags := c.StringSlice("tag")

		ctx := c.Context
		appCtx := jctx.FromCtx(ctx)

		return appCtx.Store.Tx(ctx, func(j *journal.Journal) error {
			q := friend.ListReminderQuery{
				Tags: tags,
			}

			if state != "" {
				q.State = state
			}

			if linkedType != "" {
				q.LinkedEntityType = linkedType
			}

			reminders := j.ListReminders(q)

			if len(reminders) == 0 {
				log.Info("No reminders found")
				return nil
			}

			log.Infof("Found %d reminder(s):\n", len(reminders))

			for _, r := range reminders {
				printReminder(&r)
			}

			return nil
		})
	},
}

func printReminder(r *friend.Reminder) {
	stateIcon := "⏳"

	switch r.State {
	case friend.ReminderStateFired:
		stateIcon = "✅"
	case friend.ReminderStateAcknowledged:
		stateIcon = "✓"
	}

	fmt.Printf("%s [%s] %s:%s\n", stateIcon, r.ID[:8], r.LinkedEntityType, r.LinkedEntityID[:8])
	fmt.Printf("   Trigger: %s", r.TriggerAt.Format("2006-01-02"))

	if r.Recurrence != friend.RecurrenceOnce && r.Recurrence != "" {
		fmt.Printf(" (%s)", r.Recurrence)
	}

	fmt.Println()

	if r.Desc != "" {
		fmt.Printf("   Desc: %s\n", r.Desc)
	}

	if len(r.Tags) > 0 {
		fmt.Printf("   Tags: %v\n", r.Tags)
	}

	fmt.Println()
}
