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

package lang

import (
	"testing"
	"time"

	"github.com/roma-glushko/frens/internal/friend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractReminderSchedule(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		wantRecurrence friend.Recurrence
		wantOffsetAmt  int
		wantOffsetUnit string
		wantOffsetDir  friend.OffsetDirection
		wantAbsolute   bool
		wantInFuture   bool
		wantErr        bool
	}{
		{
			name:           "absolute date",
			input:          "!r[2025-03-15]",
			wantRecurrence: friend.RecurrenceOnce,
			wantAbsolute:   true,
		},
		{
			name:           "1 week before",
			input:          "!r[1w before]",
			wantRecurrence: friend.RecurrenceOnce,
			wantOffsetAmt:  1,
			wantOffsetUnit: "w",
			wantOffsetDir:  friend.OffsetDirectionBefore,
		},
		{
			name:           "3 days after",
			input:          "!r[3d after]",
			wantRecurrence: friend.RecurrenceOnce,
			wantOffsetAmt:  3,
			wantOffsetUnit: "d",
			wantOffsetDir:  friend.OffsetDirectionAfter,
		},
		{
			name:           "yearly recurrence",
			input:          "!r[yearly]",
			wantRecurrence: friend.RecurrenceYearly,
		},
		{
			name:           "monthly recurrence",
			input:          "!r[monthly]",
			wantRecurrence: friend.RecurrenceMonthly,
		},
		{
			name:           "yearly 1 week before",
			input:          "!r[yearly 1w before]",
			wantRecurrence: friend.RecurrenceYearly,
			wantOffsetAmt:  1,
			wantOffsetUnit: "w",
			wantOffsetDir:  friend.OffsetDirectionBefore,
		},
		{
			name:           "in 3 days",
			input:          "!r[in 3d]",
			wantRecurrence: friend.RecurrenceOnce,
			wantOffsetAmt:  3,
			wantOffsetUnit: "d",
			wantOffsetDir:  friend.OffsetDirectionAfter,
			wantInFuture:   true,
		},
		{
			name:           "in 2 weeks",
			input:          "!r[in 2w]",
			wantRecurrence: friend.RecurrenceOnce,
			wantOffsetAmt:  2,
			wantOffsetUnit: "w",
			wantOffsetDir:  friend.OffsetDirectionAfter,
			wantInFuture:   true,
		},
		{
			name:           "default direction is before",
			input:          "!r[2d]",
			wantRecurrence: friend.RecurrenceOnce,
			wantOffsetAmt:  2,
			wantOffsetUnit: "d",
			wantOffsetDir:  friend.OffsetDirectionBefore,
		},
		{
			name:    "invalid syntax - no brackets",
			input:   "!r",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schedule, err := ExtractReminderSchedule(tt.input)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantRecurrence, schedule.Recurrence)
			assert.Equal(t, tt.wantOffsetAmt, schedule.OffsetAmount)
			assert.Equal(t, tt.wantOffsetUnit, schedule.OffsetUnit)
			assert.Equal(t, tt.wantOffsetDir, schedule.OffsetDirection)
			assert.Equal(t, tt.wantInFuture, schedule.InFuture)

			if tt.wantAbsolute {
				assert.NotNil(t, schedule.AbsoluteDate)
			}
		})
	}
}

func TestReminderSchedule_ToDuration(t *testing.T) {
	tests := []struct {
		name     string
		schedule ReminderSchedule
		want     time.Duration
	}{
		{
			name:     "1 day",
			schedule: ReminderSchedule{OffsetAmount: 1, OffsetUnit: "d"},
			want:     24 * time.Hour,
		},
		{
			name:     "2 weeks",
			schedule: ReminderSchedule{OffsetAmount: 2, OffsetUnit: "w"},
			want:     14 * 24 * time.Hour,
		},
		{
			name:     "1 month",
			schedule: ReminderSchedule{OffsetAmount: 1, OffsetUnit: "m"},
			want:     30 * 24 * time.Hour,
		},
		{
			name:     "1 year",
			schedule: ReminderSchedule{OffsetAmount: 1, OffsetUnit: "y"},
			want:     365 * 24 * time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.schedule.ToDuration()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestReminderSchedule_ComputeTriggerDate(t *testing.T) {
	baseDate := time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC)
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		schedule ReminderSchedule
		want     time.Time
	}{
		{
			name: "absolute date",
			schedule: ReminderSchedule{
				AbsoluteDate: func() *time.Time { t := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC); return &t }(),
			},
			want: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "1 week before",
			schedule: ReminderSchedule{
				OffsetAmount:    1,
				OffsetUnit:      "w",
				OffsetDirection: friend.OffsetDirectionBefore,
			},
			want: time.Date(2025, 3, 8, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "3 days after",
			schedule: ReminderSchedule{
				OffsetAmount:    3,
				OffsetUnit:      "d",
				OffsetDirection: friend.OffsetDirectionAfter,
			},
			want: time.Date(2025, 3, 18, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "in 3 days from now",
			schedule: ReminderSchedule{
				OffsetAmount:    3,
				OffsetUnit:      "d",
				OffsetDirection: friend.OffsetDirectionAfter,
				InFuture:        true,
			},
			want: time.Date(2025, 1, 4, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.schedule.ComputeTriggerDate(baseDate, now)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHasReminderExpr(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"!r[yearly]", true},
		{"birthday !r[yearly 1w before]", true},
		{"no reminder here", false},
		{"!r[]", false}, // Empty brackets is not valid
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := HasReminderExpr(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRemoveReminderExpr(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"birthday !r[yearly]", "birthday"},
		{"!r[1w before] anniversary", "anniversary"},
		{"no reminder here", "no reminder here"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := RemoveReminderExpr(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRenderReminderSchedule(t *testing.T) {
	tests := []struct {
		name     string
		reminder *friend.Reminder
		want     string
	}{
		{
			name: "yearly recurrence",
			reminder: &friend.Reminder{
				Recurrence: friend.RecurrenceYearly,
			},
			want: "!r[yearly]",
		},
		{
			name: "monthly with offset",
			reminder: &friend.Reminder{
				Recurrence:      friend.RecurrenceMonthly,
				OffsetDuration:  7 * 24 * time.Hour,
				OffsetDirection: friend.OffsetDirectionBefore,
			},
			want: "!r[monthly 1w before]",
		},
		{
			name: "absolute date",
			reminder: &friend.Reminder{
				TriggerAt: time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC),
			},
			want: "!r[2025-03-15]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RenderReminderSchedule(tt.reminder)
			assert.Equal(t, tt.want, got)
		})
	}
}
