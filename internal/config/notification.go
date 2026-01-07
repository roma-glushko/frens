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

import "strings"

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

// Notifications holds notification-related configuration
type Notifications struct {
	Channels  []NotificationChannel  `toml:"channels"`
	Templates []NotificationTemplate `toml:"templates"`
	Rules     []RoutingRule          `toml:"rules"`
}

// GetChannel returns a channel by ID
func (n *Notifications) GetChannel(id string) *NotificationChannel {
	for i := range n.Channels {
		if n.Channels[i].ID == id {
			return &n.Channels[i]
		}
	}

	return nil
}

// GetTemplate returns a template by ID
func (n *Notifications) GetTemplate(id string) *NotificationTemplate {
	for i := range n.Templates {
		if n.Templates[i].ID == id {
			return &n.Templates[i]
		}
	}

	return nil
}

// GetEnabledChannels returns all enabled channels
func (n *Notifications) GetEnabledChannels() []NotificationChannel {
	enabled := make([]NotificationChannel, 0)

	for _, ch := range n.Channels {
		if ch.Enabled {
			enabled = append(enabled, ch)
		}
	}

	return enabled
}

// MatchRuleCtx contains the context for matching routing rules
type MatchRuleCtx struct {
	Tags    []string
	Content string // Description or other searchable content
}

// MatchRule finds the first matching rule for given tags (convenience wrapper)
func (n *Notifications) MatchRule(tags []string) *RoutingRule {
	return n.MatchRuleWithCtx(MatchRuleCtx{Tags: tags})
}

// MatchRuleWithCtx finds the first matching rule for given context
func (n *Notifications) MatchRuleWithCtx(ctx MatchRuleCtx) *RoutingRule {
	tagSet := make(map[string]bool)

	for _, t := range ctx.Tags {
		tagSet[t] = true
	}

	contentLower := strings.ToLower(ctx.Content)

	for i := range n.Rules {
		rule := &n.Rules[i]

		// Rule with no conditions matches everything
		if len(rule.MatchTags) == 0 && len(rule.MatchKeywords) == 0 {
			return rule
		}

		if rule.MatchAll {
			// AND logic - all conditions must match
			if !allTagsMatch(rule.MatchTags, tagSet) {
				continue
			}

			if !allKeywordsMatch(rule.MatchKeywords, contentLower) {
				continue
			}

			return rule
		}

		// OR logic - any condition matches
		if anyTagMatches(rule.MatchTags, tagSet) {
			return rule
		}

		if anyKeywordMatches(rule.MatchKeywords, contentLower) {
			return rule
		}
	}

	return nil
}

func allTagsMatch(matchTags []string, tagSet map[string]bool) bool {
	for _, mt := range matchTags {
		if !tagSet[mt] {
			return false
		}
	}

	return true
}

func anyTagMatches(matchTags []string, tagSet map[string]bool) bool {
	for _, mt := range matchTags {
		if tagSet[mt] {
			return true
		}
	}

	return false
}

func allKeywordsMatch(keywords []string, content string) bool {
	for _, kw := range keywords {
		if !strings.Contains(content, strings.ToLower(kw)) {
			return false
		}
	}

	return true
}

func anyKeywordMatches(keywords []string, content string) bool {
	for _, kw := range keywords {
		if strings.Contains(content, strings.ToLower(kw)) {
			return true
		}
	}

	return false
}
