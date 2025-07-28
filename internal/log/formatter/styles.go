package formatter

import "github.com/charmbracelet/lipgloss"

var (
	// symbolStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("2")) // green
	labelStyle    = lipgloss.NewStyle().Bold(true)
	tagStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("6")) // cyan
	locationStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5")) // magenta
	countLabel    = lipgloss.NewStyle().
			Faint(true).
			PaddingLeft(1).
			Foreground(lipgloss.AdaptiveColor{Light: "#666666", Dark: "#999999"})
)
