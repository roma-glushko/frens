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

package telegram

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	jctx "github.com/roma-glushko/frens/internal/context"

	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/lang"
	"github.com/roma-glushko/frens/internal/version"
	"gopkg.in/telebot.v4/middleware"

	"github.com/pkg/errors"

	"github.com/urfave/cli/v2"
	tele "gopkg.in/telebot.v4"
)

const helpMessage = `Welcome to the Frens Bot! Here is what you can do:
/stats - Journal statistics.

/addfriend - Add a new friend.
/addlocation - Add a new location.
/addnote - Add a new note.
/addactivity - Add a new activity.

/listfrens - List my friends.
/listlocs - List my locations.
/listnotes - List my notes.
/listactivities - List my activities.

/version - Show the current version of frens.
`

var BotCommand = &cli.Command{
	Name:  "bot",
	Usage: "Start a Telegram bot",
	Flags: []cli.Flag{
		&cli.Int64SliceFlag{
			Name:    "user_id",
			Aliases: []string{"u"},
			Usage:   "List of user IDs that bot will respond to (all others will be ignored).",
		},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context
		userList := c.Int64Slice("user_id")
		token := os.Getenv("TOKEN")

		if token == "" {
			return errors.New("TOKEN environment variable is not set")
		}

		pref := tele.Settings{
			Token:  token,
			Poller: &tele.LongPoller{Timeout: 10 * time.Second},
		}

		bot, err := tele.NewBot(pref)
		if err != nil {
			return fmt.Errorf("failed to start telegram bot: %w", err)
		}

		startedAt := time.Now()
		appCtx := jctx.FromCtx(ctx)
		s := appCtx.Store

		if len(userList) > 0 {
			bot.Use(middleware.Whitelist(userList...))
		} else {
			fmt.Println("Warning: No user IDs specified, bot will respond to all users.")
		}

		bot.Handle("/help", func(c tele.Context) error {
			return c.Send(helpMessage)
		})

		bot.Handle("/start", func(c tele.Context) error {
			return c.Send(helpMessage)
		})

		bot.Handle("/version", func(c tele.Context) error {
			var sb strings.Builder

			sb.WriteString("Frens Version: " + version.Version + "\n")
			sb.WriteString("Built At: " + version.BuildDate + "\n")
			sb.WriteString("Commit: " + version.GitCommit + "\n")
			sb.WriteString("Uptime: " + time.Since(startedAt).String() + "\n")

			return c.Send(sb.String())
		})

		bot.Handle("/stats", func(c tele.Context) error {
			ctx := context.Background()

			return s.Tx(ctx, func(jr *journal.Journal) error {
				stats := jr.Stats()

				var sb strings.Builder

				sb.WriteString("Frens Stats:\n")
				sb.WriteString(fmt.Sprintf("Friends: %d\n", stats.Friends))
				sb.WriteString(fmt.Sprintf("Locations: %d\n", stats.Locations))
				sb.WriteString(fmt.Sprintf("Notes: %d\n", stats.Notes))
				sb.WriteString(fmt.Sprintf("Activities: %d\n", stats.Activities))

				return c.Send(sb.String())
			})
		})

		bot.Handle("/listfrens", func(c tele.Context) error {
			payload := c.Message().Payload

			q, err := lang.ExtractPersonQuery(payload)
			if err != nil {
				return c.Send(fmt.Sprintf("Failed to parse query: %v", err))
			}

			ctx := context.Background()

			return s.Tx(ctx, func(jr *journal.Journal) error {
				friends := jr.ListFriends(q)

				if len(friends) == 0 {
					return c.Send("No friends found matching your query.")
				}

				var sb strings.Builder

				sb.WriteString("Here are your friends:\n")

				for _, f := range friends {
					sb.WriteString(fmt.Sprintf("👤 %s\n", f.String()))

					if len(f.Tags) > 0 {
						sb.WriteString(fmt.Sprintf("  Tags: %s\n", strings.Join(f.Tags, ", ")))
					}

					if len(f.Locations) > 0 {
						sb.WriteString(
							fmt.Sprintf("  Locations: %s\n", strings.Join(f.Locations, ", ")),
						)
					}

					if f.Desc != "" {
						sb.WriteString(fmt.Sprintf("  Description: %s\n", f.Desc))
					}

					sb.WriteString("\n")
				}

				return c.Send(sb.String())
			})
		})

		bot.Handle("/listlocs", func(c tele.Context) error {
			payload := c.Message().Payload

			q, err := lang.ExtractLocationQuery(payload)
			if err != nil {
				return c.Send(fmt.Sprintf("Failed to parse query: %v", err))
			}

			ctx := context.Background()

			return s.Tx(ctx, func(jr *journal.Journal) error {
				locs := jr.ListLocations(q)

				if len(locs) == 0 {
					return c.Send("No locations found matching your query.")
				}

				var sb strings.Builder

				sb.WriteString("Here are your locations:\n")

				for _, l := range locs {
					sb.WriteString(fmt.Sprintf("📍 %s\n", l.String()))

					if len(l.Tags) > 0 {
						sb.WriteString(fmt.Sprintf("  Tags: %s\n", strings.Join(l.Tags, ", ")))
					}

					sb.WriteString("\n")
				}

				return c.Send(sb.String())
			})
		})

		bot.Handle("/listnotes", func(c tele.Context) error {
			payload := c.Message().Payload

			q, err := lang.ExtractEventQuery(payload)
			if err != nil {
				return c.Send(fmt.Sprintf("Failed to parse query: %v", err))
			}

			q.Type = friend.EventTypeNote

			ctx := context.Background()

			return s.Tx(ctx, func(jr *journal.Journal) error {
				notes, err := jr.ListEvents(q)
				if err != nil {
					return c.Send(fmt.Sprintf("Failed to list notes: %v", err))
				}

				if len(notes) == 0 {
					return c.Send("No notes found matching your query.")
				}

				var sb strings.Builder

				sb.WriteString("Here are your notes:\n")

				for _, nt := range notes {
					sb.WriteString(nt.Desc + "\n")

					if len(nt.Tags) > 0 {
						sb.WriteString(fmt.Sprintf("  Tags: %s\n", strings.Join(nt.Tags, ", ")))
					}

					sb.WriteString("\n")
				}

				return c.Send(sb.String())
			})
		})

		bot.Handle("/listactivities", func(c tele.Context) error {
			payload := c.Message().Payload

			q, err := lang.ExtractEventQuery(payload)
			if err != nil {
				return c.Send(fmt.Sprintf("Failed to parse query: %v", err))
			}

			q.Type = friend.EventTypeActivity

			ctx := context.Background()

			return s.Tx(ctx, func(jr *journal.Journal) error {
				activities, err := jr.ListEvents(q)
				if err != nil {
					return c.Send(fmt.Sprintf("Failed to list activities: %v", err))
				}

				if len(activities) == 0 {
					return c.Send("No activities found matching your query.")
				}

				var sb strings.Builder

				sb.WriteString("Here are your activities:\n")

				for _, nt := range activities {
					sb.WriteString(nt.Desc + "\n")

					if len(nt.Tags) > 0 {
						sb.WriteString(fmt.Sprintf("  Tags: %s\n", strings.Join(nt.Tags, ", ")))
					}

					sb.WriteString("\n")
				}

				return c.Send(sb.String())
			})
		})

		bot.Handle("/addfriend", func(c tele.Context) error {
			payload := c.Message().Payload

			if payload == "" {
				return c.Send(
					"Please provide a friend information in the format:\n/addfriend " + lang.FormatPersonInfo,
				)
			}

			f, err := lang.ExtractPerson(payload)

			if err != nil && !errors.Is(err, lang.ErrNoInfo) {
				return c.Send(fmt.Sprintf("failed to parse friend info: %v", err))
			}

			ctx := context.Background()

			return s.Tx(ctx, func(j *journal.Journal) error {
				j.AddFriend(f)

				if err != nil {
					return c.Send(fmt.Sprintf("failed to add friend: %v", err))
				}

				fmt.Println("Added new friend: " + f.String())

				var sb strings.Builder

				sb.WriteString("✅ Added new friend: " + f.String())

				if len(f.Locations) > 0 {
					sb.WriteString("\n📍 Locations: " + strings.Join(f.Locations, ", "))
				}

				if len(f.Tags) > 0 {
					sb.WriteString("\n🏷️ Tags: " + strings.Join(f.Tags, ", "))
				}

				return c.Send(sb.String())
			})
		})

		bot.Handle("/addlocation", func(c tele.Context) error {
			payload := c.Message().Payload

			if payload == "" {
				return c.Send(
					"Please provide a location information in the format:\n/addlocation " + lang.FormatLocationInfo,
				)
			}

			fmt.Println("Received payload:", payload)

			l, err := lang.ExtractLocation(payload)

			if err != nil && !errors.Is(err, lang.ErrNoInfo) {
				return c.Send(fmt.Sprintf("failed to parse friend info: %v", err))
			}

			ctx := context.Background()

			return s.Tx(ctx, func(j *journal.Journal) error {
				j.AddLocation(l)

				if err != nil {
					return c.Send(fmt.Sprintf("failed to add friend: %v", err))
				}

				var sb strings.Builder

				sb.WriteString("✅ Added new location: " + l.String())

				if len(l.Aliases) > 0 {
					sb.WriteString("\n📍 Aliases: " + strings.Join(l.Aliases, ", "))
				}

				if len(l.Tags) > 0 {
					sb.WriteString("\n🏷️ Tags: " + strings.Join(l.Tags, ", "))
				}

				if l.Desc != "" {
					sb.WriteString("\n🧭 Description: \n" + l.Desc)
				}

				return c.Send(sb.String())
			})
		})

		bot.Handle("/addnote", func(c tele.Context) error {
			payload := c.Message().Payload

			if payload == "" {
				return c.Send(
					"Please provide a note information in the format:\n/addnote " + lang.FormatEventInfo,
				)
			}

			e, err := lang.ExtractEvent(friend.EventTypeNote, payload)

			if err != nil && !errors.Is(err, lang.ErrNoInfo) {
				return c.Send(fmt.Sprintf("failed to parse note info: %v", err))
			}

			ctx := context.Background()

			return s.Tx(ctx, func(j *journal.Journal) error {
				e, err = j.AddEvent(e)
				if err != nil {
					return c.Send(fmt.Sprintf("failed to add note: %v", err))
				}

				fmt.Println("✅ Added new note: " + e.ID)

				var sb strings.Builder

				sb.WriteString("✅ Added new note: " + e.ID)
				sb.WriteString("\n🧭 Description: \n" + e.Desc)

				if len(e.Tags) > 0 {
					sb.WriteString("\n🏷️ Tags: " + strings.Join(e.Tags, ", "))
				}

				return c.Send(sb.String())
			})
		})

		bot.Handle("/addactivity", func(c tele.Context) error {
			payload := c.Message().Payload

			if payload == "" {
				return c.Send(
					"Please provide a activity information in the format:\n/addactivity " + lang.FormatEventInfo,
				)
			}

			e, err := lang.ExtractEvent(friend.EventTypeActivity, payload)

			if err != nil && !errors.Is(err, lang.ErrNoInfo) {
				return c.Send(fmt.Sprintf("failed to parse activity info: %v", err))
			}

			ctx := context.Background()

			return s.Tx(ctx, func(j *journal.Journal) error {
				e, err = j.AddEvent(e)
				if err != nil {
					return c.Send(fmt.Sprintf("failed to add activity: %v", err))
				}

				fmt.Println("✅ Added new activity: " + e.ID)

				var sb strings.Builder

				sb.WriteString("✅ Added new activity: " + e.ID)
				sb.WriteString("\n🧭 Description: \n" + e.Desc)

				if len(e.Tags) > 0 {
					sb.WriteString("\n🏷️ Tags: " + strings.Join(e.Tags, ", "))
				}

				return c.Send(sb.String())
			})
		})

		bot.Handle("/edit", func(c tele.Context) error {
			// Assume person is selected; in real bot use DB lookup
			// personID := "person_123"
			// editMgr.BeginEdit(c.Sender().ID, personID)

			markup := &tele.ReplyMarkup{}
			btnName := markup.Data("✏️ Name", "edit_name")
			btnDesc := markup.Data("📝 Desc", "edit_desc")
			btnTags := markup.Data("🏷️ Tags", "edit_tags")
			btnLoc := markup.Data("📍 Location", "edit_location")

			markup.Inline(markup.Row(btnName, btnDesc), markup.Row(btnTags, btnLoc))

			return c.Send("What field do you want to edit?", markup)
		})

		go func() {
			fmt.Println("Starting Telegram bot...")
			bot.Start()
		}()

		<-ctx.Done()

		fmt.Println("\nStopping Telegram bot...")
		bot.Stop()

		return nil
	},
}
