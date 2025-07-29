// Copyright 2025 Roma Hlushko
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
	"strings"
	"time"

	"github.com/roma-glushko/frens/internal/friend"

	"github.com/markusmobius/go-dateparser"
)

var Separator = "::"

var FormatDateInfo = "[LABEL ::] DATE [$cal:gregorian|hebrew]"

// ExtractDate parses a free-form & relative date string and returns a time.Time object.
// If some parts of the timestamp are missing, it will assume the current year, month, and day, etc.
func ExtractDate(s string, def ...time.Time) time.Time {
	ts := time.Time{}

	if len(def) > 0 {
		ts = def[0].UTC()
	}

	if s != "" {
		parsedDate, err := dateparser.Parse(nil, s)
		if err != nil {
			ts = time.Now().UTC()
		} else {
			ts = parsedDate.Time.UTC()
		}
	}

	return ts
}

type dateProps struct {
	Calendar string `frentxt:"cal"`
}

func ParseExactDate(s string) (time.Time, error) {
	timeMarker := time.Time{}
	c := dateparser.Configuration{
		CurrentTime: timeMarker,
	} // empty date as default, so we can infer which parts of the date are missing

	parsedDate, err := dateparser.Parse(&c, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse date: %w", err)
	}

	return parsedDate.Time, nil
}

func ExtractDateInfo(s string) (*friend.Date, error) {
	props, err := ExtractProps[dateProps](s)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date properties: %w", err)
	}

	s = RemoveProps(s)

	parts := strings.SplitN(s, Separator, 2)

	var label string

	var dateExpr string

	if len(parts) == 0 {
		return nil, fmt.Errorf("invalid date info format, expected format: %s", FormatDateInfo)
	}

	if len(parts) == 2 {
		label = strings.TrimSpace(parts[0])
		dateExpr = strings.TrimSpace(parts[1])
	} else {
		dateExpr = strings.TrimSpace(parts[0])
	}

	cal := friend.CalendarGregorian

	if props.Calendar != "" { // TODO: extract validation into a separate function
		switch strings.ToLower(props.Calendar) {
		case friend.CalendarGregorian:
			cal = friend.CalendarGregorian
		case friend.CalendarHebrew:
			cal = friend.CalendarHebrew
		default:
			return nil, fmt.Errorf(
				"unsupported calendar type: %s",
				props.Calendar,
			) // TODO: list available calendars
		}
	}

	return &friend.Date{
		Label:    label,
		DateExpr: dateExpr,
		Calendar: cal,
	}, nil
}
