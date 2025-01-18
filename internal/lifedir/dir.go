package lifedir

import (
	"fmt"
	"os"
	"path/filepath"
)

const DefaultFrensDir = ".frens"

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
