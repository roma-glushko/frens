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
	"github.com/urfave/cli/v2"
)

func addFriend(t *testing.T, app *cli.App, jDir, info string) {
	t.Helper()

	err := app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "add",
		info,
	})
	require.NoError(t, err)
}

func addDate(t *testing.T, app *cli.App, jDir, friendID, info string) {
	t.Helper()

	err := app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "date", "add",
		friendID,
		info,
	})
	require.NoError(t, err)
}

// --- Add ---

func TestFriendDate_Add(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	addFriend(t, &app, jDir, "John Doe :: A good friend #friends @NewYork $id:john_doe")
	addDate(t, &app, jDir, "john_doe", "May 13th :: birthday")

	j := LoadJournal(t, jDir)
	require.Len(t, j.Friends, 1)
	require.Len(t, j.Friends[0].Dates, 1)

	dt := j.Friends[0].Dates[0]
	require.NotEmpty(t, dt.ID)
	require.Equal(t, "May 13th", dt.DateExpr)
	require.Equal(t, "birthday", dt.Desc)
	require.Equal(t, "gregorian", dt.Calendar)
}

func TestFriendDate_Add_WithFlags(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	addFriend(t, &app, jDir, "John Doe :: A good friend $id:john_doe")

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

	addFriend(t, &app, jDir, "John Doe :: A good friend #friends @NewYork $id:john_doe")
	addDate(t, &app, jDir, "john_doe", "2009-9-09 :: anniversary #important")

	j := LoadJournal(t, jDir)
	require.Len(t, j.Friends[0].Dates, 1)

	dt := j.Friends[0].Dates[0]
	require.Equal(t, "2009-9-09", dt.DateExpr)
	require.Equal(t, "anniversary", dt.Desc)
	require.Equal(t, []string{"important"}, dt.Tags)
}

func TestFriendDate_Add_WithCalendar(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	addFriend(t, &app, jDir, "John Doe :: A good friend $id:john_doe")
	addDate(t, &app, jDir, "john_doe", "Av 16 5784 :: birthday $cal:hebrew")

	j := LoadJournal(t, jDir)
	dt := j.Friends[0].Dates[0]
	require.Equal(t, "hebrew", dt.Calendar)
	require.Equal(t, "Av 16 5784", dt.DateExpr)
	require.Equal(t, "birthday", dt.Desc)
}

func TestFriendDate_Add_MultipleDates(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	addFriend(t, &app, jDir, "John Doe :: A good friend $id:john_doe")
	addDate(t, &app, jDir, "john_doe", "May 13th :: birthday #birthday")
	addDate(t, &app, jDir, "john_doe", "2009-9-09 :: anniversary #anniversary")
	addDate(t, &app, jDir, "john_doe", "2015-06-15 :: graduation")

	j := LoadJournal(t, jDir)
	require.Len(t, j.Friends[0].Dates, 3)
}

func TestFriendDate_Add_NonexistentFriend(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "date", "add",
		"nobody",
		"May 13th :: birthday",
	})
	require.Error(t, err)
}

// --- List ---

func TestFriendDate_List(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	addFriend(t, &app, jDir, "John Doe :: A good friend $id:john_doe")
	addDate(t, &app, jDir, "john_doe", "May 13th :: birthday")
	addDate(t, &app, jDir, "john_doe", "2009-9-09 :: anniversary")

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "date", "list",
	})
	require.NoError(t, err)
}

func TestFriendDate_List_Empty(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	addFriend(t, &app, jDir, "John Doe :: A good friend $id:john_doe")

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

	addFriend(t, &app, jDir, "John Doe :: A good friend $id:john_doe")
	addFriend(t, &app, jDir, "Jane Smith :: Work colleague $id:jane_smith")
	addDate(t, &app, jDir, "john_doe", "May 13th :: birthday")
	addDate(t, &app, jDir, "jane_smith", "June 20th :: birthday")

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

	addFriend(t, &app, jDir, "John Doe :: A good friend $id:john_doe")
	addDate(t, &app, jDir, "john_doe", "May 13th :: birthday #birthday")
	addDate(t, &app, jDir, "john_doe", "2009-9-09 :: anniversary #anniversary")

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

	addFriend(t, &app, jDir, "John Doe :: A good friend $id:john_doe")
	addDate(t, &app, jDir, "john_doe", "May 13th :: birthday")
	addDate(t, &app, jDir, "john_doe", "2009-9-09 :: wedding anniversary")

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "date", "list",
		"--search", "wedding",
	})
	require.NoError(t, err)
}

