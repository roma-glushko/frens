// Copyright 2026 Roma Hlushko
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

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/BurntSushi/toml"
	"github.com/roma-glushko/frens/internal/log"
)

const ConfigFile = "config.toml"

// Config is the main application configuration
type Config struct {
	Notifications Notifications `toml:"notifications"`
}

// Load loads the configuration from the config directory
func Load(configDir string) (*Config, error) {
	configPath := filepath.Join(configDir, ConfigFile)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{}, nil
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config: %w", err)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Error(fmt.Sprintf("failed to close config file: %v", cerr))
		}
	}()

	var cfg Config

	decoder := toml.NewDecoder(file)
	if _, err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	// Sort notification rules by priority
	sort.SliceStable(cfg.Notifications.Rules, func(i, j int) bool {
		return cfg.Notifications.Rules[i].Priority < cfg.Notifications.Rules[j].Priority
	})

	return &cfg, nil
}

// Save saves the configuration to the config directory
func Save(configDir string, cfg *Config) error {
	configPath := filepath.Join(configDir, ConfigFile)

	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	defer func() {
		_ = file.Close()
	}()

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	return nil
}
