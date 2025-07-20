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

package sync

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Git struct
type Git struct {
	RepoPath string
}

// NewGit creates a new Git instance
func NewGit(repoPath string) *Git {
	return &Git{
		RepoPath: repoPath,
	}
}

func (g Git) Installed() error {
	o, err := exec.LookPath("git")
	if err != nil {
		return fmt.Errorf("%s\n%s", o, err)
	}

	return nil
}

func (g Git) Inited() error {
	gitDir := filepath.Join(g.RepoPath, ".git")

	if f, err := os.Stat(gitDir); err != nil || !f.IsDir() {
		return fmt.Errorf(
			"no git repository found under %s. Please initialize or connect to a remote repository first",
			g.RepoPath,
		)
	}

	return nil
}

func (g Git) Init() error {
	cmd := exec.Command("git", "init")

	cmd.Dir = g.RepoPath
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
	}

	return nil
}

func (g Git) Clone(url string) error {
	cmd := exec.Command("git", "clone", url, g.RepoPath)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}

	return nil
}

func (g Git) GetBranchName() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")

	cmd.Dir = g.RepoPath
	cmd.Stdin = os.Stdin

	o, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s\n%s", o, err)
	}

	return strings.TrimSpace(string(o)), nil
}

func (g Git) Branch(branch string) error {
	cmd := exec.Command("git", "branch", "-M", branch)

	cmd.Dir = g.RepoPath
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run and wait
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git branch failed: %w", err)
	}

	return nil
}

func (g Git) GetStatus() (string, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = g.RepoPath

	status, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git status failed: %w", err)
	}

	return strings.TrimSpace(string(status)), nil
}

func (g Git) Commit(message string) error {
	cmd := exec.Command("git", "add", ".")
	cmd.Dir = g.RepoPath
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git add failed: %w", err)
	}

	cmd = exec.Command("git", "commit", "-m", message)
	cmd.Dir = g.RepoPath
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git commit failed: %w", err)
	}

	return nil
}

func (g Git) AddRemote(origin, url string) error {
	cmd := exec.Command("git", "remote", "add", origin, url)

	cmd.Dir = g.RepoPath
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run and wait
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git remote add failed: %w", err)
	}

	return nil
}

func (g Git) Push(origin, branch string) error {
	cmd := exec.Command("git", "push", origin, branch)
	cmd.Dir = g.RepoPath
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git push failed: %w", err)
	}

	return nil
}

func (g Git) Pull(origin, branch string) error {
	cmd := exec.Command(
		"git",
		"pull",
		origin,
		branch,
		"--rebase",
		"--autostash",
		"--allow-unrelated-histories",
	)

	cmd.Dir = g.RepoPath
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run and wait
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git pull failed: %w", err)
	}

	return nil
}
