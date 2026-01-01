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
	"github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/journaldir/toml"
)

// Repository provides access to journal data with proper transaction handling.
type Repository interface {
	Dir() string
	IsLoaded() bool
	Load() (*journal.Journal, error)
	Journal() *journal.Journal
	Update(fn func(*journal.Journal) error) error
	Reload() error
}

// NewRepository creates a new Repository for the given journal directory.
func NewRepository(dir string) Repository {
	return toml.NewRepository(dir)
}

// Exists checks if a journal exists at the given path.
func Exists(path string) bool {
	return toml.Exists(path)
}

// Init initializes a new journal at the given path.
func Init(path string) error {
	return toml.Init(path)
}
