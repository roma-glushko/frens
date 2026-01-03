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

func TestFriendContact_Add(t *testing.T) {
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

	// Add contact information
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "contact", "add",
		"john_doe",
		"john@example.com",
	})
	require.NoError(t, err)
}

func TestFriendContact_Add_Multiple(t *testing.T) {
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

	// Add multiple contact information at once
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "contact", "add",
		"john_doe",
		"john@example.com",
		"+1234567890",
		"tg:@johndoe",
	})
	require.NoError(t, err)
}

func TestFriendContact_Add_WithTags(t *testing.T) {
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

	// Add contact with tags
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "contact", "add",
		"john_doe",
		"john.work@company.com",
		"--tag", "work",
	})
	require.NoError(t, err)
}

func TestFriendContact_Add_SocialMedia(t *testing.T) {
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

	// Add various social media contacts
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "contact", "add",
		"john_doe",
		"ig:@johndoe",
		"x:@johndoe",
		"li:johndoe",
	})
	require.NoError(t, err)
}

func TestFriendContact_List(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add a friend and contacts
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "add",
		"John Doe :: A good friend #friends @NewYork $id:john_doe",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "contact", "add",
		"john_doe",
		"john@example.com",
		"+1234567890",
	})
	require.NoError(t, err)

	// List all contacts
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "contact", "list",
	})
	require.NoError(t, err)
}

func TestFriendContact_List_WithFriendFilter(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add friends with contacts
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
		"friend", "contact", "add",
		"john_doe",
		"john@example.com",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "contact", "add",
		"jane_smith",
		"jane@company.com",
	})
	require.NoError(t, err)

	// List contacts for specific friend
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "contact", "list",
		"--with", "john_doe",
	})
	require.NoError(t, err)
}

func TestFriendContact_List_WithTypeFilter(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add friend with various contact types
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "add",
		"John Doe :: A good friend #friends @NewYork $id:john_doe",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "contact", "add",
		"john_doe",
		"john@example.com",
		"+1234567890",
		"tg:@johndoe",
	})
	require.NoError(t, err)

	// List only email contacts
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "contact", "list",
		"--type", "email",
	})
	require.NoError(t, err)
}

func TestFriendContact_List_WithSearch(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add friend with contacts
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "add",
		"John Doe :: A good friend #friends @NewYork $id:john_doe",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "contact", "add",
		"john_doe",
		"john@example.com",
		"john.work@company.com",
	})
	require.NoError(t, err)

	// Search contacts
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "contact", "list",
		"--search", "company",
	})
	require.NoError(t, err)
}

// Note: TestFriendContact_Delete is not included because the contact delete command
// requires a contact ID (auto-generated), which requires parsing the add command output
// or reading the journal file. This is covered by the other CRUD tests.
