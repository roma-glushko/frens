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

func TestLocation_List(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add locations
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"location", "add",
		"New York City :: The Big Apple #city @USA $id:nyc",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"location", "add",
		"San Francisco :: Tech hub #city @USA $id:sf",
	})
	require.NoError(t, err)

	// List all locations
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"location", "list",
	})
	require.NoError(t, err)
}

func TestLocation_List_WithSearch(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add locations
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"location", "add",
		"New York City :: The Big Apple #city @USA $id:nyc",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"location", "add",
		"San Francisco :: Tech hub #city @USA $id:sf",
	})
	require.NoError(t, err)

	// Search locations
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"location", "list",
		"--search", "York",
	})
	require.NoError(t, err)
}

func TestLocation_List_WithTagFilter(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add locations with different tags
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"location", "add",
		"New York City :: The Big Apple #city #metro @USA $id:nyc",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"location", "add",
		"Central Park :: Green space #park #outdoor @USA $id:central_park",
	})
	require.NoError(t, err)

	// Filter by tag
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"location", "list",
		"--tag", "park",
	})
	require.NoError(t, err)
}

func TestLocation_List_WithCountryFilter(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add locations with different countries
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"location", "add",
		"New York City :: The Big Apple #city @USA $id:nyc",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"location", "add",
		"London :: Capital of UK #city @UK $id:london",
	})
	require.NoError(t, err)

	// Filter by country
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"location", "list",
		"--country", "USA",
	})
	require.NoError(t, err)
}

func TestLocation_List_WithSortAndReverse(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add locations
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"location", "add",
		"Boston :: Historic city #city @USA $id:boston",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"location", "add",
		"Austin :: Texas capital #city @USA $id:austin",
	})
	require.NoError(t, err)

	// List with alpha sort
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"location", "list",
		"--sort", "alpha",
	})
	require.NoError(t, err)

	// List with reverse sort
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"location", "list",
		"--sort", "alpha",
		"--reverse",
	})
	require.NoError(t, err)
}

func TestLocation_Delete(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add a location
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"location", "add",
		"New York City :: The Big Apple #city @USA $id:nyc",
	})
	require.NoError(t, err)

	// Delete the location with force flag
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"location", "delete",
		"--force",
		"nyc",
	})
	require.NoError(t, err)
}

func TestLocation_Delete_Multiple(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add multiple locations
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"location", "add",
		"New York City :: The Big Apple #city @USA $id:nyc",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"location", "add",
		"San Francisco :: Tech hub #city @USA $id:sf",
	})
	require.NoError(t, err)

	// Delete multiple locations at once
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"location", "delete",
		"--force",
		"nyc",
		"sf",
	})
	require.NoError(t, err)
}
