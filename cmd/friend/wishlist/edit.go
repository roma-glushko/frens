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

	"github.com/roma-glushko/frens/internal/log/formatter"

	jctx "github.com/roma-glushko/frens/internal/context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/lang"
	"github.com/roma-glushko/frens/internal/log"
	"github.com/roma-glushko/frens/internal/tui"
	"github.com/urfave/cli/v2"
)

var EditCommand = &cli.Command{
	Name:      "edit",
	Aliases:   []string{"e", "modify", "update"},
	Usage:     "Update a wishlist item",
	Args:      true,
	ArgsUsage: `<WISHLIST_ITEM_ID>`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "desc",
			Aliases: []string{"d"},
			Usage:   "Description of the wishlist item",
		},
		&cli.StringFlag{
			Name:    "link",
			Aliases: []string{"l"},
			Usage:   "Link to the wishlist item",
		},
		&cli.StringFlag{
			Name:    "price",
			Aliases: []string{"p"},
			Usage:   "Price of the wishlist item",
		},
		&cli.StringSliceFlag{
			Name:    "tag",
			Aliases: []string{"t"},
			Usage:   "Set tags for the wishlist item",
		},
	},
	Action: func(ctx *cli.Context) error {
		if ctx.NArg() < 1 {
			return cli.Exit(
				"You must provide a wishlist item ID. Execute `frens friend wishlist ls` to find out.",
				1,
			)
		}

		wID := strings.Join(ctx.Args().Slice(), " ")

		appCtx := jctx.FromCtx(ctx.Context)
		jr := appCtx.Repository.Journal()

		wOld, err := jr.GetFriendWishlistItem(wID)
		if err != nil {
			return err
		}

		inputForm := tui.NewEditorForm(tui.EditorOptions{
			Title:      "Edit wishlist item (" + wOld.ID + "):",
			SyntaxHint: lang.FormatWishlistItem,
		})

		inputForm.Textarea.SetValue(lang.RenderWishlistItem(wOld))

		teaUI := tea.NewProgram(inputForm, tea.WithMouseAllMotion())

		if _, err := teaUI.Run(); err != nil {
			log.Errorf("uh oh: %v", err)
			return err
		}

		infoTxt := inputForm.Textarea.Value()

		if infoTxt == "" {
			return errors.New("no wishlist item info provided")
		}

		wNew, err := lang.ExtractWishlistItem(infoTxt)

		desc := ctx.String("desc")
		link := ctx.String("link")
		price := ctx.String("price")
		tags := ctx.StringSlice("tag")

		if desc != "" {
			wNew.Desc = desc
		}

		if link != "" {
			wNew.Link = link
		}

		if price != "" {
			wNew.Price = price
		}

		if len(tags) > 0 {
			wNew.Tags = tags
		}

		if err != nil && !errors.Is(err, lang.ErrNoInfo) {
			log.Errorf(" failed to parse wishlist item info: %v", err)
			return err
		}

		if err := wNew.Validate(); err != nil {
			return err
		}

		err = appCtx.Repository.Update(func(j *journal.Journal) error {
			wNew, err = j.UpdateFriendWishlistItem(wOld, wNew)

			return err
		})
		if err != nil {
			return err
		}

		log.Info(" Wishlist item updated")
		log.Info("==> Wishlist Item Information\n")

		fmtr := formatter.WishlistItemTextFormatter{}

		o, _ := fmtr.FormatSingle(&wNew)
		fmt.Println(o)

		return nil
	},
}
