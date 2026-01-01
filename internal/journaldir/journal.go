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

package journaldir

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/journaldir/toml"
)

// Init loads the journal from the specific path or `~/.config/frens/` is used by default
func Init(path string) error {
	_, err := os.Stat(path)
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("life already exists at %s", path)
	}

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create the life directory at %s: %w", path, err)
	}

	err = toml.Init(path)
	if err != nil {
		return fmt.Errorf("failed to initialize life: %w", err)
	}

	return nil
}

func Exists(path string) bool {
	return toml.Exists(path)
}

// load loads the journal from the specific path or `~/.config/frens/` is used by default
func load(path string) (*journal.Journal, error) {
	_, err := os.Stat(path)
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf(
			"couldn't find life space at %s. Please init a life space via the init command",
			path,
		)
	}

	data, err := toml.Load(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load journal directory: %w", err)
	}

	// Convert toml.Data to journal.Journal
	jr := &journal.Journal{
		DirPath:    data.DirPath,
		Tags:       data.Tags,
		Friends:    data.Friends,
		Locations:  data.Locations,
		Activities: data.Activities,
		Notes:      data.Notes,
	}

	jr.Init()

	return jr, nil
}

// save saves the life files from the specific path or `~/.config/frens/` is used by default
func save(jr *journal.Journal) error {
	// Convert journal.Journal to toml.Data for saving
	data := &toml.Data{
		DirPath:    jr.DirPath,
		Tags:       jr.Tags,
		Friends:    jr.Friends,
		Locations:  jr.Locations,
		Activities: jr.Activities,
		Notes:      jr.Notes,
	}

	return toml.Save(data)
}