func TestFriendDate_List_MultipleFriends(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	addFriend(t, &app, jDir, "John Doe :: A good friend $id:john_doe")
	addFriend(t, &app, jDir, "Jane Smith :: Work colleague $id:jane_smith")
	addDate(t, &app, jDir, "john_doe", "May 13th :: birthday #birthday")
	addDate(t, &app, jDir, "jane_smith", "June 20th :: birthday #birthday")
	addDate(t, &app, jDir, "jane_smith", "2015-03-14 :: anniversary #anniversary")

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "date", "list",
	})
	require.NoError(t, err)
}

func TestFriendDate_List_CombinedFilters(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	addFriend(t, &app, jDir, "John Doe :: A good friend $id:john_doe")
	addFriend(t, &app, jDir, "Jane Smith :: Work colleague $id:jane_smith")
	addDate(t, &app, jDir, "john_doe", "May 13th :: birthday #birthday")
	addDate(t, &app, jDir, "john_doe", "2009-9-09 :: anniversary #anniversary")
	addDate(t, &app, jDir, "jane_smith", "June 20th :: birthday #birthday")

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "date", "list",
		"--with", "john_doe",
		"--tag", "birthday",
	})
	require.NoError(t, err)
}

// --- List with output formats ---

func TestFriendDate_List_JSONFormat(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	addFriend(t, &app, jDir, "John Doe :: A good friend $id:john_doe")
	addDate(t, &app, jDir, "john_doe", "May 13th :: birthday #birthday")

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"--format", "json",
		"friend", "date", "list",
	})
	require.NoError(t, err)
}

func TestFriendDate_List_MarkdownFormat(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	addFriend(t, &app, jDir, "John Doe :: A good friend $id:john_doe")
	addDate(t, &app, jDir, "john_doe", "May 13th :: birthday #birthday")

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"--format", "markdown",
		"friend", "date", "list",
	})
	require.NoError(t, err)
}

func TestFriendDate_List_Compact(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	addFriend(t, &app, jDir, "John Doe :: A good friend $id:john_doe")
	addDate(t, &app, jDir, "john_doe", "May 13th :: birthday")
	addDate(t, &app, jDir, "john_doe", "2009-9-09 :: anniversary")

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"--compact",
		"friend", "date", "list",
	})
	require.NoError(t, err)
}

// --- Delete ---

func TestFriendDate_Delete(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	addFriend(t, &app, jDir, "John Doe :: A good friend $id:john_doe")
	addDate(t, &app, jDir, "john_doe", "May 13th :: birthday")

	j := LoadJournal(t, jDir)
	require.Len(t, j.Friends[0].Dates, 1)

	dateID := j.Friends[0].Dates[0].ID

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "date", "delete",
		"--force",
		dateID,
	})
	require.NoError(t, err)

	j = LoadJournal(t, jDir)
	require.Empty(t, j.Friends[0].Dates)
}

func TestFriendDate_Delete_Multiple(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	addFriend(t, &app, jDir, "John Doe :: A good friend $id:john_doe")
	addDate(t, &app, jDir, "john_doe", "May 13th :: birthday")
	addDate(t, &app, jDir, "john_doe", "2009-9-09 :: anniversary")

	j := LoadJournal(t, jDir)
	require.Len(t, j.Friends[0].Dates, 2)

	id1 := j.Friends[0].Dates[0].ID
	id2 := j.Friends[0].Dates[1].ID

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "date", "delete",
		"--force",
		id1,
		id2,
	})
	require.NoError(t, err)

	j = LoadJournal(t, jDir)
	require.Empty(t, j.Friends[0].Dates)
}

func TestFriendDate_Delete_NotFound(t *testing.T) {
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	addFriend(t, &app, jDir, "John Doe :: A good friend $id:john_doe")

	err = app.RunContext(t.Context(), []string{
		"frens", "-j", jDir,
		"friend", "date", "delete",
		"--force",
		"nonexistent_id",
	})
	require.Error(t, err)
}
