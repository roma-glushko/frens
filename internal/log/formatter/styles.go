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

package formatter

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/roma-glushko/frens/internal/log"
)

// Re-export shared styles from log package for backward compatibility
var (
	labelStyle    = log.LabelStyle
	tagStyle      = log.TagStyle
	locationStyle = log.LocationStyle
	idStyle       = log.IDStyle
	countLabel    = lipgloss.NewStyle().
			Faint(true).
			PaddingLeft(1).
			Foreground(lipgloss.AdaptiveColor{Light: "#666666", Dark: "#999999"})
	friendStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("3"))
)
