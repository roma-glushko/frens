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

package activity

import (
	"github.com/urfave/cli/v2"
)

var Commands = &cli.Command{
	Name:        "activity",
	Aliases:     []string{"a", "act"},
	Usage:       "Manage your activity logs",
	Description: `Activity logs are events, interactions with people or locations. They may optionally contain any number of friends, locations, and tags.`,
	Subcommands: []*cli.Command{
		AddCommand,
		EditCommand,
		ListCommand,
		DeleteCommand,
	},
}
