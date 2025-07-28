package formatter

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/lang"
	"github.com/roma-glushko/frens/internal/log"
	"strings"
)

var (
	symbolStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("2")) // green
	labelStyle    = lipgloss.NewStyle().Bold(true)
	valueStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("7")) // gray
	descStyle     = lipgloss.NewStyle().PaddingLeft(4)
	tagStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("6")) // cyan
	locationStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5")) // magenta
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

	sb.WriteString(fmt.Sprintf(" %s [%s]", labelStyle.Render(person.String()), person.ID))
	sb.WriteString("\n")

	if len(person.Tags) > 0 {
		sb.WriteString(" " + tagStyle.Render(lang.RenderTags(person.Tags)))
		sb.WriteString(" ")
	}

	if len(person.Locations) > 0 {
		sb.WriteString(locationStyle.Render(lang.RenderLocMarkers(person.Locations)))
	}

	if person.Desc != "" {
		sb.WriteString("\n")
		wrapped := wrapText(person.Desc, 80)
		for _, line := range wrapped {
			sb.WriteString(" " + line + "\n")
		}
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

func (p PersonTextFormatter) FormatList(entities []any) (string, error) {
	//TODO implement me
	panic("implement me")
}
