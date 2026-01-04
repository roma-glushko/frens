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
	"os"

	jctx "github.com/roma-glushko/frens/internal/context"
	"github.com/roma-glushko/frens/internal/sync"
	"github.com/urfave/cli/v2"
)

var SyncCommand = &cli.Command{
	Name:    "sync",
	Aliases: []string{"s"},
	Usage:   "Synchronize your journal with a remote git repository",
	Action: func(c *cli.Context) error {
		ctx := c.Context
		jCtx := jctx.FromCtx(ctx)
		jDir := jCtx.JournalDir

		git := sync.NewGit(jDir)

		if err := git.Installed(); err != nil {
			return fmt.Errorf("git is not installed or not found in PATH: %w", err)
		}

		if err := git.Inited(); err != nil {
			return err
		}

		origin := "origin"
		branch, err := git.GetBranchName(ctx)

		if err == nil {
			fmt.Println("Pulling latest changes from the remote repository...")

			if err := git.Pull(ctx, origin, branch); err != nil {
				return fmt.Errorf("git pull failed: %w", err)
			}
		} else {
			fmt.Println("failed to get current branch name, assuming no branch is set yet")
		}

		if branch == "" {
			branch = "main"
		}

		status, err := git.GetStatus(ctx)
		if err != nil {
			return err
		}

		if len(status) == 0 {
			fmt.Println("No changes to commit.")
			return nil
		} else {
			fmt.Println("Committing changes to the remote repository...")

			hostname, _ := os.Hostname()
			commit := "🔄 sync: synchronize journal @ " + hostname

			err := git.Commit(ctx, commit)
			if err != nil {
				return err
			}

			if err = git.Branch(ctx, branch); err != nil {
				return fmt.Errorf("git branch failed: %w", err)
			}
		}

		fmt.Println("Pushing changes to the remote repository...")

		if err = git.Push(ctx, origin, branch); err != nil {
			return fmt.Errorf("git push failed: %w", err)
		}

		fmt.Println("✅ Journal synchronized with the remote repository successfully.")

		return nil
	},
}
