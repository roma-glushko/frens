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

package note

import "github.com/urfave/cli/v2"

var Commands = &cli.Command{
	Name:        "note",
	Aliases:     []string{"n", "nt"},
	Usage:       "Manage your notes",
	UsageText:   "frens note [command] [options]",
	Description: `Notes helps to remember things about friends and locations with deeper meaning, insights, background, preferences, longer-term context.`,
	Subcommands: []*cli.Command{
		AddCommand,
		EditCommand,
		ListCommand,
		DeleteCommand,
	},
}
