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
	"strings"
	"time"

	jctx "github.com/roma-glushko/frens/internal/context"
	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/lang"
	"github.com/roma-glushko/frens/internal/log"
	"github.com/urfave/cli/v2"
)

var AddCommand = &cli.Command{
	Name:    "add",
	Aliases: []string{"a", "new"},
	Usage:   "Add a new reminder",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "link",
			Aliases:  []string{"l"},
			Usage:    "Entity to link: date:<id>, wishlist:<id>, activity:<id>, note:<id>",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "schedule",
			Aliases: []string{"s"},
			Usage:   "Schedule expression: 2025-03-15, 1w before, yearly, in 3d, etc.",
		},
		&cli.StringFlag{
			Name:  "desc",
			Usage: "Optional description",
		},
		&cli.StringSliceFlag{
			Name:    "tag",
			Aliases: []string{"t"},
			Usage:   "Tags for routing (can be specified multiple times)",
		},
	},
	Action: func(c *cli.Context) error {
		linkStr := c.String("link")
		scheduleStr := c.String("schedule")
		desc := c.String("desc")
		tags := c.StringSlice("tag")

		// Parse link
		entityType, entityID, err := parseLink(linkStr)
		if err != nil {
			return cli.Exit(err.Error(), 1)
		}

		// Parse schedule
		var schedule *lang.ReminderSchedule
		if scheduleStr != "" {
			schedule, err = lang.ParseReminderSchedule(scheduleStr)
			if err != nil {
				return cli.Exit(fmt.Sprintf("Invalid schedule: %v", err), 1)
			}
		}

		ctx := c.Context
		appCtx := jctx.FromCtx(ctx)

		return appCtx.Store.Tx(ctx, func(j *journal.Journal) error {
			// Resolve base date from linked entity
			baseDate, friendID, err := resolveBaseDate(j, entityType, entityID)
			if err != nil {
				return cli.Exit(fmt.Sprintf("Failed to resolve linked entity: %v", err), 1)
			}

			now := time.Now()

			// Compute trigger date
			var triggerAt time.Time
			if schedule != nil {
				triggerAt = schedule.ComputeTriggerDate(baseDate, now)
			} else {
				triggerAt = baseDate
			}

			// Build reminder
			r := friend.Reminder{
				LinkedEntityType: entityType,
				LinkedEntityID:   entityID,
				FriendID:         friendID,
				TriggerAt:        triggerAt,
				ScheduleExpr:     scheduleStr,
				State:            friend.ReminderStatePending,
				Desc:             desc,
				Tags:             tags,
			}

			// Set recurrence and offset from schedule
			if schedule != nil {
				r.Recurrence = schedule.Recurrence
				r.OffsetDirection = schedule.OffsetDirection
				r.OffsetDuration = schedule.ToDuration()
			} else {
				r.Recurrence = friend.RecurrenceOnce
			}

			// Add reminder
			r, err = j.AddReminder(r)
			if err != nil {
				return fmt.Errorf("failed to add reminder: %v", err)
			}

			log.Infof(" ✔ Reminder added")
			log.Infof("   Linked to: %s:%s", entityType, entityID)
			log.Infof("   Trigger: %s", triggerAt.Format("2006-01-02"))

			if r.Recurrence != friend.RecurrenceOnce {
				log.Infof("   Recurrence: %s", r.Recurrence)
			}

			return nil
		})
	},
}

func parseLink(link string) (friend.LinkedEntityType, string, error) {
	parts := strings.SplitN(link, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid link format, expected: type:id (e.g., date:abc123)")
	}

	entityType := friend.LinkedEntityType(parts[0])
	entityID := parts[1]

	switch entityType {
	case friend.LinkedEntityDate, friend.LinkedEntityWishlist,
		friend.LinkedEntityActivity, friend.LinkedEntityNote:
		return entityType, entityID, nil
	default:
		return "", "", fmt.Errorf("unsupported entity type: %s (supported: date, wishlist, activity, note)", parts[0])
	}
}

func resolveBaseDate(j *journal.Journal, entityType friend.LinkedEntityType, entityID string) (time.Time, string, error) {
	switch entityType {
	case friend.LinkedEntityDate:
		date, err := j.GetFriendDate(entityID)
		if err != nil {
			return time.Time{}, "", err
		}

		// Parse the date expression
		parsedDate := lang.ExtractDate(date.DateExpr)

		return parsedDate, date.Person, nil

	case friend.LinkedEntityWishlist:
		item, err := j.GetFriendWishlistItem(entityID)
		if err != nil {
			return time.Time{}, "", err
		}

		// Wishlist items don't have a natural date, use now
		return time.Now(), item.Person, nil

	case friend.LinkedEntityActivity:
		event, err := j.GetEvent(friend.EventTypeActivity, entityID)
		if err != nil {
			return time.Time{}, "", err
		}

		var friendID string
		if len(event.FriendIDs) > 0 {
			friendID = event.FriendIDs[0]
		}

		return event.Date, friendID, nil

	case friend.LinkedEntityNote:
		event, err := j.GetEvent(friend.EventTypeNote, entityID)
		if err != nil {
			return time.Time{}, "", err
		}

		var friendID string
		if len(event.FriendIDs) > 0 {
			friendID = event.FriendIDs[0]
		}

		return event.Date, friendID, nil

	default:
		return time.Time{}, "", fmt.Errorf("unknown entity type: %s", entityType)
	}
}
