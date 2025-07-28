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
	"time"

	"github.com/markusmobius/go-dateparser"
)

var DescSeparator = "::"

var (
	FormatDateInfo = fmt.Sprintf(
		"[LABEL ::] DATE [$cal:gregorian|hebrew]",
	)
)

// ExtractDate parses a free-form & relative date string and returns a time.Time object.
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

//func ExtractDateInfo(s string) (string, time.Time) {
//
//}
