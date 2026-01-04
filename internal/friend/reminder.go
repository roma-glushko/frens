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

import "time"

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

type Reminder struct {
	Recurrence      Recurrence      `toml:"recurrence,omitempty"`
	OffsetDirection OffsetDirection `toml:"offset_direction,omitempty"`
	Offset          *time.Duration  `toml:"offset,omitempty"`
}
