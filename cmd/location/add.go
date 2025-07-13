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
	"errors"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/journaldir"
	"github.com/roma-glushko/frens/internal/lang"
	"github.com/roma-glushko/frens/internal/tui"
	"github.com/urfave/cli/v2"
)

var AddCommand = &cli.Command{
	Name:    "add",
	Aliases: []string{"a", "new", "create"},
	Usage:   "Add a new location",
	Args:    true,
	ArgsUsage: `<INFO>
		If no arguments are provided, a textarea will be shown to fill in the details interactively.
		Otherwise, the information will be parsed from the command options.
		
		<INFO> format:
			` + lang.FormatLocationInfo + `

		For example:
			Scranton, USA (a.k.a. The Electric City, Scranton) :: Located a branch of Dunder Mifflin #theoffice
			New York City (aka NYC, The Big Apple) :: A bustling metropolis known for its skyscrapers and culture
	`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "name",
			Aliases: []string{"n"},
			Usage:   "Name of the location",
		},
		&cli.StringFlag{
			Name:    "country",
			Aliases: []string{"cn"},
			Usage:   "Country of the location",
		},
		&cli.StringFlag{
			Name:    "desc",
			Aliases: []string{"d"},
			Usage:   "Description of the location",
		},
		&cli.StringSliceFlag{
			Name:    "alias",
			Aliases: []string{"a", "aka", "nick"},
			Usage:   "Aliases for the location (can be used to search for it or refer to it)",
		},
		&cli.StringSliceFlag{
			Name:    "tags",
			Aliases: []string{"t"},
			Usage:   "Tags associated with the location (for categorization or search purposes)",
		},
	},
	Action: func(ctx *cli.Context) error {
		journalDir, err := journaldir.DefaultDir()
		if err != nil {
			return err
		}

		j, err := journaldir.Load(journalDir)
		if err != nil {
			return err
		}

		var info string

		if ctx.NArg() == 0 {
			// TODO: also check if we are in the interactive mode
			inputForm := tui.NewEditorForm(tui.EditorOptions{
				Title:      "Add a new location information:",
				SyntaxHint: lang.FormatLocationInfo,
			})
			teaUI := tea.NewProgram(inputForm, tea.WithMouseAllMotion())

			if _, err := teaUI.Run(); err != nil {
				log.Error("uh oh", "err", err)
				return err
			}

			info = inputForm.Textarea.Value()
		} else {
			info = strings.Join(ctx.Args().Slice(), " ")
		}

		var l friend.Location

		if info != "" {
			l, err = lang.ExtractLocation(info)

			if err != nil && !errors.Is(err, lang.ErrNoInfo) {
				log.Error("failed to parse location info", "err", err)
				return err
			}
		}

		// apply CLI flags
		name := ctx.String("name")
		country := ctx.String("country")
		desc := ctx.String("desc")
		aliases := ctx.StringSlice("alias")
		tags := ctx.StringSlice("tag")

		if name != "" {
			l.Name = name
		}

		if country != "" {
			l.Country = country
		}

		if desc != "" {
			l.Desc = desc
		}

		if len(aliases) > 0 {
			l.Aliases = aliases
		}

		if len(tags) > 0 {
			l.Tags = tags
		}

		if err := l.Validate(); err != nil {
			return err
		}

		err = journaldir.Update(j, func(j *journal.Data) error {
			j.AddLocation(l)

			return nil
		})
		if err != nil {
			return err
		}

		log.Info("✅ Added location: " + l.String())

		if len(l.Aliases) > 0 {
			log.Info("📍 Aliases: " + strings.Join(l.Aliases, ", "))
		}

		if len(l.Tags) > 0 {
			log.Info("🏷️ Tags: " + strings.Join(l.Tags, ", "))
		}

		if l.Desc != "" {
			log.Info("🧭 Description: \n" + l.Desc)
		}

		return nil
	},
}
