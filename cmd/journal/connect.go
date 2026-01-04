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

package journal

import (
	"fmt"

	jctx "github.com/roma-glushko/frens/internal/context"
	"github.com/roma-glushko/frens/internal/sync"
	"github.com/urfave/cli/v2"
)

var ConnectCommand = &cli.Command{
	Name:      "connect",
	Aliases:   []string{"con"},
	Usage:     "Connect an existing journal to an empty remote git repository",
	ArgsUsage: "<REPOSITORY>",
	Args:      true,
	Action: func(c *cli.Context) error {
		ctx := c.Context
		jCtx := jctx.FromCtx(ctx)
		jDir := jCtx.JournalDir

		git := sync.NewGit(jDir)

		if err := git.Installed(); err != nil {
			return fmt.Errorf("git is not installed or not found in PATH: %w", err)
		}

		repoURL := c.Args().First()

		if err := git.Inited(); err != nil {
			// if the .git directory does not exist, we should init git first

			if err := git.Init(ctx); err != nil {
				return err
			}
		}

		fmt.Println("Connecting to git repository", repoURL)

		return git.AddRemote(ctx, "origin", repoURL)
	},
}
