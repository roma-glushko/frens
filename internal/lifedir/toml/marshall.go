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

package toml

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/roma-glushko/frens/internal/life"
)

type Files interface {
	FriendsFile | ActivitiesFile
}

func saveFile[T Files](filePath string, content T) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to save file %s: %w", filePath, err)
	}

	defer file.Close()

	encoder := toml.NewEncoder(file)

	if err = encoder.Encode(content); err != nil {
		return fmt.Errorf("failed to encode content: %w", err)
	}

	return nil
}

func loadFile[T Files](filePath string) (*T, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}

	defer file.Close()

	var content T

	decoder := toml.NewDecoder(file)

	if _, err = decoder.Decode(&content); err != nil {
		return nil, err
	}

	return &content, nil
}

func Init(lifeDir string) error {
	var errs []error

	var entities FriendsFile

	if err := saveFile(filepath.Join(lifeDir, FileNameFriends), entities); err != nil {
		errs = append(errs, fmt.Errorf("failed to create friends file: %w", err))
	}

	var activities ActivitiesFile

	if err := saveFile(filepath.Join(lifeDir, FileNameActivities), activities); err != nil {
		errs = append(errs, fmt.Errorf("failed to create activities file: %w", err))
	}

	return errors.Join(errs...)
}

func Load(lifeDir string) (*life.Data, error) {
	var errs []error

	entities, err := loadFile[FriendsFile](filepath.Join(lifeDir, FileNameFriends))
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to load friends file: %w", err))
	}

	activities, err := loadFile[ActivitiesFile](filepath.Join(lifeDir, FileNameActivities))
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to load activities file: %w", err))
	}

	if len(errs) != 0 {
		return nil, errors.Join(errs...)
	}

	return &life.Data{
		DirPath:    lifeDir,
		Tags:       entities.Tags,
		Friends:    entities.Friends,
		Locations:  entities.Locations,
		Activities: activities.Activities,
	}, nil
}

func Save(l *life.Data) error {
	var errs []error

	entities := FriendsFile{
		Tags:      l.Tags,
		Friends:   l.Friends,
		Locations: l.Locations,
	}

	if err := saveFile(filepath.Join(l.DirPath, FileNameFriends), entities); err != nil {
		errs = append(errs, fmt.Errorf("failed to create friends file: %w", err))
	}

	activities := ActivitiesFile{
		Activities: l.Activities,
	}

	if err := saveFile(filepath.Join(l.DirPath, FileNameActivities), activities); err != nil {
		errs = append(errs, fmt.Errorf("failed to create activities file: %w", err))
	}

	return errors.Join(errs...)
}
