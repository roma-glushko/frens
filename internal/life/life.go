package life

import (
	"errors"
	"fmt"
	"github.com/roma-glushko/frens/internal/life/toml"
	"io/fs"
	"os"
)

type Life struct {
}

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
func Load(lifePath string) (*Life, error) {
	// TODO: implement

	return nil, nil
}
