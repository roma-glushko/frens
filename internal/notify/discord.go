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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/roma-glushko/frens/internal/config"
)

// DiscordNotifier sends notifications via Discord webhooks
type DiscordNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewDiscordNotifier creates a new Discord dispatcher
func NewDiscordNotifier(webhookURL string) *DiscordNotifier {
	return &DiscordNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

// Type returns the channel type
func (n *DiscordNotifier) Type() config.ChannelType {
	return config.ChannelDiscord
}

// discordWebhookPayload represents the Discord webhook message format
type discordWebhookPayload struct {
	Content string `json:"content"`
}

// Send sends a message via Discord webhook
func (n *DiscordNotifier) Send(
	ctx context.Context,
	_ *ReminderContext,
	destination string,
	message string,
) error {
	webhookURL := n.webhookURL

	if destination != "" {
		// Allow override via destination
		webhookURL = destination
	}

	if webhookURL == "" {
		return errors.New("discord webhook URL is required")
	}

	payload := discordWebhookPayload{
		Content: message,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal discord payload: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		webhookURL,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to create discord request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send discord message: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("discord webhook returned status %d", resp.StatusCode)
	}

	return nil
}

// ValidateConfig validates the Discord channel configuration
func (n *DiscordNotifier) ValidateConfig(channelConfig map[string]interface{}) error {
	if _, ok := channelConfig["webhook_url"]; !ok {
		return errors.New("discord config requires 'webhook_url'")
	}

	return nil
}
