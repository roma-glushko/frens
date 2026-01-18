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

package lang

import (
	"bufio"
	"strings"

	"github.com/roma-glushko/frens/internal/friend"
)

type blockType string

const (
	blockTypeFriend   blockType = "/f"
	blockTypeLocation blockType = "/l"
	blockTypeNote     blockType = "/n"
	blockTypeActivity blockType = "/act"
)

// FrenTXTImport holds the parsed results from a FrenTXT file.
type FrenTXTImport struct {
	Friends    []friend.Person
	Locations  []friend.Location
	Notes      []friend.Event
	Activities []friend.Event
}

// ParseFrenTXT parses content in FrenTXT format and returns structured data.
func ParseFrenTXT(content string) (*FrenTXTImport, error) {
	result := &FrenTXTImport{
		Friends:    make([]friend.Person, 0),
		Locations:  make([]friend.Location, 0),
		Notes:      make([]friend.Event, 0),
		Activities: make([]friend.Event, 0),
	}

	blocks := splitIntoBlocks(content)

	for _, block := range blocks {
		if err := parseBlock(block, result); err != nil {
			// Skip blocks that fail to parse rather than failing entirely
			continue
		}
	}

	return result, nil
}

type rawBlock struct {
	blockType blockType
	content   string
}

func splitIntoBlocks(content string) []rawBlock {
	var blocks []rawBlock
	var currentType blockType
	var currentLines []string

	scanner := bufio.NewScanner(strings.NewReader(content))

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Check if this line is a block marker
		if newType := detectBlockType(trimmed); newType != "" {
			// Save the previous block if it exists
			if currentType != "" && len(currentLines) > 0 {
				blocks = append(blocks, rawBlock{
					blockType: currentType,
					content:   strings.Join(currentLines, "\n"),
				})
			}

			currentType = newType
			currentLines = nil

			continue
		}

		// Accumulate lines for the current block
		if currentType != "" {
			currentLines = append(currentLines, line)
		}
	}

	// Don't forget the last block
	if currentType != "" && len(currentLines) > 0 {
		blocks = append(blocks, rawBlock{
			blockType: currentType,
			content:   strings.Join(currentLines, "\n"),
		})
	}

	return blocks
}

func detectBlockType(line string) blockType {
	switch line {
	case string(blockTypeFriend):
		return blockTypeFriend
	case string(blockTypeLocation):
		return blockTypeLocation
	case string(blockTypeNote):
		return blockTypeNote
	case string(blockTypeActivity):
		return blockTypeActivity
	default:
		return ""
	}
}

func parseBlock(block rawBlock, result *FrenTXTImport) error {
	// Normalize multi-line content: join lines that are part of the description
	content := normalizeBlockContent(block.content)

	if content == "" {
		return ErrNoInfo
	}

	switch block.blockType {
	case blockTypeFriend:
		person, err := ExtractPerson(content)
		if err != nil {
			return err
		}

		result.Friends = append(result.Friends, person)

	case blockTypeLocation:
		location, err := ExtractLocation(content)
		if err != nil {
			return err
		}

		result.Locations = append(result.Locations, location)

	case blockTypeNote:
		event, err := ExtractEvent(friend.EventTypeNote, content)
		if err != nil {
			return err
		}

		result.Notes = append(result.Notes, event)

	case blockTypeActivity:
		event, err := ExtractEvent(friend.EventTypeActivity, content)
		if err != nil {
			return err
		}

		result.Activities = append(result.Activities, event)
	}

	return nil
}

// normalizeBlockContent handles multi-line content by:
// 1. Preserving the first line as-is (header with name/date)
// 2. Joining subsequent lines to extend the description
func normalizeBlockContent(content string) string {
	lines := strings.Split(content, "\n")

	if len(lines) == 0 {
		return ""
	}

	// Find the separator position in the first line
	firstLine := strings.TrimSpace(lines[0])

	if len(lines) == 1 {
		return firstLine
	}

	// Check if the first line contains the separator
	sepIdx := strings.Index(firstLine, Separator)

	var result strings.Builder

	if sepIdx >= 0 {
		// First line has separator - append additional lines to description
		result.WriteString(firstLine)

		for i := 1; i < len(lines); i++ {
			line := strings.TrimSpace(lines[i])
			if line != "" {
				result.WriteString(" ")
				result.WriteString(line)
			}
		}
	} else {
		// First line has no separator (just name/date)
		// Look for separator in subsequent lines or treat as description continuation
		result.WriteString(firstLine)

		for i := 1; i < len(lines); i++ {
			line := strings.TrimSpace(lines[i])
			if line != "" {
				// If this line starts with "::", it's the description start
				if strings.HasPrefix(line, Separator) {
					result.WriteString(" ")
					result.WriteString(line)
				} else {
					// Otherwise, check if we already have a separator
					if strings.Contains(result.String(), Separator) {
						result.WriteString(" ")
						result.WriteString(line)
					} else {
						// Add as part of description with separator
						result.WriteString(" ")
						result.WriteString(Separator)
						result.WriteString(" ")
						result.WriteString(line)
					}
				}
			}
		}
	}

	return result.String()
}

// RenderFrenTXT renders journal data to FrenTXT format.
func RenderFrenTXT(data *FrenTXTImport) string {
	var sb strings.Builder

	// Render friends
	for i, f := range data.Friends {
		if i > 0 || sb.Len() > 0 {
			sb.WriteString("\n")
		}

		sb.WriteString(string(blockTypeFriend))
		sb.WriteString("\n")
		sb.WriteString(RenderPerson(f))
		sb.WriteString("\n")
	}

	// Render locations
	for _, l := range data.Locations {
		if sb.Len() > 0 {
			sb.WriteString("\n")
		}

		sb.WriteString(string(blockTypeLocation))
		sb.WriteString("\n")
		sb.WriteString(RenderLocation(l))
		sb.WriteString("\n")
	}

	// Render notes
	for _, n := range data.Notes {
		if sb.Len() > 0 {
			sb.WriteString("\n")
		}

		sb.WriteString(string(blockTypeNote))
		sb.WriteString("\n")
		sb.WriteString(RenderEvent(n))
		sb.WriteString("\n")
	}

	// Render activities
	for _, a := range data.Activities {
		if sb.Len() > 0 {
			sb.WriteString("\n")
		}

		sb.WriteString(string(blockTypeActivity))
		sb.WriteString("\n")
		sb.WriteString(RenderEvent(a))
		sb.WriteString("\n")
	}

	return sb.String()
}
