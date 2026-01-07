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
	"fmt"
	"strings"
	"time"

	"github.com/roma-glushko/frens/internal/config"
	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/log"
)

// ReminderContext contains all context needed to send a reminder notification
type ReminderContext struct {
	CreatedAt    time.Time
	Reminder     *friend.Reminder
	Event        *friend.Event
	Date         *friend.Date
	WishlistItem *friend.WishlistItem
	Friend       *friend.Person
}

// Notifier is the interface that notification channels must implement
type Notifier interface {
	// Type returns the channel type identifier
	Type() config.ChannelType

	// Send sends a notification to the specified destination
	Send(ctx context.Context, rc *ReminderContext, destination string, message string) error

	// ValidateConfig validates the channel configuration
	ValidateConfig(channelConfig map[string]interface{}) error
}

// Registry manages available notification dispatchers
type Registry struct {
	notifiers map[config.ChannelType]Notifier
}

// NewRegistry creates a new dispatcher registry
func NewRegistry() *Registry {
	return &Registry{
		notifiers: make(map[config.ChannelType]Notifier),
	}
}

// Register adds a dispatcher to the registry
func (r *Registry) Register(n Notifier) {
	r.notifiers[n.Type()] = n
}

// Get retrieves a dispatcher by channel type
func (r *Registry) Get(channelType config.ChannelType) (Notifier, bool) {
	n, ok := r.notifiers[channelType]

	return n, ok
}

// GetAll returns all registered dispatchers
func (r *Registry) GetAll() map[config.ChannelType]Notifier {
	return r.notifiers
}

// SendResult represents the result of sending a notification
type SendResult struct {
	ChannelID   string
	Destination string
	Success     bool
	Error       error
}

// NotificationSender handles sending notifications through configured channels
type NotificationSender struct {
	registry   *Registry
	notifyConf *config.Notifications
}

// NewNotificationSender creates a new notification sender
func NewNotificationSender(
	dr *Registry,
	nc *config.Notifications,
) *NotificationSender {
	return &NotificationSender{
		registry:   dr,
		notifyConf: nc,
	}
}

// Send sends a notification to appropriate channels based on routing rules
func (s *NotificationSender) Send(
	ctx context.Context,
	rc *ReminderContext,
	template *config.NotificationTemplate,
) ([]SendResult, error) {
	results := make([]SendResult, 0)

	tags := rc.Reminder.Tags

	rule := s.notifyConf.MatchRuleWithCtx(config.MatchRuleCtx{
		Tags:    tags,
		Content: buildContentForMatching(rc),
	})

	channelIDs, destination := s.resolveChannelsAndDestination(rule)

	if len(channelIDs) == 0 {
		log.Warn(
			"no routing rule matched, notification not sent (add a catch-all rule to prevent this)",
		)

		return results, nil
	}

	// Render message
	templateCtx := NewTemplateContext(rc, time.Now().UTC())

	message, err := RenderTemplate(template, templateCtx)
	if err != nil {
		return results, fmt.Errorf("failed to render template: %w", err)
	}

	// Send to each channel
	for _, chID := range channelIDs {
		channel := s.notifyConf.GetChannel(chID)

		if channel == nil {
			results = append(results, SendResult{
				ChannelID: chID,
				Success:   false,
				Error:     fmt.Errorf("channel not found: %s", chID),
			})

			continue
		}

		if !channel.Enabled {
			continue
		}

		notifier, ok := s.registry.Get(channel.Type)

		if !ok {
			results = append(results, SendResult{
				ChannelID: chID,
				Success:   false,
				Error:     fmt.Errorf("no dispatcher for channel type: %s", channel.Type),
			})

			continue
		}

		dest := destination

		if dest == "" {
			// Try to get default destination from channel config
			if defaultDest, ok := channel.Config["default_chat_id"].(string); ok {
				dest = defaultDest
			} else if defaultDest, ok := channel.Config["default_destination"].(string); ok {
				dest = defaultDest
			}
		}

		err := notifier.Send(ctx, rc, dest, message)

		results = append(results, SendResult{
			ChannelID:   chID,
			Destination: dest,
			Success:     err == nil,
			Error:       err,
		})
	}

	return results, nil
}

func (s *NotificationSender) resolveChannelsAndDestination(
	rule *config.RoutingRule,
) ([]string, string) {
	if rule == nil {
		return nil, ""
	}

	return rule.ChannelIDs, rule.Destination
}

// buildContentForMatching builds a searchable content string from reminder context
func buildContentForMatching(rc *ReminderContext) string {
	var parts strings.Builder

	// Add reminder description
	if rc.Reminder.Desc != "" {
		parts.WriteString(rc.Reminder.Desc)
		parts.WriteString(" ")
	}

	// Add friend name if available
	if rc.Friend != nil {
		parts.WriteString(rc.Friend.Name)
		parts.WriteString(" ")

		if rc.Friend.Desc != "" {
			parts.WriteString(rc.Friend.Desc)
			parts.WriteString(" ")
		}
	}

	if rc.Date != nil {
		if rc.Date.Desc != "" {
			parts.WriteString(rc.Date.Desc)
			parts.WriteString(" ")
		}

		parts.WriteString(rc.Date.DateExpr)
		parts.WriteString(" ")
	}

	if rc.WishlistItem != nil && rc.WishlistItem.Desc != "" {
		parts.WriteString(rc.WishlistItem.Desc)
		parts.WriteString(" ")
	}

	if rc.Event != nil && rc.Event.Desc != "" {
		parts.WriteString(rc.Event.Desc)
		parts.WriteString(" ")
	}

	return strings.TrimSpace(parts.String())
}
