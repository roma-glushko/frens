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
	"os"
	"path/filepath"
	"testing"

	"github.com/roma-glushko/frens/cmd"
	"github.com/stretchr/testify/require"
)

func TestJournal_Import_FromFile(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Create a FrenTXT file to import
	importContent := `/f
Alice (aka Al, Ally) :: Close friend from college.
#college, #bestie @Boston $id:alice

/f
Bob :: Work colleague
#work @NYC $id:bob

/l
Paris, France (aka City of Light) :: A city I visited.
#travel $id:paris

/n
2023-08-15 :: Had dinner with Alice.
#catchup @paris

/act
2023-08-16 :: Jogged with Bob.
#exercise @NYC
`
	importFile := filepath.Join(t.TempDir(), "import.frentxt")
	err = os.WriteFile(importFile, []byte(importContent), 0o644)
	require.NoError(t, err)

	// Import the file
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"journal", "import",
		importFile,
	})
	require.NoError(t, err)

	// Verify friends were imported by listing them
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "list",
	})
	require.NoError(t, err)

	// Verify locations were imported
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"location", "list",
	})
	require.NoError(t, err)

	// Verify notes were imported
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"note", "list",
	})
	require.NoError(t, err)

	// Verify activities were imported
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"activity", "list",
	})
	require.NoError(t, err)
}

func TestJournal_Import_DryRun(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Create a FrenTXT file
	importContent := `/f
Charlie :: Test friend
#test $id:charlie
`
	importFile := filepath.Join(t.TempDir(), "import.frentxt")
	err = os.WriteFile(importFile, []byte(importContent), 0o644)
	require.NoError(t, err)

	// Import with dry-run
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"journal", "import",
		"--dry-run",
		importFile,
	})
	require.NoError(t, err)

	// The friend should NOT be in the journal (dry-run doesn't save)
	// We can verify by checking the friends.toml file doesn't contain "Charlie"
	friendsFile := filepath.Join(jDir, "friends.toml")
	content, err := os.ReadFile(friendsFile)
	require.NoError(t, err)
	require.NotContains(t, string(content), "Charlie")
}

func TestJournal_Import_WithAlias(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	importContent := `/f
Dave :: Another friend
#friend $id:dave
`
	importFile := filepath.Join(t.TempDir(), "import.frentxt")
	err = os.WriteFile(importFile, []byte(importContent), 0o644)
	require.NoError(t, err)

	// Use the "imp" alias
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"journal", "imp",
		importFile,
	})
	require.NoError(t, err)
}

func TestJournal_Export_ToFile(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add some data first
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "add",
		"Eve (aka Evie) :: Security expert #infosec @Berlin $id:eve",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"location", "add",
		"Berlin, Germany :: Capital city #city $id:berlin",
	})
	require.NoError(t, err)

	// Export to file
	exportFile := filepath.Join(t.TempDir(), "export.frentxt")
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"journal", "export",
		exportFile,
	})
	require.NoError(t, err)

	// Verify the file was created and contains expected content
	require.FileExists(t, exportFile)

	content, err := os.ReadFile(exportFile)
	require.NoError(t, err)

	require.Contains(t, string(content), "/f")
	require.Contains(t, string(content), "Eve")
	require.Contains(t, string(content), "Evie")
	require.Contains(t, string(content), "#infosec")

	require.Contains(t, string(content), "/l")
	require.Contains(t, string(content), "Berlin")
	require.Contains(t, string(content), "Germany")
}

func TestJournal_Export_FilterFriends(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add friend and location
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "add",
		"Frank :: Test friend #test $id:frank",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"location", "add",
		"London, UK :: Capital #city $id:london",
	})
	require.NoError(t, err)

	// Export only friends
	exportFile := filepath.Join(t.TempDir(), "friends_only.frentxt")
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"journal", "export",
		"--friends",
		exportFile,
	})
	require.NoError(t, err)

	content, err := os.ReadFile(exportFile)
	require.NoError(t, err)

	require.Contains(t, string(content), "/f")
	require.Contains(t, string(content), "Frank")
	require.NotContains(t, string(content), "/l")
	require.NotContains(t, string(content), "London")
}

