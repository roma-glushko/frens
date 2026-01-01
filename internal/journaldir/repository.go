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
	"fmt"
	"sync"

	"github.com/roma-glushko/frens/internal/journal"
)

// Repository provides access to journal data with proper transaction handling.
// It encapsulates the loading, saving, and locking logic for journal operations.
type Repository struct {
	dir     string
	journal *journal.Journal
	mu      sync.RWMutex
}

// NewRepository creates a new Repository for the given journal directory.
// It does not load the journal immediately - call Load() or use Update() to load data.
func NewRepository(dir string) *Repository {
	return &Repository{
		dir: dir,
	}
}

// Dir returns the journal directory path.
func (r *Repository) Dir() string {
	return r.dir
}

// IsLoaded returns true if the journal has been loaded into memory.
func (r *Repository) IsLoaded() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.journal != nil
}

// Load loads the journal from disk if not already loaded.
// Returns the cached journal if already loaded.
func (r *Repository) Load() (*journal.Journal, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.journal != nil {
		return r.journal, nil
	}

	jr, err := load(r.dir)
	if err != nil {
		return nil, fmt.Errorf("failed to load journal: %w", err)
	}

	r.journal = jr
	return r.journal, nil
}

// Journal returns the loaded journal for read operations.
// Returns nil if the journal has not been loaded yet.
// For write operations, use Update() instead.
func (r *Repository) Journal() *journal.Journal {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.journal
}

// Update executes a function that modifies the journal and saves changes if dirty.
// This method provides transaction-like semantics: the journal is locked during
// the update, and changes are persisted only if the journal is marked dirty.
func (r *Repository) Update(fn func(*journal.Journal) error) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Load if not yet loaded
	if r.journal == nil {
		jr, err := load(r.dir)
		if err != nil {
			return fmt.Errorf("failed to load journal: %w", err)
		}
		r.journal = jr
	}

	// Execute the update function
	if err := fn(r.journal); err != nil {
		return err
	}

	// Save if dirty
	if !r.journal.IsDirty() {
		return nil
	}

	if err := save(r.journal); err != nil {
		return fmt.Errorf("failed to save journal: %w", err)
	}

	r.journal.SetDirty(false)
	return nil
}

// Reload forces a reload of the journal from disk, discarding any unsaved changes.
func (r *Repository) Reload() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	jr, err := load(r.dir)
	if err != nil {
		return fmt.Errorf("failed to reload journal: %w", err)
	}

	r.journal = jr
	return nil
}
