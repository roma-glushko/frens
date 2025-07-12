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
	"github.com/charmbracelet/log"
	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/journaldir"
	"github.com/urfave/cli/v2"
)

var AddCommand = &cli.Command{
	Name:    "add",
	Aliases: []string{"a", "new", "create"},
	Usage:   "Add a new activity",
	Args:    true,
	ArgsUsage: `<DESCR> [<DESCR2> ...]

	<DESCR> is a description of the activity to record.
	
	Examples:
		"Michael wrote a book 'Somehow I managed'" - no date, will be recorded as today
		"yesterday: Jim Halpert put my stuff in jello" - relative date & description
		"2009/09/08: "Jim and Pam got married at Niagara Falls" - absolute date & description
`,
	Action: func(ctx *cli.Context) error {
		if ctx.NArg() == 0 {
			return cli.Exit("You must provide a description for the activity.", 1)
		}

		journalDir, err := journaldir.DefaultDir()
		if err != nil {
			return err
		}

		j, err := journaldir.Load(journalDir)
		if err != nil {
			return err
		}

		log.Info("Adding a new activity..")

		for _, desc := range ctx.Args().Slice() {
			if desc == "" {
				log.Warn("Empty description provided, skipping.")
				continue
			}

			e := friend.NewEvent(friend.EventTypeActivity, desc)

			j.AddActivity(e)
		}

		return nil
	},
}
