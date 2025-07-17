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
	"os/exec"
	"path/filepath"

	jctx "github.com/roma-glushko/frens/internal/context"
	"github.com/urfave/cli/v2"
)

var SyncCommand = &cli.Command{
	Name:    "sync",
	Aliases: []string{"s"},
	Usage:   "Synchronize your journal with a remote git repository",
	Action: func(ctx *cli.Context) error {
		jCtx := jctx.FromCtx(ctx.Context)
		jDir := jCtx.JournalDir
		gitDir := filepath.Join(jDir, ".git")

		if f, err := os.Stat(gitDir); err != nil || !f.IsDir() {
			return cli.Exit("No git repository found in the current journal directory. "+
				"Please initialize or connect to a remote repository first.", 1)
		}

		cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")

		cmd.Dir = jDir
		cmd.Stdin = os.Stdin

		if err := cmd.Run(); err == nil {
			fmt.Println("Pulling latest changes from the remote repository...")
			cmd := exec.Command(
				"git",
				"pull",
				"origin",
				"main",
				"--rebase",
				"--autostash",
				"--allow-unrelated-histories",
			)

			cmd.Dir = jDir
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			// Run and wait
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("git clone failed: %w", err)
			}
		}

		cmd = exec.Command("git", "status", "--porcelain")
		cmd.Dir = jDir
		output, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("git status failed: %w", err)
		}

		if len(output) == 0 {
			fmt.Println("No changes to commit.")
			return nil
		} else {
			fmt.Println("Committing changes to the remote repository...")

			cmd = exec.Command("git", "add", ".")
			cmd.Dir = jDir
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				return fmt.Errorf("git add failed: %w", err)
			}

			hostname, _ := os.Hostname()

			cmd = exec.Command("git", "commit", "-m", "🔄 sync: synchronize journal @ "+hostname)
			cmd.Dir = jDir
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				return fmt.Errorf("git commit failed: %w", err)
			}

			cmd := exec.Command("git", "branch", "-M", "main")

			cmd.Dir = jDir
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			// Run and wait
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("git branch failed: %w", err)
			}
		}

		fmt.Println("Pushing changes to the remote repository...")

		cmd = exec.Command("git", "push", "origin", "main")
		cmd.Dir = jDir
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("git push failed: %w", err)
		}

		fmt.Println("✅ Journal synchronized with the remote repository successfully.")

		return nil
	},
}
