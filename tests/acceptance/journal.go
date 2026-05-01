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
	"context"
	"path/filepath"
	"testing"

	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/store/file"
	"github.com/urfave/cli/v2"
)

func InitJournal(t *testing.T, c cli.App) (string, error) {
	jDir := t.TempDir()
	a := []string{
		"frens",
		"-j",
		jDir,
		"journal",
		"init",
	}

	return jDir, c.RunContext(t.Context(), a)
}

// LoadReminders loads the reminders from the journal directory for test assertions
func LoadReminders(t *testing.T, jDir string) []*friend.Reminder {
	t.Helper()

	store := file.NewTOMLFileStore(jDir)

	j, err := store.Load(context.Background())
	if err != nil {
		t.Fatalf("failed to load journal: %v", err)
	}

	return j.Reminders
}

// RemindersFileExists checks that reminders.toml exists in the journal directory
func RemindersFileExists(jDir string) string {
	return filepath.Join(jDir, file.FileNameReminders)
}
