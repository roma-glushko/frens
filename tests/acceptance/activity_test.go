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

package acceptance

import (
	"github.com/roma-glushko/frens/cmd"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJournal_AddActivityWithFriend(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	a := []string{
		"frens",
		"-j",
		jDir,
		"friend",
		"add",
		"John Doe :: A good friend #friends @NewYork $id:john_doe",
	}

	err = app.RunContext(t.Context(), a)
	require.NoError(t, err)

	a = []string{
		"frens",
		"-j",
		jDir,
		"activity",
		"add",
		"Had a great time with John Doe at the park #friends @NewYork",
	}

	err = app.RunContext(t.Context(), a)
	require.NoError(t, err)
}
