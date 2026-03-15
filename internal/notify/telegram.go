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

package notify

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/roma-glushko/frens/internal/config"
	tele "gopkg.in/telebot.v4"
)

// TelegramDispatcher sends notifications via Telegram
type TelegramDispatcher struct {
	bot *tele.Bot
}

// NewTelegramDispatcher creates a new Telegram dispatcher
func NewTelegramDispatcher(token string) (*TelegramDispatcher, error) {
	pref := tele.Settings{
		Token: token,
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot: %w", err)
	}

	return &TelegramDispatcher{bot: bot}, nil
}

// Type returns the channel type
func (d *TelegramDispatcher) Type() config.ChannelType {
	return config.ChannelTelegram
}

// Send sends a message to the specified Telegram chat
func (d *TelegramDispatcher) Send(
	_ context.Context,
	_ *ReminderContext,
	destination string,
	message string,
) error {
	if destination == "" {
		return errors.New("telegram destination (chat_id) is required")
	}

	chatID, err := strconv.ParseInt(destination, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid telegram chat_id: %w", err)
	}

	chat, err := d.bot.ChatByID(chatID)
	if err != nil {
		return fmt.Errorf("failed to get telegram chat: %w", err)
	}

	_, err = d.bot.Send(chat, message)
	if err != nil {
		return fmt.Errorf("failed to send telegram message: %w", err)
	}

	return nil
}

// ValidateConfig validates the Telegram channel configuration
func (d *TelegramDispatcher) ValidateConfig(channelConfig map[string]interface{}) error {
	if _, ok := channelConfig["token"]; !ok {
		return errors.New("telegram config requires 'token'")
	}

	return nil
}
