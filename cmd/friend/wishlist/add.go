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

package wishlist

import (
	"errors"
	"fmt"
	"strings"

	"github.com/roma-glushko/frens/internal/wishlist"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/roma-glushko/frens/internal/tui"

	jctx "github.com/roma-glushko/frens/internal/context"

	"github.com/roma-glushko/frens/internal/friend"
	//"github.com/roma-glushko/frens/internal/journal"
	//"github.com/roma-glushko/frens/internal/journaldir"
	"github.com/roma-glushko/frens/internal/lang"
	"github.com/roma-glushko/frens/internal/log"
	"github.com/urfave/cli/v2"
)

var AddCommand = &cli.Command{
	Name:      "add",
	Aliases:   []string{"a", "new", "create"},
	Usage:     "Add a new item to friend's wishlist",
	UsageText: "frens friend wishlist add [OPTIONS] <FRIEND_NAME, FRIEND_NICKNAME, FRIEND_ID> [INFO]",
	Args:      true,
	ArgsUsage: `<INFO>
		If no arguments are provided, a textarea will be shown to fill in the details interactively.
		Otherwise, the information will be parsed from the command options.
		
		<INFO> format:
			` + lang.FormatWishlistItem + `

		For example:
			"https://amazon.com/cool-keyboard #techgift"
			"Cool mouse #tech #gaming"
	`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "desc",
			Aliases: []string{"d"},
			Usage:   "Description of the wishlist item",
		},
		&cli.StringFlag{
			Name:    "link",
			Aliases: []string{"l"},
			Usage:   "Link to the wishlist item (e.g., product page URL)",
		},
		&cli.StringFlag{
			Name:    "price",
			Aliases: []string{"p"},
			Usage:   "Price of the wishlist item (e.g., 29.99USD)",
		},
		&cli.StringSliceFlag{
			Name:    "tag",
			Aliases: []string{"t"},
			Usage:   "Add tags to the date (e.g., 'gift', 'tech', 'gaming')",
		},
	},
	Action: func(ctx *cli.Context) error {
		var info string

		if ctx.NArg() == 0 {
			return cli.Exit(
				"You must provide a friend name, nickname, or ID to add a wishlist item.",
				1,
			)
		}

		appCtx := jctx.FromCtx(ctx.Context)
		jr := appCtx.Repository.Journal()

		pID := ctx.Args().First()
		_, err := jr.GetFriend(pID)
		if err != nil {
			return err
		}

		if ctx.NArg() == 1 {
			// TODO: also check if we are in the interactive mode
			inputForm := tui.NewEditorForm(tui.EditorOptions{
				Title:      "Add a new friend wishlist item information:",
				SyntaxHint: lang.FormatWishlistItem,
			})
			teaUI := tea.NewProgram(inputForm, tea.WithMouseAllMotion())

			if _, err := teaUI.Run(); err != nil {
				log.Errorf("uh oh: %v", err)
				return err
			}

			info = inputForm.Textarea.Value()
		} else {
			info = strings.Join(ctx.Args().Slice()[1:], " ")
		}

		var w friend.WishlistItem

		if info != "" {
			w, err = lang.ExtractWishlistItem(info)

			if err != nil && !errors.Is(err, lang.ErrNoInfo) {
				log.Errorf("failed to parse date info: %v", err)
				return err
			}
		}

		// apply CLI flags
		desc := ctx.String("desc")
		link := ctx.String("link")
		tags := ctx.StringSlice("tag")

		if desc != "" {
			w.Desc = desc
		}

		if link != "" {
			w.Link = link
		}

		if len(tags) > 0 {
			w.Tags = tags
		}

		if err := w.Validate(); err != nil {
			return err
		}

		urlMan := wishlist.NewURLManager()

		if w.Link != "" {
			pInfo, err := urlMan.Scrape(ctx.Context, w.Link)
			if err != nil {
				return fmt.Errorf("failed to scrape product info: %v", err)
			}

			fmt.Printf("%+v", pInfo)
		}

		//err = journaldir.Update(jr, func(j *journal.Journal) error {
		//	d, err = j.AddFriendDate(p.ID, d)
		//
		//	return err
		//})
		//if err != nil {
		//	return err
		//}
		//
		log.Info(" ✔ Wishlist item added")

		return nil
	},
}
