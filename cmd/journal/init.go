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

	"github.com/roma-glushko/frens/internal/tui"

	"github.com/roma-glushko/frens/internal/journaldir"
	"github.com/urfave/cli/v2"
)

var InitCommand = &cli.Command{
	Name:    "init",
	Aliases: []string{"i"},
	Usage:   "Init a new journal",
	Flags:   []cli.Flag{},
	Action: func(_ *cli.Context) error {
		jDir, err := journaldir.DefaultDir()
		if err != nil {
			return err
		}

		if journaldir.Exists(jDir) {
			// TODO: check if interactive mode is enabled
			fmt.Println("A journal already exists at", jDir)
			if tui.ConfirmAction("\n⚠️  Do you want to overwrite the existing journal under?") {
				fmt.Println("Overwriting the existing journal...")
			} else {
				fmt.Println("\n↩️  Journal initialization cancelled.")
				return nil
			}
		}

		err = journaldir.Init(jDir)
		if err != nil {
			return fmt.Errorf("failed to initialize the journal at %s: %w", jDir, err)
		}

		fmt.Println("✅ A new journal's initialized at", jDir)

		return nil
	},
}
