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
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/roma-glushko/frens/internal/friend"
)

// Syntax: !r[SCHEDULE]
// Examples:
//   !r[2025-03-15]          - Absolute date
//   !r[1w before]           - 1 week before linked date
//   !r[3d after]            - 3 days after linked date
//   !r[yearly]              - Yearly on the linked date
//   !r[monthly]             - Monthly on the linked date
//   !r[yearly 1w before]    - Yearly, 1 week before
//   !r[in 3d]               - 3 days from now
//   !r[Friday]              - Next Friday
//   !r[monday]              - Next Monday

const (
	markerReminderStart = "!r["
	markerReminderEnd   = "]"
)

var (
	reminderRe     = regexp.MustCompile(`!r\[([^\]]+)\]`)
	offsetRe       = regexp.MustCompile(`(\d+)([dwmy])\s*(before|after)?`)
	absoluteDateRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	inFutureRe     = regexp.MustCompile(`in\s+(\d+)([dwmy])`)
)

var dayOfWeekMap = map[string]time.Weekday{
	"sunday":    time.Sunday,
	"sun":       time.Sunday,
	"monday":    time.Monday,
	"mon":       time.Monday,
	"tuesday":   time.Tuesday,
	"tue":       time.Tuesday,
	"wednesday": time.Wednesday,
	"wed":       time.Wednesday,
	"thursday":  time.Thursday,
	"thu":       time.Thursday,
	"friday":    time.Friday,
	"fri":       time.Friday,
	"saturday":  time.Saturday,
	"sat":       time.Saturday,
}

var FormatReminderSchedule = fmt.Sprintf(
	"%sSCHEDULE%s where SCHEDULE is: DATE, WEEKDAY, Nd/w/m before|after, yearly, monthly, in Nd/w/m",
	markerReminderStart,
	markerReminderEnd,
)

// ReminderSchedule represents a parsed reminder schedule expression
type ReminderSchedule struct {
	AbsoluteDate    *time.Time
	Weekday         *time.Weekday // for "Friday", "Monday", etc.
	Recurrence      friend.Recurrence
	OffsetAmount    int
	OffsetUnit      string // d, w, m, y
	OffsetDirection friend.OffsetDirection
	InFuture        bool // true if "in Xd" format
}

// ParseReminderSchedule parses the content inside !r[...]
func ParseReminderSchedule(s string) (*ReminderSchedule, error) {
	matches := reminderRe.FindStringSubmatch(s)

	if len(matches) < 2 {
		return nil, fmt.Errorf(
			"invalid reminder syntax, expected format: %s",
			FormatReminderSchedule,
		)
	}

	expr := strings.TrimSpace(matches[1])
	recurrence, expr := extractRecurrence(expr)

	schedule := &ReminderSchedule{
		Recurrence: recurrence,
	}

	// Check for "in Xd" format (future date)
	if amount, unit, ok := extractInFuture(expr); ok {
		schedule.OffsetAmount = amount
		schedule.OffsetUnit = unit
		schedule.OffsetDirection = friend.OffsetDirectionAfter
		schedule.InFuture = true

		return schedule, nil
	}

	// Check for absolute date
	if date := extractAbsoluteDate(expr); date != nil {
		schedule.AbsoluteDate = date

		return schedule, nil
	}

	// Check for day of week (e.g., "Friday", "mon")
	if weekday := extractDayOfWeek(expr); weekday != nil {
		schedule.Weekday = weekday

		return schedule, nil
	}

	// If only recurrence keyword was given, return
	if expr == "" {
		return schedule, nil
	}

	// Check for offset expression (e.g., "1w before", "3d after")
	amount, unit, direction, ok := extractOffset(expr)
	if !ok {
		return nil, fmt.Errorf(
			"invalid reminder offset syntax, expected format: %s",
			FormatReminderSchedule,
		)
	}

	schedule.OffsetAmount = amount
	schedule.OffsetUnit = unit
	schedule.OffsetDirection = direction

	return schedule, nil
}

// extractRecurrence extracts recurrence keyword and returns (recurrence, remainingExpr)
func extractRecurrence(expr string) (friend.Recurrence, string) {
	lowerExpr := strings.ToLower(expr)

	if strings.Contains(lowerExpr, "yearly") {
		return friend.RecurrenceYearly, strings.TrimSpace(
			strings.ReplaceAll(lowerExpr, "yearly", ""),
		)
	}

	if strings.Contains(lowerExpr, "monthly") {
		return friend.RecurrenceMonthly, strings.TrimSpace(
			strings.ReplaceAll(lowerExpr, "monthly", ""),
		)
	}

	return friend.RecurrenceOnce, expr
}

// extractInFuture tries to parse "in Xd" format, returns (amount, unit, ok)
func extractInFuture(expr string) (int, string, bool) {
	matches := inFutureRe.FindStringSubmatch(expr)
	if len(matches) < 3 {
		return 0, "", false
	}

	amount, _ := strconv.Atoi(matches[1])

	return amount, matches[2], true
}

// extractAbsoluteDate tries to parse absolute date, returns nil if not matched
func extractAbsoluteDate(expr string) *time.Time {
	if !absoluteDateRe.MatchString(expr) {
		return nil
	}

	t, err := time.Parse("2006-01-02", expr)
	if err != nil {
		return nil
	}

	return &t
}