func TestJournal_Export_FilterLocations(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add friend and location
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "add",
		"Grace :: Test friend #test $id:grace",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"location", "add",
		"Tokyo, Japan :: Capital #city $id:tokyo",
	})
	require.NoError(t, err)

	// Export only locations
	exportFile := filepath.Join(t.TempDir(), "locations_only.frentxt")
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"journal", "export",
		"--locations",
		exportFile,
	})
	require.NoError(t, err)

	content, err := os.ReadFile(exportFile)
	require.NoError(t, err)

	require.Contains(t, string(content), "/l")
	require.Contains(t, string(content), "Tokyo")
	require.NotContains(t, string(content), "/f")
	require.NotContains(t, string(content), "Grace")
}

func TestJournal_Export_WithAlias(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "add",
		"Henry :: Test #test $id:henry",
	})
	require.NoError(t, err)

	// Use the "exp" alias
	exportFile := filepath.Join(t.TempDir(), "alias_export.frentxt")
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"journal", "exp",
		exportFile,
	})
	require.NoError(t, err)

	require.FileExists(t, exportFile)
}

func TestJournal_ImportExport_RoundTrip(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	// Initialize first journal and add data
	jDir1, err := InitJournal(t, app)
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir1,
		"friend", "add",
		"Ivy (aka Ives) :: Best friend #bestie @Seattle $id:ivy",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir1,
		"location", "add",
		"Seattle, USA :: Emerald City #city $id:seattle",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir1,
		"note", "add",
		"2024-01-15 :: Met Ivy for coffee #meetup @seattle",
	})
	require.NoError(t, err)

	// Export from first journal
	exportFile := filepath.Join(t.TempDir(), "roundtrip.frentxt")
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir1,
		"journal", "export",
		exportFile,
	})
	require.NoError(t, err)

	// Initialize second journal
	jDir2, err := InitJournal(t, app)
	require.NoError(t, err)

	// Import into second journal
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir2,
		"journal", "import",
		exportFile,
	})
	require.NoError(t, err)

	// Export from second journal
	exportFile2 := filepath.Join(t.TempDir(), "roundtrip2.frentxt")
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir2,
		"journal", "export",
		exportFile2,
	})
	require.NoError(t, err)

	// Verify the second export contains the same data
	content, err := os.ReadFile(exportFile2)
	require.NoError(t, err)

	require.Contains(t, string(content), "Ivy")
	require.Contains(t, string(content), "Ives")
	require.Contains(t, string(content), "#bestie")
	require.Contains(t, string(content), "Seattle")
	require.Contains(t, string(content), "USA")
}

func TestJournal_Export_MultipleFilters(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Add all types of data
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"friend", "add",
		"Jack :: Friend #friend $id:jack",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"location", "add",
		"Miami, USA :: Beach city #beach $id:miami",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"note", "add",
		"2024-02-01 :: Note about Jack #note @miami",
	})
	require.NoError(t, err)

	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"activity", "add",
		"2024-02-02 :: Beach day with Jack #activity @miami",
	})
	require.NoError(t, err)

	// Export friends and locations only
	exportFile := filepath.Join(t.TempDir(), "friends_locations.frentxt")
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"journal", "export",
		"--friends",
		"--locations",
		exportFile,
	})
	require.NoError(t, err)

	content, err := os.ReadFile(exportFile)
	require.NoError(t, err)

	require.Contains(t, string(content), "/f")
	require.Contains(t, string(content), "Jack")
	require.Contains(t, string(content), "/l")
	require.Contains(t, string(content), "Miami")
	require.NotContains(t, string(content), "/n")
	require.NotContains(t, string(content), "/act")
}

func TestJournal_Import_EmptyFile(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Create an empty file
	importFile := filepath.Join(t.TempDir(), "empty.frentxt")
	err = os.WriteFile(importFile, []byte(""), 0o644)
	require.NoError(t, err)

	// Import should succeed (with no data)
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"journal", "import",
		importFile,
	})
	require.NoError(t, err)
}

func TestJournal_Export_EmptyJournal(t *testing.T) {
	ctx := t.Context()
	app := cmd.NewApp()

	jDir, err := InitJournal(t, app)
	require.NoError(t, err)

	// Export from empty journal
	exportFile := filepath.Join(t.TempDir(), "empty_export.frentxt")
	err = app.RunContext(ctx, []string{
		"frens", "-j", jDir,
		"journal", "export",
		exportFile,
	})
	require.NoError(t, err)

	// File should exist but be empty or minimal
	require.FileExists(t, exportFile)
}
