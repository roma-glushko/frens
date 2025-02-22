package lifedir

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/roma-glushko/frens/internal/lifedir/toml"

	"github.com/roma-glushko/frens/internal/life"
)

// Init loads the life from the specific path or `~/.frens/` is used by default
func Init(lifePath string) error {
	_, err := os.Stat(lifePath)

	if err != nil && errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("life already exists at %s", lifePath)
	}

	err = os.MkdirAll(lifePath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create the life directory at %s: %w", lifePath, err)
	}

	err = toml.Init(lifePath)
	if err != nil {
		return fmt.Errorf("failed to initialize life: %w", err)
	}

	return nil
}

// Load loads the life from the specific path or `~/.frens/` is used by default
func Load(lifePath string) (*life.Data, error) {
	_, err := os.Stat(lifePath)

	if err != nil && errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("couldn't find life space at %s. Please init a life space via the init command", lifePath)
	}

	data, err := toml.Load(lifePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load life space: %w", err)
	}

	data.Init()

	return data, nil
}

// Save saves the life files from the specific path or `~/.frens/` is used by default
func Save(lifePath string, data *life.Data) error {
	return toml.Save(lifePath, data)
}
