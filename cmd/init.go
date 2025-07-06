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

package cmd

import (
	"fmt"

	"github.com/roma-glushko/frens/internal/journaldir"
	"github.com/urfave/cli/v2"
)

var InitCommand = &cli.Command{
	Name:    "init",
	Aliases: []string{"i"},
	Usage:   "Init a new journal",
	Flags:   []cli.Flag{},
	Action: func(_ *cli.Context) error {
		journalDir, err := journaldir.DefaultDir()
		if err != nil {
			return err
		}

		err = journaldir.Init(journalDir)
		if err != nil {
			return err
		}

		fmt.Println("A new journal's initialized at", journalDir)

		return nil
	},
}
