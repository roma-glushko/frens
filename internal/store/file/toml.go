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

package file

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/store"
	"github.com/roma-glushko/frens/internal/tag"
)

type Files interface {
	FriendsFile | EventsFile
}

const (
	FileNameFriends    = "friends.toml"
	FileNameActivities = "activities.toml"
)

type FriendsFile struct {
	Tags      []tag.Tag          `toml:"tags"`
	Friends   []*friend.Person   `toml:"friends"`
	Locations []*friend.Location `toml:"locations"`
}

type EventsFile struct {
	Activities []*friend.Event `toml:"activities"`
	Notes      []*friend.Event `toml:"notes"`
}

type TOMLFileStore struct {
	dir string
	mu  sync.Mutex // TODO: use file locking instead
}

var _ store.Store = (*TOMLFileStore)(nil)

func NewTOMLFileStore(dir string) *TOMLFileStore {
	return &TOMLFileStore{
		dir: dir,
	}
}

func (s *TOMLFileStore) Init(ctx context.Context) error {
	var errs []error

	var entities FriendsFile

	if err := saveFile(ctx, s.dir, FileNameFriends, entities); err != nil {
		errs = append(errs, fmt.Errorf("failed to create friends file: %w", err))
	}

	var activities EventsFile

	if err := saveFile(ctx, s.dir, FileNameActivities, activities); err != nil {
		errs = append(errs, fmt.Errorf("failed to create activities file: %w", err))
	}

	return errors.Join(errs...)
}

func (s *TOMLFileStore) Exist(ctx context.Context) bool {
	friendsFilePath := filepath.Join(s.dir, FileNameFriends)
	activitiesFilePath := filepath.Join(s.dir, FileNameActivities)

	if _, err := os.Stat(friendsFilePath); err == nil {
		return true
	}

	if _, err := os.Stat(activitiesFilePath); err == nil {
		return true
	}

	return false
}

func (s *TOMLFileStore) Load(ctx context.Context) (*journal.Journal, error) {
	var errs []error

	entities, err := loadFile[FriendsFile](ctx, filepath.Join(s.dir, FileNameFriends))
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to load friends file: %w", err))
	}

	events, err := loadFile[EventsFile](ctx, filepath.Join(s.dir, FileNameActivities))
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to load activities file: %w", err))
	}

	if len(errs) != 0 {
		return nil, errors.Join(errs...)
	}

	j := &journal.Journal{
		Tags:       entities.Tags,
		Friends:    entities.Friends,
		Locations:  entities.Locations,
		Activities: events.Activities,
		Notes:      events.Notes,
	}

	j.Init()

	return j, nil
}

func (s *TOMLFileStore) Save(ctx context.Context, j *journal.Journal) error {
	var errs []error

	entities := FriendsFile{
		Tags:      j.Tags,
		Friends:   j.Friends,
		Locations: j.Locations,
	}

	if err := saveFile(ctx, s.dir, FileNameFriends, entities); err != nil {
		errs = append(errs, fmt.Errorf("failed to create friends file: %w", err))
	}

	events := EventsFile{
		Notes:      j.Notes,
		Activities: j.Activities,
	}

	if err := saveFile(ctx, s.dir, FileNameActivities, events); err != nil {
		errs = append(errs, fmt.Errorf("failed to create events file: %w", err))
	}

	return errors.Join(errs...)
}

func (s *TOMLFileStore) Tx(ctx context.Context, fn store.JournalUpdater) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	j, err := s.Load(ctx)
	if err != nil {
		return fmt.Errorf("failed to load journal: %w", err)
	}

	if err := fn(j); err != nil {
		return fmt.Errorf("failed to execute transaction function: %w", err)
	}

	if !j.IsDirty() {
		return nil
	}

	if err := s.Save(ctx, j); err != nil {
		return fmt.Errorf("failed to save journal: %w", err)
	}

	j.SetDirty(false)

	return nil
}

func saveFile[T Files](_ context.Context, dirPath, fileName string, content T) (err error) {
	path := filepath.Join(dirPath, fileName)

	tmpFile, err := os.CreateTemp(dirPath, fileName+".*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create a temp file for %s: %w", path, err)
	}

	encoder := toml.NewEncoder(tmpFile)

	if err = encoder.Encode(content); err != nil {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())

		return fmt.Errorf("failed to encode content of %s: %w", path, err)
	}

	if err := tmpFile.Sync(); err != nil {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())

		return fmt.Errorf("failed to sync temp file for %s: %w", path, err)
	}

	if err := tmpFile.Close(); err != nil {
		_ = os.Remove(tmpFile.Name())

		return fmt.Errorf("failed to close temp file for %s: %w", path, err)
	}

	if err := os.Rename(tmpFile.Name(), path); err != nil {
		_ = os.Remove(tmpFile.Name())

		return fmt.Errorf("failed to rename temp file to %s: %w", path, err)
	}

	return nil
}

func loadFile[T Files](_ context.Context, filePath string) (c *T, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close file %s: %w", filePath, closeErr)
		}
	}()

	var content T

	decoder := toml.NewDecoder(file)

	if _, err = decoder.Decode(&content); err != nil {
		return nil, err
	}

	return &content, nil
}
