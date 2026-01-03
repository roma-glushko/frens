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
	"testing"

	"github.com/roma-glushko/frens/cmd"
	"github.com/stretchr/testify/require"
)

func TestJournal_AddActivityWithFriend(t *testing.T) {
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
		"John Doe :: A good friend #friends @NewYork $id:john_doe",
	}

	err = app.RunContext(ctx, a)
	require.NoError(t, err)

	a = []string{
		"frens",
		"-j",
		jDir,
		"activity",
		"add",
		"Had a great time with John Doe at the park #friends @NewYork",
	}

	err = app.RunContext(ctx, a)
	require.NoError(t, err)
}

func TestActivity_List(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"activity", "add",
		"Had coffee with friends #social @CoffeeShop",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"activity", "add",
		"Team meeting at office #work @Office",
	})
	require.NoError(t, err)

	// List all activities
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"activity", "list",
	})
	require.NoError(t, err)
}

func TestActivity_List_WithSearch(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"activity", "add",
		"Had coffee with friends #social",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"activity", "add",
		"Team meeting at office #work",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"activity", "list",
		"--search", "coffee",
	})
	require.NoError(t, err)
}

func TestActivity_List_WithTagFilter(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"activity", "add",
		"Had coffee with friends #social",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"activity", "add",
		"Team meeting at office #work",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"activity", "list",
		"--tag", "work",
	})
	require.NoError(t, err)
}

func TestActivity_List_WithSortAndReverse(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add activities
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"activity", "add",
		"Alpha activity #test",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"activity", "add",
		"Beta activity #test",
	})
	require.NoError(t, err)

	// List with alpha sort
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"activity", "list",
		"--sort", "alpha",
	})
	require.NoError(t, err)

	// List with recency sort (default)
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"activity", "list",
		"--sort", "recency",
	})
	require.NoError(t, err)

	// List with reverse sort
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"activity", "list",
		"--reverse",
	})
	require.NoError(t, err)
}

func TestActivity_List_WithDateFilters(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add activities
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"activity", "add",
		"Today's activity #test",
	})
	require.NoError(t, err)

	// List with from filter
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"activity", "list",
		"--from", "yesterday",
	})
	require.NoError(t, err)

	// List with to filter
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"activity", "list",
		"--to", "tomorrow",
	})
	require.NoError(t, err)

	// List with both from and to filters
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"activity", "list",
		"--from", "yesterday",
		"--to", "tomorrow",
	})
	require.NoError(t, err)
}

func TestActivity_AddWithFriendReference(t *testing.T) {
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

	// Add activity with friend reference
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"activity", "add",
		"Had lunch with &john_doe #social @Restaurant",
	})
	require.NoError(t, err)

	// List activities to verify
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"activity", "list",
	})
	require.NoError(t, err)
}
