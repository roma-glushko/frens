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

var EditCommand = &cli.Command{
	Name:    "edit",
	Aliases: []string{"e"},
	Usage:   "Edit journal raw files",
	Flags:   []cli.Flag{},
	Action: func(ctx *cli.Context) error {
		jctx := jctx.FromCtx(ctx.Context)

		editor := GetEditor()

		cmd := exec.Command(
			editor,
			filepath.Join(jctx.JournalDir, "friends.toml"),
		) // TODO: make it configurable
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("error running editor: %s", err)
		}

		return nil
	},
}

func GetEditor() string {
	editor := os.Getenv("EDITOR")

	if editor == "" {
		editor = "vim"
	}

	return editor
}
