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
	"github.com/roma-glushko/frens/internal/journaldir"
	"github.com/urfave/cli/v2"
)

var EditCommand = &cli.Command{
	Name:      "edit",
	Aliases:   []string{"e", "modify", "update"},
	Usage:     "Update activity log",
	Args:      true,
	ArgsUsage: `<FRIEND_NAME, FRIEND_NICKNAME, FRIEND_ID>`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "name",
			Aliases: []string{"n"},
			Usage:   "Set friend's name",
		},
		&cli.StringFlag{
			Name:    "desc",
			Aliases: []string{"d"},
			Usage:   "Set description of the friend",
		},
	},
	Action: func(ctx *cli.Context) error {
		journalDir, err := journaldir.DefaultDir()
		if err != nil {
			return err
		}

		_, err = journaldir.Load(journalDir)
		if err != nil {
			return err
		}

		// TODO: implement activity editing

		return nil
	},
}
