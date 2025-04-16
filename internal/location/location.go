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

package location

import (
	"fmt"
	"strings"
)

type Location struct {
	Name    string   `toml:"name"`
	Country string   `toml:"country,omitempty"`
	Alias   []string `toml:"alias,omitempty"`
	Tags    []string `toml:"tags,omitempty"`
}

func (l Location) String() string {
	var sb strings.Builder

	sb.WriteString(l.Name)

	if len(l.Country) > 0 {
		sb.WriteString(", " + l.Country)
	}

	if len(l.Alias) > 0 {
		sb.WriteString(fmt.Sprintf(" (%s)", strings.Join(l.Alias, ", ")))
	}

	return sb.String()
}
