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

package journal

import (
	"fmt"
	"os"

	jctx "github.com/roma-glushko/frens/internal/context"
	"github.com/roma-glushko/frens/internal/sync"
	"github.com/roma-glushko/frens/internal/tui"
	"github.com/urfave/cli/v2"
)

var CloneCommand = &cli.Command{
	Name:      "clone",
	Aliases:   []string{"cl"},
	Usage:     "Clone a journal from a remote git repository",
	ArgsUsage: "<REPOSITORY>",
	Args:      true,
	Action: func(c *cli.Context) error {
		ctx := c.Context
		jCtx := jctx.FromCtx(ctx)
		jDir := jCtx.JournalDir

		repoURL := c.Args().First()

		git := sync.NewGit(jDir)

		if err := git.Installed(); err != nil {
			return err
		}

		if err := git.Inited(); err == nil {
			// TODO: check if interactive mode is enabled
			if tui.ConfirmAction("\n⚠️  Do you want to overwrite the existing journal under?") {
				fmt.Println("Overwriting the existing journal...")

				if err = os.RemoveAll(jDir); err != nil {
					return cli.Exit(
						"Failed to remove existing journal .git directory: "+err.Error(),
						1,
					)
				}
			} else {
				fmt.Println("\n↩️  Journal initialization canceled.")
				return nil
			}
		}

		if err := git.Clone(ctx, repoURL); err != nil {
			return err
		}

		fmt.Println("✅ A new journal's initialized at", jDir)

		return nil
	},
}
