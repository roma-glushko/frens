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

var (
	reminderRe     = regexp.MustCompile(`!r\[([^\]]+)\]`)
	offsetRe       = regexp.MustCompile(`(\d+)([dwmy])\s*(before|after)?`)
	absoluteDateRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	inFutureRe     = regexp.MustCompile(`in\s+(\d+)([dwmy])`)
)

var FormatReminderSchedule = "!r[SCHEDULE] where SCHEDULE is: DATE, Nd/w/m before|after, yearly, monthly, in Nd/w/m"

// ReminderSchedule represents a parsed reminder schedule expression
type ReminderSchedule struct {
	AbsoluteDate    *time.Time
	Recurrence      friend.Recurrence
	OffsetAmount    int
	OffsetUnit      string // d, w, m, y
	OffsetDirection friend.OffsetDirection
	InFuture        bool // true if "in Xd" format
}

// ExtractReminderSchedule extracts schedule from !r[...] syntax
func ExtractReminderSchedule(s string) (*ReminderSchedule, error) {
	matches := reminderRe.FindStringSubmatch(s)
	if len(matches) < 2 {
		return nil, fmt.Errorf("invalid reminder syntax, expected format: %s", FormatReminderSchedule)
	}

	return ParseReminderSchedule(matches[1])
}

// ParseReminderSchedule parses the content inside !r[...]
func ParseReminderSchedule(expr string) (*ReminderSchedule, error) {
	expr = strings.TrimSpace(expr)
	schedule := &ReminderSchedule{
		Recurrence: friend.RecurrenceOnce,
	}

	lowerExpr := strings.ToLower(expr)

	// Check for recurrence keywords
	if strings.Contains(lowerExpr, "yearly") {
		schedule.Recurrence = friend.RecurrenceYearly
		expr = strings.ReplaceAll(lowerExpr, "yearly", "")
		expr = strings.TrimSpace(expr)
	} else if strings.Contains(lowerExpr, "monthly") {
		schedule.Recurrence = friend.RecurrenceMonthly
		expr = strings.ReplaceAll(lowerExpr, "monthly", "")
		expr = strings.TrimSpace(expr)
	}

	// Check for "in Xd" format (future date)
	if inMatches := inFutureRe.FindStringSubmatch(expr); len(inMatches) >= 3 {
		amount, _ := strconv.Atoi(inMatches[1])
		schedule.OffsetAmount = amount
		schedule.OffsetUnit = inMatches[2]
		schedule.OffsetDirection = friend.OffsetDirectionAfter
		schedule.InFuture = true

		return schedule, nil
	}

	// Check for absolute date
	if absoluteDateRe.MatchString(expr) {
		t, err := time.Parse("2006-01-02", expr)
		if err == nil {
			schedule.AbsoluteDate = &t

			return schedule, nil
		}
	}

	// If only recurrence keyword was given, return
	if expr == "" {
		return schedule, nil
	}

	// Check for offset expression (e.g., "1w before", "3d after")
	if offsetMatches := offsetRe.FindStringSubmatch(expr); len(offsetMatches) >= 3 {
		amount, _ := strconv.Atoi(offsetMatches[1])
		schedule.OffsetAmount = amount
		schedule.OffsetUnit = offsetMatches[2]

		if len(offsetMatches) >= 4 && offsetMatches[3] != "" {
			if offsetMatches[3] == "before" {
				schedule.OffsetDirection = friend.OffsetDirectionBefore
			} else {
				schedule.OffsetDirection = friend.OffsetDirectionAfter
			}
		} else {
			// Default to "before" for relative offsets without direction
			schedule.OffsetDirection = friend.OffsetDirectionBefore
		}
	}

	return schedule, nil
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

// HasReminderExpr checks if the string contains a reminder expression
func HasReminderExpr(s string) bool {
	return reminderRe.MatchString(s)
}

// RemoveReminderExpr removes reminder expression from string
func RemoveReminderExpr(s string) string {
	return strings.TrimSpace(reminderRe.ReplaceAllString(s, ""))
}

// RenderReminderSchedule converts a Reminder back to !r[...] syntax
func RenderReminderSchedule(r *friend.Reminder) string {
	var parts []string

	if r.Recurrence != "" && r.Recurrence != friend.RecurrenceOnce {
		parts = append(parts, string(r.Recurrence))
	}

	if r.OffsetDuration > 0 {
		days := int(r.OffsetDuration.Hours() / 24)

		var offsetStr string
		if days >= 365 && days%365 == 0 {
			offsetStr = fmt.Sprintf("%dy", days/365)
		} else if days >= 30 && days%30 == 0 {
			offsetStr = fmt.Sprintf("%dm", days/30)
		} else if days >= 7 && days%7 == 0 {
			offsetStr = fmt.Sprintf("%dw", days/7)
		} else {
			offsetStr = fmt.Sprintf("%dd", days)
		}

		parts = append(parts, offsetStr)

		if r.OffsetDirection != "" {
			parts = append(parts, string(r.OffsetDirection))
		}
	}

	if len(parts) == 0 && !r.TriggerAt.IsZero() {
		return fmt.Sprintf("!r[%s]", r.TriggerAt.Format("2006-01-02"))
	}

	return fmt.Sprintf("!r[%s]", strings.Join(parts, " "))
}
