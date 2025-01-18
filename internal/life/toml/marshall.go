package toml

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
)

type Files interface {
	FriendsFile | ActivitiesFile
}

func createFile[T Files](lifeDir string, fileName string) error {
	file, err := os.Create(filepath.Join(lifeDir, fileName))

	if err != nil {
		return err
	}

	defer file.Close()

	var content T

	encoder := toml.NewEncoder(file)

	if err := encoder.Encode(content); err != nil {
		return err
	}

	return nil
}

func Init(lifeDir string) error {
	var errs []error

	if err := createFile[FriendsFile](lifeDir, FileNameFriends); err != nil {
		errs = append(errs, fmt.Errorf("failed to create friends file: %w", err))
	}

	if err := createFile[ActivitiesFile](lifeDir, FileNameActivities); err != nil {
		errs = append(errs, fmt.Errorf("failed to create activities file: %w", err))
	}

	return errors.Join(errs...)
}

func Load(lifeDir string) {

}
