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

package acceptance

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/roma-glushko/frens/cmd"
	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/store/file"
	"github.com/stretchr/testify/require"
	cli "github.com/urfave/cli/v2"
)

// seedFriendWithDate adds a friend and a date via CLI, returning the date ID.
func seedFriendWithDate(t *testing.T, app cli.App, jDir string) string {
	t.Helper()
	ctx := t.Context()

	err := app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "add",
		"John Doe :: A good friend #friends @NewYork $id:john_doe",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "date", "add",
		"john_doe",
		"birthday :: 2000-05-13 #birthday",
	})
	require.NoError(t, err)

	dates := loadDates(t, jDir)
	require.NotEmpty(t, dates, "expected at least one date after seedFriendWithDate")

	return dates[0].ID
}

// seedReminders directly writes reminder records into reminders.toml for controlled test state.
func seedReminders(t *testing.T, jDir string, reminders []*friend.Reminder) {
	t.Helper()

	store := file.NewTOMLFileStore(jDir)

	j, err := store.Load(context.Background())
	require.NoError(t, err)

	j.Reminders = append(j.Reminders, reminders...)
	j.SetDirty(true)

	err = store.Save(context.Background(), j)
	require.NoError(t, err)
}

func loadDates(t *testing.T, jDir string) []*friend.Date {
	t.Helper()

	store := file.NewTOMLFileStore(jDir)

	j, err := store.Load(context.Background())
	require.NoError(t, err)

	var dates []*friend.Date

	for _, f := range j.Friends {
		dates = append(dates, f.Dates...)
	}

	return dates
}

func writeConfig(t *testing.T, jDir string) {
	t.Helper()

	configPath := filepath.Join(jDir, "config.toml")

	err := os.WriteFile(configPath, []byte("[notifications]\n"), 0o644)
	require.NoError(t, err)
}

// --- Inline reminder creation via date add ---

func TestReminder_CreatedInlineWithDateAdd(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "add",
		"John Doe :: A good friend #friends $id:john_doe",
	})
	require.NoError(t, err)

	// Add a date with an inline reminder
	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "date", "add",
		"john_doe",
		"birthday :: 2000-05-13 !r[yearly 1w before] #birthday",
	})
	require.NoError(t, err)

	// Verify the reminder was persisted
	require.FileExists(t, RemindersFileExists(jDir))

	reminders := LoadReminders(t, jDir)
	require.Len(t, reminders, 1)

	r := reminders[0]
	require.Equal(t, friend.LinkedEntityDate, r.LinkedEntityType)
	require.Equal(t, "john_doe", r.FriendID)
	require.Equal(t, friend.RecurrenceYearly, r.Recurrence)
	require.Equal(t, friend.OffsetDirectionBefore, r.OffsetDirection)
	require.Equal(t, 7*24*time.Hour, r.OffsetDuration)
	require.Equal(t, friend.ReminderStatePending, r.State)
}

func TestReminder_CreatedInlineWithAbsoluteDate(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "add",
		"Jane Doe :: $id:jane_doe",
	})
	require.NoError(t, err)

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "date", "add",
		"jane_doe",
		"wedding :: 2027-06-15 !r[2027-06-01]",
	})
	require.NoError(t, err)

	reminders := LoadReminders(t, jDir)
	require.Len(t, reminders, 1)

	r := reminders[0]
	require.Equal(t, friend.RecurrenceOnce, r.Recurrence)
	require.Equal(t, time.Date(2027, 6, 1, 0, 0, 0, 0, time.UTC), r.TriggerAt)
}

func TestReminder_NotCreatedWithoutReminderSyntax(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "add",
		"John Doe :: $id:john_doe",
	})
	require.NoError(t, err)

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "date", "add",
		"john_doe",
		"birthday :: 2000-05-13 #birthday",
	})
	require.NoError(t, err)

	reminders := LoadReminders(t, jDir)
	require.Empty(t, reminders)
}

func TestReminder_DatePersistedAfterAdd(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "add",
		"John Doe :: $id:john_doe",
	})
	require.NoError(t, err)

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "date", "add",
		"john_doe",
		"birthday :: 2000-05-13",
	})
	require.NoError(t, err)

	// Verify the date was actually persisted (regression test for AddFriendDate bug)
	dates := loadDates(t, jDir)
	require.Len(t, dates, 1)
	require.NotEmpty(t, dates[0].ID)
}

// --- reminder list ---

