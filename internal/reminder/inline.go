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
	"fmt"
	"time"

	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/lang"
)

// SyncAction describes what happened during a reminder sync.
type SyncAction int

const (
	// SyncNone means no reminder change occurred.
	SyncNone SyncAction = iota
	// SyncCreated means a new reminder was created (or replaced an old one).
	SyncCreated
	// SyncRemoved means an existing reminder was removed by the user.
	SyncRemoved
)

// SyncResult holds the outcome of a reminder sync operation.
type SyncResult struct {
	Action   SyncAction
	Reminder *friend.Reminder // set when Action == SyncCreated
	Err      error
}

// FindForEntity returns existing pending reminders linked to the given entity.
func FindForEntity(j *journal.Journal, entityID string) []friend.Reminder {
	return j.ListReminders(friend.ListReminderQuery{
		LinkedEntityID: entityID,
		State:          friend.ReminderStatePending,
	})
}

// RenderForEditor appends existing reminder !r[...] syntax to the editor text.
func RenderForEditor(text string, reminders []friend.Reminder) string {
	if len(reminders) == 0 {
		return text
	}

	// Only render the first pending reminder (most common case)
	return text + " " + lang.RenderReminder(&reminders[0])
}

// CreateFromAdd extracts an inline reminder from text and creates it.
func CreateFromAdd(
	j *journal.Journal,
	infoTxt string,
	entityType friend.LinkedEntityType,
	entityID string,
	friendID string,
	baseDate time.Time,
	tags []string,
) SyncResult {
	now := time.Now()

	newReminder, err := lang.ExtractReminder(infoTxt, entityType, entityID, friendID, baseDate, now, tags)
	if err != nil {
		return SyncResult{Err: fmt.Errorf("failed to parse reminder: %w", err)}
	}

	if newReminder == nil {
		return SyncResult{Action: SyncNone}
	}

	created, err := j.AddReminder(*newReminder)
	if err != nil {
		return SyncResult{Err: fmt.Errorf("failed to create reminder: %w", err)}
	}

	return SyncResult{Action: SyncCreated, Reminder: &created}
}

// SyncFromEdit reconciles the reminder state after an entity edit.
// It compares the editor text against existing reminders and creates, replaces,
// or removes reminders accordingly.
func SyncFromEdit(
	j *journal.Journal,
	infoTxt string,
	entityType friend.LinkedEntityType,
	entityID string,
	friendID string,
	baseDate time.Time,
	tags []string,
	existingReminders []friend.Reminder,
) SyncResult {
	now := time.Now()

	if !lang.HasReminderExpr(infoTxt) {
		if len(existingReminders) == 0 {
			return SyncResult{Action: SyncNone}
		}

		// User removed !r[...] from the editor — delete existing reminders
		if err := j.RemoveReminders(existingReminders); err != nil {
			return SyncResult{Action: SyncRemoved, Err: fmt.Errorf("failed to remove reminder: %w", err)}
		}

		return SyncResult{Action: SyncRemoved}
	}

	// User has a !r[...] in the editor — create or replace
	newReminder, err := lang.ExtractReminder(infoTxt, entityType, entityID, friendID, baseDate, now, tags)
	if err != nil {
		return SyncResult{Err: fmt.Errorf("failed to parse reminder: %w", err)}
	}

	if newReminder == nil {
		return SyncResult{Action: SyncNone}
	}

	// Remove old pending reminders for this entity
	if len(existingReminders) > 0 {
		if err := j.RemoveReminders(existingReminders); err != nil {
			return SyncResult{Err: fmt.Errorf("failed to remove old reminder: %w", err)}
		}
	}

	created, err := j.AddReminder(*newReminder)
	if err != nil {
		return SyncResult{Err: fmt.Errorf("failed to create reminder: %w", err)}
	}

	return SyncResult{Action: SyncCreated, Reminder: &created}
}
