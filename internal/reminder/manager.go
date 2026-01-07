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
	"context"
	"fmt"
	"time"

	"github.com/roma-glushko/frens/internal/config"
	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/notify"
)

// Manager handles reminder operations including checking and firing reminders
type Manager struct {
	journal    *journal.Journal
	notifyConf *config.Notifications
	sender     *notify.NotificationSender
}

// NewManager creates a new reminder manager
func NewManager(
	j *journal.Journal,
	notifyConf *config.Notifications,
	sender *notify.NotificationSender,
) *Manager {
	return &Manager{
		journal:    j,
		notifyConf: notifyConf,
		sender:     sender,
	}
}

// FireResult represents the result of firing a single reminder
type FireResult struct {
	Reminder    *friend.Reminder
	Success     bool
	Error       error
	SendResults []notify.SendResult
}

// CheckAndFire checks for due reminders and fires them
func (m *Manager) CheckAndFire(
	ctx context.Context,
	now time.Time,
	dryRun bool,
) ([]FireResult, error) {
	dueReminders := m.journal.GetDueReminders(now)
	results := make([]FireResult, 0, len(dueReminders))

	for i := range dueReminders {
		r := &dueReminders[i]
		result := FireResult{Reminder: r}

		if dryRun {
			result.Success = true
			results = append(results, result)

			continue
		}

		fireResult, err := m.Fire(ctx, r, now)
		if err != nil {
			result.Error = err
			result.Success = false
		} else {
			result.Success = true
			result.SendResults = fireResult
		}

		results = append(results, result)
	}

	return results, nil
}

// Fire sends notifications for a reminder and updates its state
func (m *Manager) Fire(
	ctx context.Context,
	r *friend.Reminder,
	now time.Time,
) ([]notify.SendResult, error) {
	// Resolve linked entity and friend
	linkedEntity, friendRef, err := m.journal.ResolveLinkedEntity(r)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve linked entity: %w", err)
	}

	// Create reminder context
	rc := &notify.ReminderContext{
		CreatedAt:    now,
		Friend:       friendRef,
		Reminder:     r,
		Event:        linkedEntity.Event,
		Date:         linkedEntity.Date,
		WishlistItem: linkedEntity.WishlistItem,
	}

	// Get template
	rule := m.notifyConf.MatchRule(r.Tags)
	template := notify.GetTemplateForReminder(m.notifyConf, rule, r.LinkedEntityType)

	// Send notifications
	sendResults, err := m.sender.Send(ctx, rc, template)
	if err != nil {
		return sendResults, fmt.Errorf("failed to send notifications: %w", err)
	}

	// Update reminder state
	for _, jr := range m.journal.Reminders {
		if jr.ID == r.ID {
			jr.LastFiredAt = now
			jr.FiredCount++

			if jr.Recurrence == friend.RecurrenceOnce {
				jr.State = friend.ReminderStateFired
			} else {
				// Compute next trigger date for recurring reminders
				jr.TriggerAt = m.computeNextOccurrence(jr, now)
				jr.State = friend.ReminderStatePending
			}

			m.journal.SetDirty(true)

			break
		}
	}

	return sendResults, nil
}

func (m *Manager) computeNextOccurrence(r *friend.Reminder, now time.Time) time.Time {
	baseDate := r.TriggerAt

	switch r.Recurrence {
	case friend.RecurrenceYearly:
		// Add one year
		nextYear := baseDate.AddDate(1, 0, 0)
		// If next occurrence is still in the past, keep adding years
		for nextYear.Before(now) {
			nextYear = nextYear.AddDate(1, 0, 0)
		}

		return nextYear

	case friend.RecurrenceMonthly:
		// Add one month
		nextMonth := baseDate.AddDate(0, 1, 0)
		// If next occurrence is still in the past, keep adding months
		for nextMonth.Before(now) {
			nextMonth = nextMonth.AddDate(0, 1, 0)
		}

		return nextMonth

	default:
		return baseDate
	}
}

// Acknowledge marks a reminder as acknowledged
func (m *Manager) Acknowledge(id string) error {
	r, err := m.journal.GetReminder(id)
	if err != nil {
		return err
	}

	r.State = friend.ReminderStateAcknowledged

	_, err = m.journal.UpdateReminder(r, r)

	return err
}

// GetUpcoming returns reminders due within the specified number of days
func (m *Manager) GetUpcoming(now time.Time, days int) []friend.Reminder {
	return m.journal.GetUpcomingReminders(now, days)
}

// GetDue returns reminders that are due now
func (m *Manager) GetDue(now time.Time) []friend.Reminder {
	return m.journal.GetDueReminders(now)
}
