package formatter

import (
	"fmt"
	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/lang"
	"github.com/roma-glushko/frens/internal/log"
	"strings"
)

func init() {
	log.RegisterFormatter(log.FormatText, friend.Person{}, PersonTextFormatter{})
}

func wrapText(text string, width int) []string {
	words := strings.Fields(text)
	var lines []string
	var current string
	for _, word := range words {
		if len(current)+len(word)+1 > width {
			lines = append(lines, current)
			current = word
		} else {
			if current != "" {
				current += " "
			}
			current += word
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

func writeField(sb *strings.Builder, key, value string) {
	if strings.TrimSpace(value) == "" {
		return
	}
	sb.WriteString(fmt.Sprintf("%-12s: %s\n", key, value))
}

type PersonTextFormatter struct{}

var _ log.Formatter = (*PersonTextFormatter)(nil)

func (p PersonTextFormatter) FormatSingle(e any) (string, error) {
	person, ok := e.(friend.Person)

	if !ok {
		return "", fmt.Errorf("expected 'friend.Person'")
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%-14s %s\n", "NAME", person.String()))

	if len(person.Tags) > 0 {
		sb.WriteString(fmt.Sprintf("%-14s %s\n", "TAGS", lang.RenderTags(person.Tags)))
	}

	if len(person.Locations) > 0 {
		sb.WriteString(fmt.Sprintf("%-14s %s\n", "LOCATIONS", lang.RenderLocMarkers(person.Locations)))
	}

	if person.Desc != "" {
		sb.WriteString("DESCRIPTION \n\n")
		indent := strings.Repeat(" ", 4)
		wrapped := wrapText(person.Desc, 80-4)
		for _, line := range wrapped {
			sb.WriteString(indent + line + "\n")
		}
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

func (p PersonTextFormatter) FormatList(entities []any) (string, error) {
	//TODO implement me
	panic("implement me")
}