func TestReminder_List_Empty(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"reminder", "list",
	})
	require.NoError(t, err)
}

func TestReminder_List_ShowsReminders(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	dateID := seedFriendWithDate(t, app, jDir)

	seedReminders(t, jDir, []*friend.Reminder{
		{
			ID:               "rem-aaa-00000001",
			LinkedEntityType: friend.LinkedEntityDate,
			LinkedEntityID:   dateID,
			FriendID:         "john_doe",
			TriggerAt:        time.Now().Add(24 * time.Hour),
			Recurrence:       friend.RecurrenceYearly,
			State:            friend.ReminderStatePending,
			Tags:             []string{"birthday"},
		},
		{
			ID:               "rem-bbb-00000002",
			LinkedEntityType: friend.LinkedEntityDate,
			LinkedEntityID:   dateID,
			FriendID:         "john_doe",
			TriggerAt:        time.Now().Add(-24 * time.Hour),
			Recurrence:       friend.RecurrenceOnce,
			State:            friend.ReminderStateFired,
		},
	})

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"reminder", "list",
	})
	require.NoError(t, err)
}

func TestReminder_List_FilterByState(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	dateID := seedFriendWithDate(t, app, jDir)

	seedReminders(t, jDir, []*friend.Reminder{
		{
			ID:               "rem-pending-0001",
			LinkedEntityType: friend.LinkedEntityDate,
			LinkedEntityID:   dateID,
			TriggerAt:        time.Now().Add(24 * time.Hour),
			State:            friend.ReminderStatePending,
		},
		{
			ID:               "rem-fired-00002",
			LinkedEntityType: friend.LinkedEntityDate,
			LinkedEntityID:   dateID,
			TriggerAt:        time.Now().Add(-24 * time.Hour),
			State:            friend.ReminderStateFired,
		},
	})

	// Filter by pending
	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"reminder", "list",
		"--state", "pending",
	})
	require.NoError(t, err)

	// Filter by fired
	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"reminder", "list",
		"--state", "fired",
	})
	require.NoError(t, err)
}

func TestReminder_List_FilterByType(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	dateID := seedFriendWithDate(t, app, jDir)

	seedReminders(t, jDir, []*friend.Reminder{
		{
			ID:               "rem-date-000001",
			LinkedEntityType: friend.LinkedEntityDate,
			LinkedEntityID:   dateID,
			TriggerAt:        time.Now().Add(24 * time.Hour),
			State:            friend.ReminderStatePending,
		},
	})

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"reminder", "list",
		"--type", "date",
	})
	require.NoError(t, err)
}

func TestReminder_List_FilterByTag(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	dateID := seedFriendWithDate(t, app, jDir)

	seedReminders(t, jDir, []*friend.Reminder{
		{
			ID:               "rem-tagged-00001",
			LinkedEntityType: friend.LinkedEntityDate,
			LinkedEntityID:   dateID,
			TriggerAt:        time.Now().Add(24 * time.Hour),
			State:            friend.ReminderStatePending,
			Tags:             []string{"birthday"},
		},
		{
			ID:               "rem-work-0000002",
			LinkedEntityType: friend.LinkedEntityDate,
			LinkedEntityID:   dateID,
			TriggerAt:        time.Now().Add(48 * time.Hour),
			State:            friend.ReminderStatePending,
			Tags:             []string{"work"},
		},
	})

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"reminder", "list",
		"--tag", "birthday",
	})
	require.NoError(t, err)
}

func TestReminder_List_InvalidState(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"reminder", "list",
		"--state", "bogus",
	})
	require.Error(t, err)
}

func TestReminder_List_InvalidType(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"reminder", "list",
		"--type", "bogus",
	})
	require.Error(t, err)
}

// --- reminder upcoming ---

func TestReminder_Upcoming_Empty(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"reminder", "upcoming",
	})
	require.NoError(t, err)
}

func TestReminder_Upcoming_ShowsReminders(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	dateID := seedFriendWithDate(t, app, jDir)

	seedReminders(t, jDir, []*friend.Reminder{
		{
			ID:               "rem-upcoming-001",
			LinkedEntityType: friend.LinkedEntityDate,
			LinkedEntityID:   dateID,
			FriendID:         "john_doe",
			TriggerAt:        time.Now().Add(3 * 24 * time.Hour),
			Recurrence:       friend.RecurrenceYearly,
			State:            friend.ReminderStatePending,
			Desc:             "John's birthday is coming",
		},
		{
			ID:               "rem-upcoming-002",
			LinkedEntityType: friend.LinkedEntityDate,
			LinkedEntityID:   dateID,
			FriendID:         "john_doe",
			TriggerAt:        time.Now().Add(10 * 24 * time.Hour),
			Recurrence:       friend.RecurrenceOnce,
			State:            friend.ReminderStatePending,
		},
	})

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"reminder", "upcoming",
	})
	require.NoError(t, err)
}

