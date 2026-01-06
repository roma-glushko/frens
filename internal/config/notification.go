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

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
)

const NotificationConfigFile = "notifications.toml"

type ChannelType string

const (
	ChannelTelegram ChannelType = "telegram"
	ChannelDiscord  ChannelType = "discord"
	ChannelWhatsApp ChannelType = "whatsapp"
	ChannelEmail    ChannelType = "email"
)

// NotificationChannel represents a configured notification channel
type NotificationChannel struct {
	ID      string                 `toml:"id"`
	Type    ChannelType            `toml:"type"`
	Name    string                 `toml:"name"`
	Enabled bool                   `toml:"enabled"`
	Config  map[string]interface{} `toml:"config"`
}

// NotificationTemplate represents a customizable notification message template
type NotificationTemplate struct {
	ID      string `toml:"id"`
	Name    string `toml:"name"`
	Subject string `toml:"subject,omitempty"` // For email
	Body    string `toml:"body"`
	Format  string `toml:"format"` // "text", "markdown", "html"
}

// RoutingRule defines how reminders are routed to notification channels
type RoutingRule struct {
	ID            string   `toml:"id"`
	Name          string   `toml:"name"`
	Priority      int      `toml:"priority"`              // Lower = higher priority
	MatchTags     []string `toml:"match_tags"`            // Tags to match
	MatchKeywords []string `toml:"match_keywords"`        // Keywords to match in description
	MatchAll      bool     `toml:"match_all"`             // If true, all conditions must match (AND logic)
	ChannelIDs    []string `toml:"channel_ids"`           // Channels to send to
	Destination   string   `toml:"destination"`           // Channel-specific target (chat_id, email, etc.)
	TemplateID    string   `toml:"template_id,omitempty"` // Optional template override
}

// NotificationConfig holds the complete notification configuration
type NotificationConfig struct {
	Channels  []NotificationChannel  `toml:"channels"`
	Templates []NotificationTemplate `toml:"templates"`
	Rules     []RoutingRule          `toml:"rules"`
	Default   struct {
		ChannelIDs []string `toml:"channel_ids"`
	} `toml:"default"`
}

// LoadNotificationConfig loads the notification configuration from the config directory
func LoadNotificationConfig(configDir string) (*NotificationConfig, error) {
	configPath := filepath.Join(configDir, NotificationConfigFile)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Return empty config if file doesn't exist
		return &NotificationConfig{}, nil
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open notification config: %w", err)
	}
	defer file.Close()

	var config NotificationConfig

	decoder := toml.NewDecoder(file)
	if _, err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode notification config: %w", err)
	}

	// Sort rules by priority
	sort.SliceStable(config.Rules, func(i, j int) bool {
		return config.Rules[i].Priority < config.Rules[j].Priority
	})

	return &config, nil
}

// SaveNotificationConfig saves the notification configuration to the config directory
func SaveNotificationConfig(configDir string, config *NotificationConfig) error {
	configPath := filepath.Join(configDir, NotificationConfigFile)

	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create notification config file: %w", err)
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to encode notification config: %w", err)
	}

	return nil
}

// GetChannel returns a channel by ID
func (c *NotificationConfig) GetChannel(id string) *NotificationChannel {
	for i := range c.Channels {
		if c.Channels[i].ID == id {
			return &c.Channels[i]
		}
	}

	return nil
}

// GetTemplate returns a template by ID
func (c *NotificationConfig) GetTemplate(id string) *NotificationTemplate {
	for i := range c.Templates {
		if c.Templates[i].ID == id {
			return &c.Templates[i]
		}
	}

	return nil
}

// GetEnabledChannels returns all enabled channels
func (c *NotificationConfig) GetEnabledChannels() []NotificationChannel {
	enabled := make([]NotificationChannel, 0)

	for _, ch := range c.Channels {
		if ch.Enabled {
			enabled = append(enabled, ch)
		}
	}

	return enabled
}

// MatchRuleContext contains the context for matching routing rules
type MatchRuleContext struct {
	Tags    []string
	Content string // Description or other searchable content
}

// MatchRule finds the first matching rule for given tags (convenience wrapper)
func (c *NotificationConfig) MatchRule(tags []string) *RoutingRule {
	return c.MatchRuleWithContext(MatchRuleContext{Tags: tags})
}

// MatchRuleWithContext finds the first matching rule for given context
func (c *NotificationConfig) MatchRuleWithContext(ctx MatchRuleContext) *RoutingRule {
	tagSet := make(map[string]bool)
	for _, t := range ctx.Tags {
		tagSet[t] = true
	}

	contentLower := strings.ToLower(ctx.Content)

	for i := range c.Rules {
		rule := &c.Rules[i]

		// Rule with no conditions matches everything
		if len(rule.MatchTags) == 0 && len(rule.MatchKeywords) == 0 {
			return rule
		}

		if rule.MatchAll {
			// AND logic - all conditions must match
			if !c.allTagsMatch(rule.MatchTags, tagSet) {
				continue
			}

			if !c.allKeywordsMatch(rule.MatchKeywords, contentLower) {
				continue
			}

			return rule
		}

		// OR logic - any condition matches
		if c.anyTagMatches(rule.MatchTags, tagSet) {
			return rule
		}

		if c.anyKeywordMatches(rule.MatchKeywords, contentLower) {
			return rule
		}
	}

	return nil
}

func (c *NotificationConfig) allTagsMatch(matchTags []string, tagSet map[string]bool) bool {
	for _, mt := range matchTags {
		if !tagSet[mt] {
			return false
		}
	}

	return true
}

func (c *NotificationConfig) anyTagMatches(matchTags []string, tagSet map[string]bool) bool {
	for _, mt := range matchTags {
		if tagSet[mt] {
			return true
		}
	}

	return false
}

func (c *NotificationConfig) allKeywordsMatch(keywords []string, content string) bool {
	for _, kw := range keywords {
		if !strings.Contains(content, strings.ToLower(kw)) {
			return false
		}
	}

	return true
}

func (c *NotificationConfig) anyKeywordMatches(keywords []string, content string) bool {
	for _, kw := range keywords {
		if strings.Contains(content, strings.ToLower(kw)) {
			return true
		}
	}

	return false
}

// GetDefaultChannels returns the default channel IDs
func (c *NotificationConfig) GetDefaultChannels() []string {
	return c.Default.ChannelIDs
}
