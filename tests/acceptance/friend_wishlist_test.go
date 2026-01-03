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

func TestFriendWishlist_Add(t *testing.T) {
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

	// Add a wishlist item with description
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "wishlist", "add",
		"john_doe",
		"Cool mechanical keyboard #techgift",
	})
	require.NoError(t, err)
}

func TestFriendWishlist_Add_WithFlags(t *testing.T) {
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

	// Add a wishlist item using CLI flags
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "wishlist", "add",
		"john_doe",
		"--desc", "Cool mechanical keyboard",
		"--price", "150USD",
		"--tag", "tech",
		"--tag", "gift",
	})
	require.NoError(t, err)
}

func TestFriendWishlist_Add_MultipleItems(t *testing.T) {
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

	// Add multiple wishlist items
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "wishlist", "add",
		"john_doe",
		"Cool mechanical keyboard #tech",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "wishlist", "add",
		"john_doe",
		"Gaming mouse #gaming",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "wishlist", "add",
		"john_doe",
		"Noise cancelling headphones #audio",
	})
	require.NoError(t, err)
}

func TestFriendWishlist_List(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add a friend with wishlist items
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "add",
		"John Doe :: A good friend #friends @NewYork $id:john_doe",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "wishlist", "add",
		"john_doe",
		"Cool keyboard #tech",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "wishlist", "add",
		"john_doe",
		"Gaming mouse #gaming",
	})
	require.NoError(t, err)

	// List all wishlist items
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "wishlist", "list",
	})
	require.NoError(t, err)
}

func TestFriendWishlist_List_WithFriendFilter(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add multiple friends with wishlist items
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

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "wishlist", "add",
		"john_doe",
		"Gaming keyboard #tech",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "wishlist", "add",
		"jane_smith",
		"Book on leadership #books",
	})
	require.NoError(t, err)

	// List wishlist items for specific friend
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "wishlist", "list",
		"--with", "john_doe",
	})
	require.NoError(t, err)
}

func TestFriendWishlist_List_WithTagFilter(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add friend with tagged wishlist items
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "add",
		"John Doe :: A good friend #friends @NewYork $id:john_doe",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "wishlist", "add",
		"john_doe",
		"Mechanical keyboard #tech #office",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "wishlist", "add",
		"john_doe",
		"Gaming headset #gaming #audio",
	})
	require.NoError(t, err)

	// Filter by tag
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "wishlist", "list",
		"--tag", "gaming",
	})
	require.NoError(t, err)
}

func TestFriendWishlist_List_WithSearch(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add friend with wishlist items
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "add",
		"John Doe :: A good friend #friends @NewYork $id:john_doe",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "wishlist", "add",
		"john_doe",
		"Mechanical keyboard",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "wishlist", "add",
		"john_doe",
		"Wireless mouse",
	})
	require.NoError(t, err)

	// Search wishlist items
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "wishlist", "list",
		"--search", "keyboard",
	})
	require.NoError(t, err)
}
