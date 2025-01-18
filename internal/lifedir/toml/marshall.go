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

func createFile[T Files](filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}

	defer file.Close()

	var content T

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

	if err := createFile[FriendsFile](filepath.Join(lifeDir, FileNameFriends)); err != nil {
		errs = append(errs, fmt.Errorf("failed to create friends file: %w", err))
	}

	if err := createFile[ActivitiesFile](filepath.Join(lifeDir, FileNameActivities)); err != nil {
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
		Tags:       entities.Tags,
		Friends:    entities.Friends,
		Locations:  entities.Locations,
		Activities: activities.Activities,
	}, nil
}
