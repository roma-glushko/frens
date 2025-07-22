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
	"fmt"
	"gopkg.in/telebot.v4/middleware"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/roma-glushko/frens/internal/version"
	"github.com/urfave/cli/v2"
	tele "gopkg.in/telebot.v4"
)

const helpMessage = `Welcome to the Frens Bot! Here is what you can do:
/listfrens - List my friends.
/listlocs - List my locations.
/listnotes - List my notes.
/listactivities - List my activities.

/addfriend - Add a new friend.
/addlocation - Add a new location.
/addnote - Add a new note.
/addactivity - Add a new activity.

/version - Show the current version of frens.
`

var BotCommand = &cli.Command{
	Name:  "bot",
	Usage: "Start a Telegram bot",
	Action: func(c *cli.Context) error {
		userList := []int64{283564721}
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

		bot.Use(middleware.Logger())
		bot.Use(middleware.Whitelist(userList...))

		bot.Handle("/help", func(c tele.Context) error {
			return c.Send(helpMessage)
		})

		bot.Handle("/start", func(c tele.Context) error {
			return c.Send(helpMessage)
		})

		bot.Handle("/version", func(c tele.Context) error {
			return c.Send("Frens Version: " + version.FullVersion)
		})

		go func() {
			fmt.Println("Starting Telegram bot...")
			bot.Start()
		}()

		<-c.Done()

		fmt.Println("\nStopping Telegram bot...")
		bot.Stop()

		return nil
	},
}
