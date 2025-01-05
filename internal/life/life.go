package life

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const DefaultFrensDir = ".frens"

type Life struct {
}

func DefaultDir() (string, error) {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	lifePath := filepath.Join(homeDir, DefaultFrensDir)

	// ensure directory exists
	err = os.MkdirAll(lifePath, os.ModePerm)

	if err != nil {
		return "", fmt.Errorf("failed to create the life directory at %s: %w", lifePath, err)
	}

	return lifePath, nil
}

// Init loads the life from the specific path or `~/.frens/` is used by default
func Init() (*Life, error) {
	lifePath, err := DefaultDir()

	if err != nil {
		return nil, err
	}

	_, err = os.Stat(lifePath)

	if err != nil && errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("life already exists at %s", lifePath)
	}

	err = os.MkdirAll(lifePath, os.ModePerm)

	if err != nil {
		return nil, fmt.Errorf("failed to create the life directory at %s: %w", lifePath, err)
	}

	return nil, nil
}

// Load loads the life from the specific path or `~/.frens/` is used by default
func Load() (*Life, error) {
	// TODO: implement

	return nil, nil
}
