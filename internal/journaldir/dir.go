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
	"os"
	"path/filepath"
)

const DefaultFrensDir = "frens"

func DefaultDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	lifePath := filepath.Join(configDir, DefaultFrensDir)

	// ensure directory exists
	err = os.MkdirAll(lifePath, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create the life directory at %s: %w", lifePath, err)
	}

	return lifePath, nil
}
