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

func TestJournal_AddNoteWithFriend(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	a := []string{
		"frens",
		"-j",
		jDir,
		"friend",
		"add",
		"John Doe (aka John D.) :: A good friend #friends @NewYork $id:jdoe",
	}

	err = app.RunContext(ctx, a)
	require.NoError(t, err)

	a = []string{
		"frens",
		"-j",
		jDir,
		"note",
		"add",
		"John D. likes vanilla ice cream #friends @NewYork",
	}

	err = app.RunContext(ctx, a)
	require.NoError(t, err)
}

func TestNote_List(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add notes
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"note", "add",
		"Remember to buy groceries #personal",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"note", "add",
		"Meeting scheduled for next week #work",
	})
	require.NoError(t, err)

	// List all notes
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"note", "list",
	})
	require.NoError(t, err)
}

func TestNote_List_WithSearch(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add notes
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"note", "add",
		"Remember to buy groceries #personal",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"note", "add",
		"Meeting scheduled for next week #work",
	})
	require.NoError(t, err)

	// Search notes
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"note", "list",
		"--search", "groceries",
	})
	require.NoError(t, err)
}

func TestNote_List_WithTagFilter(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add notes with different tags
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"note", "add",
		"Personal reminder #personal",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"note", "add",
		"Work task reminder #work",
	})
	require.NoError(t, err)

	// Filter by tag
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"note", "list",
		"--tag", "work",
	})
	require.NoError(t, err)
}

func TestNote_List_WithSortAndReverse(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add notes
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"note", "add",
		"Alpha note #test",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"note", "add",
		"Beta note #test",
	})
	require.NoError(t, err)

	// List with alpha sort
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"note", "list",
		"--sort", "alpha",
	})
	require.NoError(t, err)

	// List with recency sort (default)
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"note", "list",
		"--sort", "recency",
	})
	require.NoError(t, err)

	// List with reverse sort
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"note", "list",
		"--reverse",
	})
	require.NoError(t, err)
}

func TestNote_List_WithDateFilters(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add notes
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"note", "add",
		"Today's note #test",
	})
	require.NoError(t, err)

	// List with from filter
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"note", "list",
		"--from", "yesterday",
	})
	require.NoError(t, err)

	// List with to filter
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"note", "list",
		"--to", "tomorrow",
	})
	require.NoError(t, err)

	// List with both from and to filters
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"note", "list",
		"--from", "yesterday",
		"--to", "tomorrow",
	})
	require.NoError(t, err)
}

func TestNote_AddWithDate(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add note with date prefix
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"note", "add",
		"yesterday :: This is a note from yesterday #reminder",
	})
	require.NoError(t, err)

	// List notes
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"note", "list",
	})
	require.NoError(t, err)
}