// extractDayOfWeek tries to parse day of week, returns nil if not matched
func extractDayOfWeek(expr string) *time.Weekday {
	weekday, ok := dayOfWeekMap[strings.ToLower(expr)]
	if !ok {
		return nil
	}

	return &weekday
}

// extractOffset tries to parse offset expression (e.g., "1w before"), returns (amount, unit, direction, ok)
func extractOffset(expr string) (int, string, friend.OffsetDirection, bool) {
	matches := offsetRe.FindStringSubmatch(expr)
	if len(matches) < 3 {
		return 0, "", "", false
	}

	amount, _ := strconv.Atoi(matches[1])
	unit := matches[2]
	direction := friend.OffsetDirectionBefore

	if len(matches) >= 4 && matches[3] == "after" {
		direction = friend.OffsetDirectionAfter
	}

	return amount, unit, direction, true
}

// ToDuration converts the offset to a time.Duration
func (s *ReminderSchedule) ToDuration() time.Duration {
	multiplier := time.Duration(s.OffsetAmount)

	switch s.OffsetUnit {
	case "d":
		return multiplier * 24 * time.Hour
	case "w":
		return multiplier * 7 * 24 * time.Hour
	case "m":
		return multiplier * 30 * 24 * time.Hour // Approximate
	case "y":
		return multiplier * 365 * 24 * time.Hour // Approximate
	default:
		return 0
	}
}

// ComputeTriggerDate computes the actual trigger date based on schedule and base date
func (s *ReminderSchedule) ComputeTriggerDate(baseDate time.Time, now time.Time) time.Time {
	// For absolute dates
	if s.AbsoluteDate != nil {
		return *s.AbsoluteDate
	}

	// For day of week (e.g., "Friday"), compute next occurrence from now
	if s.Weekday != nil {
		return nextWeekday(now, *s.Weekday)
	}

	// For "in Xd" format, compute from now
	if s.InFuture {
		return now.Add(s.ToDuration())
	}

	// For relative offsets from a base date
	offset := s.ToDuration()
	if s.OffsetDirection == friend.OffsetDirectionBefore {
		return baseDate.Add(-offset)
	}

	return baseDate.Add(offset)
}

// nextWeekday returns the next occurrence of the given weekday from the given time.
// If today is the target weekday, it returns 7 days from now (next week).
func nextWeekday(from time.Time, target time.Weekday) time.Time {
	current := from.Weekday()
	daysUntil := int(target - current)

	if daysUntil <= 0 {
		daysUntil += 7
	}

	return from.AddDate(0, 0, daysUntil)
}

// ExtractReminder extracts !r[schedule] from text and creates a Reminder struct.
// Returns (nil, nil) if no reminder expression is found.
// Returns (nil, error) if parsing fails.
// Returns (reminder, nil) on success.
func ExtractReminder(
	text string,
	entityType friend.LinkedEntityType,
	entityID, friendID string,
	baseDate time.Time,
	tags []string,
) (*friend.Reminder, error) {
	if !reminderRe.MatchString(text) {
		return nil, nil
	}

	schedule, err := ParseReminderSchedule(text)
	if err != nil {
		return nil, err
	}

	triggerAt := schedule.ComputeTriggerDate(baseDate, time.Now())

	return &friend.Reminder{
		LinkedEntityType: entityType,
		LinkedEntityID:   entityID,
		FriendID:         friendID,
		TriggerAt:        triggerAt,
		ScheduleExpr:     text,
		Recurrence:       schedule.Recurrence,
		OffsetDirection:  schedule.OffsetDirection,
		OffsetDuration:   schedule.ToDuration(),
		State:            friend.ReminderStatePending,
		Tags:             tags,
	}, nil
}

// RemoveReminderExpr removes reminder expression from string
func RemoveReminderExpr(s string) string {
	return strings.TrimSpace(reminderRe.ReplaceAllString(s, ""))
}

// RenderReminder converts a Reminder back to !r[...] syntax
func RenderReminder(r *friend.Reminder) string {
	var parts []string

	if r.Recurrence != "" && r.Recurrence != friend.RecurrenceOnce {
		parts = append(parts, r.Recurrence)
	}

	if r.OffsetDuration > 0 {
		offsetStr := renderOffset(r.OffsetDuration)
		parts = append(parts, offsetStr)

		if r.OffsetDirection != "" {
			parts = append(parts, r.OffsetDirection)
		}
	}

	if len(parts) == 0 && !r.TriggerAt.IsZero() {
		return markerReminderStart + r.TriggerAt.Format("2006-01-02") + markerReminderEnd
	}

	return markerReminderStart + strings.Join(parts, " ") + markerReminderEnd
}

func renderOffset(d time.Duration) string {
	days := int(d.Hours() / 24)

	switch {
	case days >= 365 && days%365 == 0:
		return strconv.Itoa(days/365) + "y"
	case days >= 30 && days%30 == 0:
		return strconv.Itoa(days/30) + "m"
	case days >= 7 && days%7 == 0:
		return strconv.Itoa(days/7) + "w"
	default:
		return strconv.Itoa(days) + "d"
	}
}
