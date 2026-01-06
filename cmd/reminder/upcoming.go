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
	"time"

	jctx "github.com/roma-glushko/frens/internal/context"
	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/log"
	"github.com/urfave/cli/v2"
)

var UpcomingCommand = &cli.Command{
	Name:    "upcoming",
	Aliases: []string{"up", "soon"},
	Usage:   "Show upcoming reminders",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    "days",
			Aliases: []string{"d"},
			Value:   30,
			Usage:   "Number of days to look ahead",
		},
	},
	Action: func(c *cli.Context) error {
		days := c.Int("days")
		now := time.Now()

		ctx := c.Context
		appCtx := jctx.FromCtx(ctx)

		return appCtx.Store.Tx(ctx, func(j *journal.Journal) error {
			reminders := j.GetUpcomingReminders(now, days)

			if len(reminders) == 0 {
				log.Infof("No upcoming reminders in the next %d days", days)
				return nil
			}

			log.Infof("Upcoming reminders (next %d days):\n", days)

			// Group by days until
			groups := groupByDaysUntil(reminders, now)

			for _, group := range groups {
				fmt.Printf("\nIn %s:\n", formatDaysUntil(group.DaysUntil))

				for _, r := range group.Reminders {
					printUpcomingReminder(&r, j)
				}
			}

			return nil
		})
	},
}

type reminderGroup struct {
	DaysUntil int
	Reminders []friend.Reminder
}

func groupByDaysUntil(reminders []friend.Reminder, now time.Time) []reminderGroup {
	groups := make(map[int][]friend.Reminder)

	for _, r := range reminders {
		daysUntil := int(r.TriggerAt.Sub(now).Hours() / 24)
		if daysUntil < 0 {
			daysUntil = 0
		}

		groups[daysUntil] = append(groups[daysUntil], r)
	}

	// Convert to sorted slice
	result := make([]reminderGroup, 0, len(groups))
	for d, rs := range groups {
		result = append(result, reminderGroup{DaysUntil: d, Reminders: rs})
	}

	// Sort by days until
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].DaysUntil > result[j].DaysUntil {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result
}

func formatDaysUntil(days int) string {
	switch days {
	case 0:
		return "Today"
	case 1:
		return "Tomorrow"
	default:
		return fmt.Sprintf("%d days", days)
	}
}

func printUpcomingReminder(r *friend.Reminder, j *journal.Journal) {
	desc := r.Desc
	if desc == "" {
		// Try to get description from linked entity
		desc = getLinkedEntityDesc(r, j)
	}

	recurrence := ""
	if r.Recurrence != friend.RecurrenceOnce && r.Recurrence != "" {
		recurrence = fmt.Sprintf(" (%s)", r.Recurrence)
	}

	tags := ""
	if len(r.Tags) > 0 {
		tags = fmt.Sprintf(" #%s", r.Tags[0])
		for _, t := range r.Tags[1:] {
			tags += fmt.Sprintf(" #%s", t)
		}
	}

	fmt.Printf("  - %s%s%s\n", desc, recurrence, tags)
}

func getLinkedEntityDesc(r *friend.Reminder, j *journal.Journal) string {
	switch r.LinkedEntityType {
	case friend.LinkedEntityDate:
		date, err := j.GetFriendDate(r.LinkedEntityID)
		if err == nil {
			if date.Person != "" {
				f, err := j.GetFriend(date.Person)
				if err == nil {
					if date.Desc != "" {
						return fmt.Sprintf("%s - %s", f.Name, date.Desc)
					}

					return fmt.Sprintf("%s's date (%s)", f.Name, date.DateExpr)
				}
			}

			return date.Desc
		}

	case friend.LinkedEntityWishlist:
		item, err := j.GetFriendWishlistItem(r.LinkedEntityID)
		if err == nil {
			return item.Desc
		}

	case friend.LinkedEntityActivity, friend.LinkedEntityNote:
		eventType := friend.EventTypeActivity
		if r.LinkedEntityType == friend.LinkedEntityNote {
			eventType = friend.EventTypeNote
		}

		event, err := j.GetEvent(eventType, r.LinkedEntityID)
		if err == nil {
			return event.Desc
		}
	}

	return fmt.Sprintf("%s:%s", r.LinkedEntityType, r.LinkedEntityID[:8])
}
