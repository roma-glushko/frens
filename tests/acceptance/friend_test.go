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

func TestFriend_List(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add a friend first
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "add",
		"John Doe :: A good friend #friends @NewYork $id:john_doe",
	})
	require.NoError(t, err)

	// Add another friend
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "add",
		"Jane Smith :: Work colleague #work @SanFrancisco $id:jane_smith",
	})
	require.NoError(t, err)

	// List all friends
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "list",
	})
	require.NoError(t, err)
}

func TestFriend_List_WithSearch(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add friends
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "add",
		"John Doe :: A good friend #friends @NewYork $id:john_doe",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "add",
		"Jane Smith :: Work colleague #work @SanFrancisco $id:jane_smith",
	})
	require.NoError(t, err)

	// Search for specific friend
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "list",
		"--search", "John",
	})
	require.NoError(t, err)
}

func TestFriend_List_WithTagFilter(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add friends with different tags
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "add",
		"John Doe :: A good friend #friends @NewYork $id:john_doe",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "add",
		"Jane Smith :: Work colleague #work @SanFrancisco $id:jane_smith",
	})
	require.NoError(t, err)

	// Filter by tag
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "list",
		"--tag", "work",
	})
	require.NoError(t, err)
}

func TestFriend_List_WithLocationFilter(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add friends with different locations
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "add",
		"John Doe :: A good friend #friends @NewYork $id:john_doe",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "add",
		"Jane Smith :: Work colleague #work @SanFrancisco $id:jane_smith",
	})
	require.NoError(t, err)

	// Filter by location
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "list",
		"--location", "NewYork",
	})
	require.NoError(t, err)
}

func TestFriend_List_WithSortAndReverse(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add friends
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "add",
		"John Doe :: A good friend #friends @NewYork $id:john_doe",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "add",
		"Alice Brown :: College friend #college @Boston $id:alice_brown",
	})
	require.NoError(t, err)

	// List with alpha sort
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "list",
		"--sort", "alpha",
	})
	require.NoError(t, err)

	// List with reverse sort
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "list",
		"--sort", "alpha",
		"--reverse",
	})
	require.NoError(t, err)
}

func TestFriend_Delete(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add a friend
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "add",
		"John Doe :: A good friend #friends @NewYork $id:john_doe",
	})
	require.NoError(t, err)

	// Delete the friend with force flag
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "delete",
		"--force",
		"john_doe",
	})
	require.NoError(t, err)
}

func TestFriend_Delete_Multiple(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add multiple friends
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "add",
		"John Doe :: A good friend #friends @NewYork $id:john_doe",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "add",
		"Jane Smith :: Work colleague #work @SanFrancisco $id:jane_smith",
	})
	require.NoError(t, err)

	// Delete multiple friends at once
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "delete",
		"--force",
		"john_doe",
		"jane_smith",
	})
	require.NoError(t, err)
}
