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

package friend

import (
	"errors"
	"time"
)

type Recurrence = string

var (
	RecurrenceOnce    Recurrence = "once"
	RecurrenceMonthly Recurrence = "monthly"
	RecurrenceYearly  Recurrence = "yearly"
)

type OffsetDirection = string

var (
	OffsetDirectionBefore OffsetDirection = "before"
	OffsetDirectionAfter  OffsetDirection = "after"
)

type ReminderState = string

var (
	ReminderStatePending      ReminderState = "pending"
	ReminderStateFired        ReminderState = "fired"
	ReminderStateAcknowledged ReminderState = "acknowledged"
)

type LinkedEntityType = string

var (
	LinkedEntityDate     LinkedEntityType = "date"
	LinkedEntityWishlist LinkedEntityType = "wishlist"
	LinkedEntityActivity LinkedEntityType = "activity"
	LinkedEntityNote     LinkedEntityType = "note"
)

// Reminder represents a scheduled reminder linked to an entity (date, wishlist, activity, note)
type Reminder struct {
	ID        string    `toml:"id"         json:"id"`
	CreatedAt time.Time `toml:"created_at" json:"createdAt"`

	// Linked entity reference
	LinkedEntityType LinkedEntityType `toml:"linked_type" json:"linkedType"`
	LinkedEntityID   string           `toml:"linked_id"   json:"linkedId"`
	FriendID         string           `toml:"friend_id,omitempty" json:"friendId,omitempty"`

	// Schedule configuration
	TriggerAt    time.Time `toml:"trigger_at"    json:"triggerAt"`
	ScheduleExpr string    `toml:"schedule_expr" json:"scheduleExpr"`

	// Recurrence and offset
	Recurrence      Recurrence      `toml:"recurrence"                  json:"recurrence"`
	OffsetDirection OffsetDirection `toml:"offset_direction,omitempty"  json:"offsetDirection,omitempty"`
	OffsetDuration  time.Duration   `toml:"offset_duration,omitempty"   json:"offsetDuration,omitempty"`

	// State tracking
	State       ReminderState `toml:"state"                   json:"state"`
	LastFiredAt time.Time     `toml:"last_fired_at,omitempty" json:"lastFiredAt,omitempty"`
	FiredCount  int           `toml:"fired_count"             json:"firedCount"`

	// Optional description override and tags
	Desc string   `toml:"desc,omitempty" json:"desc,omitempty"`
	Tags []string `toml:"tags,omitempty" json:"tags,omitempty"`
}

func (r *Reminder) SetTags(tags []string) { r.Tags = tags }
func (r *Reminder) GetTags() []string     { return r.Tags }

func (r *Reminder) Validate() error {
	if r.LinkedEntityType == "" {
		return errors.New("linked entity type is required")
	}

	if r.LinkedEntityID == "" {
		return errors.New("linked entity ID is required")
	}

	if r.TriggerAt.IsZero() && r.ScheduleExpr == "" {
		return errors.New("either trigger date or schedule expression is required")
	}

	return nil
}

func (r *Reminder) IsDue(now time.Time) bool {
	return r.State == ReminderStatePending && !r.TriggerAt.After(now)
}

func (r *Reminder) ComputeNextTrigger(baseDate time.Time) time.Time {
	target := baseDate

	if r.OffsetDuration > 0 {
		if r.OffsetDirection == OffsetDirectionBefore {
			target = target.Add(-r.OffsetDuration)
		} else {
			target = target.Add(r.OffsetDuration)
		}
	}

	return target
}
