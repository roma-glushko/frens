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

package friend

import (
	"fmt"
	"strings"
	"time"
)

type SortOption string

const (
	SortAlpha      SortOption = "alpha"
	SortActivities SortOption = "activities"
	SortRecency    SortOption = "recency"
)

var SortOptions = []SortOption{
	SortAlpha,
	SortActivities,
	SortRecency,
}

func ValidateSortOption(s string) error {
	validOpts := make([]string, 0, len(SortOptions))

	for _, sortOpt := range SortOptions {
		opt := string(sortOpt)

		validOpts = append(validOpts, opt)

		if s == opt {
			return nil
		}
	}

	return fmt.Errorf(
		"invalid sort value '%s' (valid: %s)",
		s,
		strings.Join(validOpts, ", "),
	)
}

type SortOrderOption string

const (
	SortOrderDirect  SortOrderOption = "direct"
	SortOrderReverse SortOrderOption = "reverse"
)

type ListFriendQuery struct {
	Keyword   string
	Locations []string
	Tags      []string
	SortBy    SortOption
	SortOrder SortOrderOption
}

type ListLocationQuery struct {
	Keyword   string
	Countries []string
	Tags      []string
	SortBy    SortOption
	SortOrder SortOrderOption
}

type ListEventQuery struct {
	Type         EventType
	Keyword      string
	Locations    []string
	Tags         []string
	Since, Until time.Time
	SortBy       SortOption
	SortOrder    SortOrderOption
}
