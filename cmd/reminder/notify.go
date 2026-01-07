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
	"os"
	"time"

	"github.com/roma-glushko/frens/internal/config"
	jctx "github.com/roma-glushko/frens/internal/context"
	"github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/log"
	"github.com/roma-glushko/frens/internal/notify"
	"github.com/roma-glushko/frens/internal/reminder"
	"github.com/urfave/cli/v2"
)

var NotifyCommand = &cli.Command{
	Name:  "notify",
	Usage: "Check for due reminders and send notifications",
	Description: `Checks all pending reminders and sends notifications for those that are due.
Run this command periodically via cron to receive timely notifications.

Example cron entry (every hour):
  0 * * * * frens reminder notify`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "dry-run",
			Usage: "Preview due reminders without sending notifications",
		},
	},
	Action: func(c *cli.Context) error {
		dryRun := c.Bool("dry-run")
		ctx := c.Context
		appCtx := jctx.FromCtx(ctx)
		now := time.Now()

		// Load config
		cfg, err := config.Load(appCtx.JournalDir)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		return appCtx.Store.Tx(ctx, func(j *journal.Journal) error {
			// Set up notification registry
			registry := notify.NewRegistry()

			if err := registerNotifiers(registry, &cfg.Notifications); err != nil {
				return err
			}

			// Create sender and manager
			sender := notify.NewNotificationSender(registry, &cfg.Notifications)
			mgr := reminder.NewManager(j, &cfg.Notifications, sender)

			// Check and fire reminders
			results, err := mgr.CheckAndFire(ctx, now, dryRun)
			if err != nil {
				return fmt.Errorf("failed to check reminders: %w", err)
			}

			if len(results) == 0 {
				log.Info("No due reminders")
				return nil
			}

			// Print results
			if dryRun {
				log.Infof("Found %d due reminder(s) (dry-run):\n", len(results))
			} else {
				log.Infof("Processed %d reminder(s):\n", len(results))
			}

			for _, r := range results {
				printNotifyResult(r, dryRun)
			}

			return nil
		})
	},
}

func registerNotifiers(registry *notify.Registry, notifications *config.Notifications) error {
	for _, ch := range notifications.GetEnabledChannels() {
		switch ch.Type {
		case config.ChannelTelegram:
			token, _ := ch.Config["token"].(string)

			if token == "" {
				token = os.Getenv("TELEGRAM_BOT_TOKEN")
			}

			if token == "" {
				log.Warnf("Telegram channel '%s' has no token configured, skipping", ch.ID)
				continue
			}

			dispatcher, err := notify.NewTelegramDispatcher(token)
			if err != nil {
				log.Warnf("Failed to create Telegram dispatcher for '%s': %v", ch.ID, err)
				continue
			}

			registry.Register(dispatcher)

		case config.ChannelDiscord:
			webhookURL, _ := ch.Config["webhook_url"].(string)

			if webhookURL == "" {
				webhookURL = os.Getenv("DISCORD_WEBHOOK_URL")
			}

			if webhookURL == "" {
				log.Warnf("Discord channel '%s' has no webhook_url configured, skipping", ch.ID)
				continue
			}

			registry.Register(notify.NewDiscordNotifier(webhookURL))

		case config.ChannelWhatsApp, config.ChannelEmail:
			log.Warnf("Channel type '%s' is not yet implemented, skipping '%s'", ch.Type, ch.ID)
		}
	}

	return nil
}

func printNotifyResult(r reminder.FireResult, dryRun bool) {
	icon := getResultIcon(r.Success, dryRun)

	fmt.Printf(
		"%s [%s] %s:%s\n",
		icon,
		r.Reminder.ID[:8],
		r.Reminder.LinkedEntityType,
		r.Reminder.LinkedEntityID[:8],
	)

	if r.Reminder.Desc != "" {
		fmt.Printf("   %s\n", r.Reminder.Desc)
	}

	if r.Error != nil {
		fmt.Printf("   Error: %v\n", r.Error)
	}

	if !dryRun {
		printSendResults(r.SendResults)
	}

	fmt.Println()
}

func getResultIcon(success, dryRun bool) string {
	if dryRun {
		return "📋"
	}

	if success {
		return "✅"
	}

	return "❌"
}

func printSendResults(results []notify.SendResult) {
	for _, sr := range results {
		icon := "✓"
		if !sr.Success {
			icon = "✗"
		}

		fmt.Printf("   %s Channel: %s", icon, sr.ChannelID)

		if sr.Destination != "" {
			fmt.Printf(" → %s", sr.Destination)
		}

		if sr.Error != nil {
			fmt.Printf(" (%v)", sr.Error)
		}

		fmt.Println()
	}
}