func TestReminder_Upcoming_WithDaysFlag(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	dateID := seedFriendWithDate(t, app, jDir)

	seedReminders(t, jDir, []*friend.Reminder{
		{
			ID:               "rem-soon-0000001",
			LinkedEntityType: friend.LinkedEntityDate,
			LinkedEntityID:   dateID,
			TriggerAt:        time.Now().Add(2 * 24 * time.Hour),
			State:            friend.ReminderStatePending,
			Desc:             "Soon reminder",
		},
		{
			ID:               "rem-later-000002",
			LinkedEntityType: friend.LinkedEntityDate,
			LinkedEntityID:   dateID,
			TriggerAt:        time.Now().Add(60 * 24 * time.Hour),
			State:            friend.ReminderStatePending,
			Desc:             "Far away reminder",
		},
	})

	// Only 7 days ahead
	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"reminder", "upcoming",
		"--days", "7",
	})
	require.NoError(t, err)
}

func TestReminder_Upcoming_ExcludesFiredReminders(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	dateID := seedFriendWithDate(t, app, jDir)

	seedReminders(t, jDir, []*friend.Reminder{
		{
			ID:               "rem-fired-00001",
			LinkedEntityType: friend.LinkedEntityDate,
			LinkedEntityID:   dateID,
			TriggerAt:        time.Now().Add(5 * 24 * time.Hour),
			State:            friend.ReminderStateFired,
			Desc:             "Already fired",
		},
	})

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"reminder", "upcoming",
	})
	require.NoError(t, err)
}

// --- reminder delete ---

func TestReminder_Delete(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	dateID := seedFriendWithDate(t, app, jDir)

	reminderID := "rem-delete-00001"
	seedReminders(t, jDir, []*friend.Reminder{
		{
			ID:               reminderID,
			LinkedEntityType: friend.LinkedEntityDate,
			LinkedEntityID:   dateID,
			TriggerAt:        time.Now().Add(24 * time.Hour),
			State:            friend.ReminderStatePending,
		},
	})

	reminders := LoadReminders(t, jDir)
	require.Len(t, reminders, 1)

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"reminder", "delete",
		reminderID,
	})
	require.NoError(t, err)

	reminders = LoadReminders(t, jDir)
	require.Empty(t, reminders)
}

func TestReminder_Delete_NonExistent(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"reminder", "delete",
		"nonexistent-id",
	})
	require.Error(t, err)
}

func TestReminder_Delete_NoArgs(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"reminder", "delete",
	})
	require.Error(t, err)
}

func TestReminder_Delete_KeepsOtherReminders(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	dateID := seedFriendWithDate(t, app, jDir)

	seedReminders(t, jDir, []*friend.Reminder{
		{
			ID:               "rem-keep-0000001",
			LinkedEntityType: friend.LinkedEntityDate,
			LinkedEntityID:   dateID,
			TriggerAt:        time.Now().Add(24 * time.Hour),
			State:            friend.ReminderStatePending,
		},
		{
			ID:               "rem-delete-00002",
			LinkedEntityType: friend.LinkedEntityDate,
			LinkedEntityID:   dateID,
			TriggerAt:        time.Now().Add(48 * time.Hour),
			State:            friend.ReminderStatePending,
		},
	})

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"reminder", "delete",
		"rem-delete-00002",
	})
	require.NoError(t, err)

	reminders := LoadReminders(t, jDir)
	require.Len(t, reminders, 1)
	require.Equal(t, "rem-keep-0000001", reminders[0].ID)
}

// --- reminder notify --dry-run ---

func TestReminder_NotifyDryRun_NoDueReminders(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	writeConfig(t, jDir)

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"reminder", "notify",
		"--dry-run",
	})
	require.NoError(t, err)
}

