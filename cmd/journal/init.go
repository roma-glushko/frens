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
	"github.com/roma-glushko/frens/internal/log"
	"github.com/roma-glushko/frens/internal/tui"
	"github.com/urfave/cli/v2"
)

var InitCommand = &cli.Command{
	Name:      "init",
	Aliases:   []string{"i"},
	Usage:     "Init a new journal",
	UsageText: "frens journal init",
	Description: `Initialize a new journal in the default location (~/.config/frens/).

Use the global --journal flag to initialize in a custom location.

Examples:
  frens journal init                         # init in default location
  frens --journal ~/my-frens journal init    # init in custom location
`,
	Flags: []cli.Flag{},
	Action: func(c *cli.Context) error {
		ctx := c.Context
		jCtx := jctx.FromCtx(ctx)

		jDir := jCtx.JournalDir
		s := jCtx.Store

		if s.Exist(ctx) {
			// TODO: check if interactive mode is enabled
			log.Infof("A journal already exists at %s\n", jDir)
			if tui.ConfirmAction(log.WarnPrompt("Do you want to overwrite the existing journal?")) {
				log.Progress("Overwriting the existing journal...")
			} else {
				log.Canceled("Journal initialization canceled.")
				return nil
			}
		}

		if err := s.Init(ctx); err != nil {
			return fmt.Errorf("failed to initialize the journal at %s: %w", jDir, err)
		}

		log.Successf("Journal initialized at %s", jDir)

		return nil
	},
}
