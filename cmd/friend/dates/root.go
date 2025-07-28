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

package dates

import (
	"github.com/urfave/cli/v2"
)

var Commands = &cli.Command{
	Name:        "dates",
	Aliases:     []string{"dt", "date"},
	Usage:       "Manage your friend's important dates",
	UsageText:   "frens friend dates [command] [options]",
	Description: `Important dates lets you keep track of your friend's birthdays, anniversaries, and other significant events.`,
	Subcommands: []*cli.Command{
		AddCommand,
	},
}
