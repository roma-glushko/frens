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
)

// ReminderContext contains all context needed to send a reminder notification
type ReminderContext struct {
	Reminder     *friend.Reminder
	LinkedEntity interface{} // *friend.Date, *friend.WishlistItem, *friend.Event
	Friend       *friend.Person
}

// Dispatcher is the interface that notification channels must implement
type Dispatcher interface {
	// Type returns the channel type identifier
	Type() config.ChannelType

	// Send sends a notification to the specified destination
	Send(ctx context.Context, rc *ReminderContext, destination string, message string) error

	// ValidateConfig validates the channel configuration
	ValidateConfig(channelConfig map[string]interface{}) error
}

// DispatcherRegistry manages available notification dispatchers
type DispatcherRegistry struct {
	dispatchers map[config.ChannelType]Dispatcher
}

// NewDispatcherRegistry creates a new dispatcher registry
func NewDispatcherRegistry() *DispatcherRegistry {
	return &DispatcherRegistry{
		dispatchers: make(map[config.ChannelType]Dispatcher),
	}
}

// Register adds a dispatcher to the registry
func (r *DispatcherRegistry) Register(d Dispatcher) {
	r.dispatchers[d.Type()] = d
}

// Get retrieves a dispatcher by channel type
func (r *DispatcherRegistry) Get(channelType config.ChannelType) (Dispatcher, bool) {
	d, ok := r.dispatchers[channelType]

	return d, ok
}

// GetAll returns all registered dispatchers
func (r *DispatcherRegistry) GetAll() map[config.ChannelType]Dispatcher {
	return r.dispatchers
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
	registry   *DispatcherRegistry
	notifyConf *config.NotificationConfig
}

// NewNotificationSender creates a new notification sender
func NewNotificationSender(
	registry *DispatcherRegistry,
	notifyConf *config.NotificationConfig,
) *NotificationSender {
	return &NotificationSender{
		registry:   registry,
		notifyConf: notifyConf,
	}
}

// Send sends a notification to appropriate channels based on routing rules
func (s *NotificationSender) Send(
	ctx context.Context,
	rc *ReminderContext,
	template *config.NotificationTemplate,
) ([]SendResult, error) {
	results := make([]SendResult, 0)

	// Collect tags from reminder
	tags := rc.Reminder.Tags

	// Build content for keyword matching
	content := buildContentForMatching(rc)

	// Find matching rule
	rule := s.notifyConf.MatchRuleWithContext(config.MatchRuleContext{
		Tags:    tags,
		Content: content,
	})

	var channelIDs []string
	var destination string

	if rule != nil {
		channelIDs = rule.ChannelIDs
		destination = rule.Destination
	} else {
		// Use default channels
		channelIDs = s.notifyConf.GetDefaultChannels()
	}

	if len(channelIDs) == 0 {
		return results, fmt.Errorf("no notification channels configured")
	}

	// Render message
	templateCtx := NewTemplateContext(rc.Reminder, rc.LinkedEntity, rc.Friend, time.Now())

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

		dispatcher, ok := s.registry.Get(channel.Type)
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

		err := dispatcher.Send(ctx, rc, dest, message)
		results = append(results, SendResult{
			ChannelID:   chID,
			Destination: dest,
			Success:     err == nil,
			Error:       err,
		})
	}

	return results, nil
}

// buildContentForMatching builds a searchable content string from reminder context
func buildContentForMatching(rc *ReminderContext) string {
	var parts []string

	// Add reminder description
	if rc.Reminder.Desc != "" {
		parts = append(parts, rc.Reminder.Desc)
	}

	// Add friend name if available
	if rc.Friend != nil {
		parts = append(parts, rc.Friend.Name)
		if rc.Friend.Desc != "" {
			parts = append(parts, rc.Friend.Desc)
		}
	}

	// Add linked entity description
	switch entity := rc.LinkedEntity.(type) {
	case *friend.Date:
		if entity.Desc != "" {
			parts = append(parts, entity.Desc)
		}

		parts = append(parts, entity.DateExpr)
	case *friend.WishlistItem:
		if entity.Desc != "" {
			parts = append(parts, entity.Desc)
		}
	case *friend.Event:
		if entity.Desc != "" {
			parts = append(parts, entity.Desc)
		}
	}

	return strings.Join(parts, " ")
}
