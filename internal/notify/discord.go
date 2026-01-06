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
	"fmt"
	"net/http"

	"github.com/roma-glushko/frens/internal/config"
)

// DiscordDispatcher sends notifications via Discord webhooks
type DiscordDispatcher struct {
	webhookURL string
	client     *http.Client
}

// NewDiscordDispatcher creates a new Discord dispatcher
func NewDiscordDispatcher(webhookURL string) *DiscordDispatcher {
	return &DiscordDispatcher{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

// Type returns the channel type
func (d *DiscordDispatcher) Type() config.ChannelType {
	return config.ChannelDiscord
}

// discordWebhookPayload represents the Discord webhook message format
type discordWebhookPayload struct {
	Content string `json:"content"`
}

// Send sends a message via Discord webhook
func (d *DiscordDispatcher) Send(
	ctx context.Context,
	_ *ReminderContext,
	destination string,
	message string,
) error {
	webhookURL := d.webhookURL
	if destination != "" {
		// Allow override via destination
		webhookURL = destination
	}

	if webhookURL == "" {
		return fmt.Errorf("discord webhook URL is required")
	}

	payload := discordWebhookPayload{
		Content: message,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal discord payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create discord request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send discord message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("discord webhook returned status %d", resp.StatusCode)
	}

	return nil
}

// ValidateConfig validates the Discord channel configuration
func (d *DiscordDispatcher) ValidateConfig(channelConfig map[string]interface{}) error {
	if _, ok := channelConfig["webhook_url"]; !ok {
		return fmt.Errorf("discord config requires 'webhook_url'")
	}

	return nil
}
