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

package acceptance

import (
	"testing"

	"github.com/roma-glushko/frens/cmd"
	"github.com/stretchr/testify/require"
)

func TestFriendDate_Add(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add a friend first
	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "add",
		"John Doe :: A good friend #friends @NewYork $id:john_doe",
	})
	require.NoError(t, err)

	// Add a date
	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "date", "add",
		"john_doe",
		"birthday :: May 13th",
	})
	require.NoError(t, err)
}

func TestFriendDate_Add_WithFlags(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add a friend first
	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "add",
		"John Doe :: A good friend #friends @NewYork $id:john_doe",
	})
	require.NoError(t, err)

	// Add a date using CLI flags
	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "date", "add",
		"john_doe",
		"--desc", "Birthday",
		"--date", "1990-05-13",
		"--tag", "birthday",
	})
	require.NoError(t, err)
}

func TestFriendDate_Add_Anniversary(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add a friend first
	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "add",
		"John Doe :: A good friend #friends @NewYork $id:john_doe",
	})
	require.NoError(t, err)

	// Add an anniversary date
	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "date", "add",
		"john_doe",
		"anniversary :: 2009-9-09 #important",
	})
	require.NoError(t, err)
}

func TestFriendDate_List(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add a friend with dates
	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "add",
		"John Doe :: A good friend #friends @NewYork $id:john_doe",
	})
	require.NoError(t, err)

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "date", "add",
		"john_doe",
		"birthday :: May 13th",
	})
	require.NoError(t, err)

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "date", "add",
		"john_doe",
		"anniversary :: 2009-9-09",
	})
	require.NoError(t, err)

	// List all dates
	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "date", "list",
	})
	require.NoError(t, err)
}

func TestFriendDate_List_WithFriendFilter(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add friends with dates
	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "add",
		"John Doe :: A good friend #friends @NewYork $id:john_doe",
	})
	require.NoError(t, err)

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "add",
		"Jane Smith :: Work colleague #work @SanFrancisco $id:jane_smith",
	})
	require.NoError(t, err)

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "date", "add",
		"john_doe",
		"birthday :: May 13th",
	})
	require.NoError(t, err)

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "date", "add",
		"jane_smith",
		"birthday :: June 20th",
	})
	require.NoError(t, err)

	// List dates for specific friend
	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "date", "list",
		"--with", "john_doe",
	})
	require.NoError(t, err)
}

func TestFriendDate_List_WithTagFilter(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add friend with tagged dates
	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "add",
		"John Doe :: A good friend #friends @NewYork $id:john_doe",
	})
	require.NoError(t, err)

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "date", "add",
		"john_doe",
		"birthday :: May 13th #birthday",
	})
	require.NoError(t, err)

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "date", "add",
		"john_doe",
		"anniversary :: 2009-9-09 #anniversary",
	})
	require.NoError(t, err)

	// Filter by tag
	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "date", "list",
		"--tag", "birthday",
	})
	require.NoError(t, err)
}

func TestFriendDate_List_WithSearch(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add friend with dates
	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "add",
		"John Doe :: A good friend #friends @NewYork $id:john_doe",
	})
	require.NoError(t, err)

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "date", "add",
		"john_doe",
		"birthday :: May 13th",
	})
	require.NoError(t, err)

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "date", "add",
		"john_doe",
		"wedding anniversary :: 2009-9-09",
	})
	require.NoError(t, err)

	// Search dates
	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "date", "list",
		"--search", "wedding",
	})
	require.NoError(t, err)
}