func TestReminder_NotifyDryRun_WithDueReminders(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	dateID := seedFriendWithDate(t, app, jDir)
	writeConfig(t, jDir)

	seedReminders(t, jDir, []*friend.Reminder{
		{
			ID:               "rem-due-00000001",
			LinkedEntityType: friend.LinkedEntityDate,
			LinkedEntityID:   dateID,
			FriendID:         "john_doe",
			TriggerAt:        time.Now().Add(-1 * time.Hour),
			Recurrence:       friend.RecurrenceOnce,
			State:            friend.ReminderStatePending,
			Desc:             "Birthday reminder",
		},
	})

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"reminder", "notify",
		"--dry-run",
	})
	require.NoError(t, err)

	// Verify the reminder state was NOT changed (dry run)
	reminders := LoadReminders(t, jDir)
	require.Len(t, reminders, 1)
	require.Equal(t, friend.ReminderStatePending, reminders[0].State)
}

func TestReminder_NotifyDryRun_SkipsNotDueReminders(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	dateID := seedFriendWithDate(t, app, jDir)
	writeConfig(t, jDir)

	seedReminders(t, jDir, []*friend.Reminder{
		{
			ID:               "rem-due-00000001",
			LinkedEntityType: friend.LinkedEntityDate,
			LinkedEntityID:   dateID,
			FriendID:         "john_doe",
			TriggerAt:        time.Now().Add(-1 * time.Hour),
			State:            friend.ReminderStatePending,
		},
		{
			ID:               "rem-future-00002",
			LinkedEntityType: friend.LinkedEntityDate,
			LinkedEntityID:   dateID,
			FriendID:         "john_doe",
			TriggerAt:        time.Now().Add(30 * 24 * time.Hour),
			State:            friend.ReminderStatePending,
		},
	})

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"reminder", "notify",
		"--dry-run",
	})
	require.NoError(t, err)

	// Both should remain pending (dry run doesn't mutate state)
	reminders := LoadReminders(t, jDir)
	require.Len(t, reminders, 2)

	for _, r := range reminders {
		require.Equal(t, friend.ReminderStatePending, r.State)
	}
}

func TestReminder_NotifyDryRun_SkipsFiredReminders(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	dateID := seedFriendWithDate(t, app, jDir)
	writeConfig(t, jDir)

	seedReminders(t, jDir, []*friend.Reminder{
		{
			ID:               "rem-already-fired",
			LinkedEntityType: friend.LinkedEntityDate,
			LinkedEntityID:   dateID,
			TriggerAt:        time.Now().Add(-1 * time.Hour),
			State:            friend.ReminderStateFired,
			FiredCount:       1,
		},
	})

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"reminder", "notify",
		"--dry-run",
	})
	require.NoError(t, err)
}

// --- Persistence ---

func TestReminder_PersistenceRoundTrip(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	dateID := seedFriendWithDate(t, app, jDir)
	now := time.Now().Truncate(time.Second)

	seedReminders(t, jDir, []*friend.Reminder{
		{
			ID:               "rem-persist-00001",
			CreatedAt:        now,
			LinkedEntityType: friend.LinkedEntityDate,
			LinkedEntityID:   dateID,
			FriendID:         "john_doe",
			TriggerAt:        now.Add(7 * 24 * time.Hour),
			ScheduleExpr:     "!r[yearly 1w before]",
			Recurrence:       friend.RecurrenceYearly,
			OffsetDirection:  friend.OffsetDirectionBefore,
			OffsetDuration:   7 * 24 * time.Hour,
			State:            friend.ReminderStatePending,
			Desc:             "Birthday reminder",
			Tags:             []string{"birthday", "important"},
		},
	})

	reminders := LoadReminders(t, jDir)
	require.Len(t, reminders, 1)

	r := reminders[0]
	require.Equal(t, "rem-persist-00001", r.ID)
	require.Equal(t, friend.LinkedEntityDate, r.LinkedEntityType)
	require.Equal(t, dateID, r.LinkedEntityID)
	require.Equal(t, "john_doe", r.FriendID)
	require.Equal(t, "!r[yearly 1w before]", r.ScheduleExpr)
	require.Equal(t, friend.RecurrenceYearly, r.Recurrence)
	require.Equal(t, friend.OffsetDirectionBefore, r.OffsetDirection)
	require.Equal(t, 7*24*time.Hour, r.OffsetDuration)
	require.Equal(t, friend.ReminderStatePending, r.State)
	require.Equal(t, "Birthday reminder", r.Desc)
	require.Equal(t, []string{"birthday", "important"}, r.Tags)
	require.Equal(t, 0, r.FiredCount)
}

func TestReminder_RemindersFileCreatedOnInit(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	require.FileExists(t, filepath.Join(jDir, file.FileNameReminders))
}
