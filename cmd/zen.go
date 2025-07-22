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

	"github.com/urfave/cli/v2"
)

var ZenCommand = &cli.Command{
	Name:  "zen",
	Usage: "Print the zen of friendship",
	Action: func(_ *cli.Context) error {
		fmt.Println("🌿 The Zen of Friendship:")
		fmt.Println("   • Treat others as you would like them to treat you.")
		fmt.Println("   • You should \"buy\" yourself a friend.")
		return nil
	},
}
