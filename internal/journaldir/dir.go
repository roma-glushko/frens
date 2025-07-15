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

func Dir(overridePath string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home dir: %w", err)
	}

	journalPath := filepath.Join(homeDir, ".config", DefaultFrensDir)

	if overridePath != "" {
		journalPath = overridePath
	}

	// ensure directory exists
	err = os.MkdirAll(journalPath, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create the journal dir at %s: %w", journalPath, err)
	}

	return journalPath, nil
}
