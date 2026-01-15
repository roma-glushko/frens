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
	"github.com/roma-glushko/frens/internal/log"
	"github.com/roma-glushko/frens/internal/sync"
	"github.com/roma-glushko/frens/internal/tui"
	"github.com/urfave/cli/v2"
)

var SyncCommand = &cli.Command{
	Name:      "sync",
	Aliases:   []string{"s"},
	Usage:     "Synchronize your journal with a remote git repository",
	UsageText: "frens journal sync",
	Description: `Pull latest changes, commit local changes, and push to remote.

Requires git to be installed and a remote repository configured.
Use 'frens journal connect' to set up the remote repository first.

Examples:
  frens journal sync                         # sync with remote
  frens --journal ~/my-frens journal sync    # sync a specific journal
`,
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
			if err := tui.RunWithProgress("Pulling latest changes from remote...", func() error {
				return git.Pull(ctx, origin, branch)
			}); err != nil {
				return fmt.Errorf("git pull failed: %w", err)
			}
		} else {
			log.Warn("Failed to get current branch name, assuming no branch is set yet\n")
		}

		if branch == "" {
			branch = "main"
		}

		status, err := git.GetStatus(ctx)
		if err != nil {
			return err
		}

		if len(status) == 0 {
			log.Info("No changes to commit.\n")
			return nil
		}

		hostname, _ := os.Hostname()
		commit := "🔄 sync: synchronize journal @ " + hostname

		if err := tui.RunWithProgress("Committing changes...", func() error {
			if err := git.Commit(ctx, commit); err != nil {
				return err
			}
			return git.Branch(ctx, branch)
		}); err != nil {
			return err
		}

		if err := tui.RunWithProgress("Pushing changes to remote...", func() error {
			return git.Push(ctx, origin, branch)
		}); err != nil {
			return fmt.Errorf("git push failed: %w", err)
		}

		log.Success("Journal synchronized successfully")

		return nil
	},
}
