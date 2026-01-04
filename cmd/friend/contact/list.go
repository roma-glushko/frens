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

package contact

import (
	"fmt"

	jctx "github.com/roma-glushko/frens/internal/context"
	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/log"
	"github.com/roma-glushko/frens/internal/log/formatter"
	"github.com/urfave/cli/v2"
)

var ListCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"l", "ls"},
	Usage:   "List contacts for all friends",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "search",
			Aliases: []string{"q"},
			Usage:   "Search by value or label",
		},
		&cli.StringSliceFlag{
			Name:    "with",
			Aliases: []string{"w"},
			Usage:   "Filter by friend(s)",
		},
		&cli.StringSliceFlag{
			Name:    "type",
			Aliases: []string{"tp"},
			Usage:   "Filter by contact type(s) (email, phone, telegram, etc.)",
		},
		&cli.StringSliceFlag{
			Name:    "tag",
			Aliases: []string{"t"},
			Usage:   "Filter by tag(s)",
		},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context
		jctx := jctx.FromCtx(ctx)
		s := jctx.Store

		return s.Tx(ctx, func(j *journal.Journal) error {
			// Parse contact types
			typeStrs := c.StringSlice("type")
			types := make([]friend.ContactType, 0, len(typeStrs))
			for _, t := range typeStrs {
				types = append(types, friend.ParseContactType(t))
			}

			contacts, err := j.ListFriendContacts(friend.ListContactQuery{
				Keyword: c.String("search"),
				Friends: c.StringSlice("with"),
				Types:   types,
				Tags:    c.StringSlice("tag"),
			})
			if err != nil {
				return err
			}

			if len(contacts) == 0 {
				log.Info("No contacts found for given query.")
				return nil
			}

			fmtr := formatter.ContactTextFormatter{}

			o, _ := fmtr.FormatList(contacts)
			fmt.Println(o)

			return nil
		})
	},
}
