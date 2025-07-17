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

var ConnectCommand = &cli.Command{
	Name:      "connect",
	Aliases:   []string{"con"},
	Usage:     "Connect an existing journal to a remote git repository",
	ArgsUsage: "<REPOSITORY>",
	Args:      true,
	Action: func(ctx *cli.Context) error {
		jCtx := jctx.FromCtx(ctx.Context)
		jDir := jCtx.JournalDir
		gitDir := filepath.Join(jDir, ".git")

		repoURL := ctx.Args().First()

		if f, err := os.Stat(gitDir); err != nil || !f.IsDir() {
			// if the .git directory does not exist, we should init git first
			cmd := exec.Command("git", "init")
			cmd.Dir = jDir
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to initialize git repository: %w", err)
			}
		}

		fmt.Println("Connecting to git repository", repoURL)

		cmd := exec.Command("git", "remote", "add", "origin", repoURL)

		cmd.Dir = jDir
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// Run and wait
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("git remote add failed: %w", err)
		}

		return nil
	},
}
