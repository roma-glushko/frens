package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/roma-glushko/frens/internal/friend"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const maxWidth = 60

var (
	red    = lipgloss.AdaptiveColor{Light: "#FE5F86", Dark: "#FE5F86"}
	indigo = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	green  = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
)

type Styles struct {
	Base,
	HeaderText,
	Status,
	StatusHeader,
	Highlight,
	ErrorHeaderText,
	Help lipgloss.Style
}

func NewStyles(lg *lipgloss.Renderer) *Styles {
	s := Styles{}

	s.Base = lg.NewStyle().
		Padding(1, 4, 0, 1)
	s.HeaderText = lg.NewStyle().
		Foreground(indigo).
		Bold(true).
		Padding(0, 1, 0, 2)
	s.Status = lg.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(indigo).
		PaddingLeft(1).
		MarginTop(1)
	s.StatusHeader = lg.NewStyle().
		Foreground(green).
		Bold(true)
	s.Highlight = lg.NewStyle().
		Foreground(lipgloss.Color("212"))
	s.ErrorHeaderText = s.HeaderText.
		Foreground(red)
	s.Help = lg.NewStyle().
		Foreground(lipgloss.Color("240"))

	return &s
}

type NewFriendModel struct {
	friend    *friend.Friend
	nicknames string
	lg        *lipgloss.Renderer
	styles    *Styles
	form      *huh.Form
	width     int
}

func NewFriendForm(f *friend.Friend) NewFriendModel {
	m := NewFriendModel{
		friend: f,
		width:  maxWidth,
	}

	if f.Nicknames != nil {
		m.nicknames = strings.Join(f.Nicknames, ", ")
	}

	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("Name").
				Description("What's your friend full name?").
				Validate(huh.ValidateMinLength(1)).
				Placeholder("Jim Halpert").
				Value(&f.Name),

			huh.NewInput().
				Key("nicknames").
				Title("Nicknames").
				Description("What's your friend nicknames?").
				Placeholder("Jimbo, Jimmy, Jimothy, Tuna").
				Value(&m.nicknames),

			huh.NewMultiSelect[string]().
				Key("locations").
				Title("Locations").
				Description("What locations would you like to associate to your friend with?").
				Options(huh.NewOptions("Scranton, Pennsylvania", "New York")...).
				Value(&f.Locations),

			huh.NewMultiSelect[string]().
				Key("tags").
				Title("Tags").
				Description("What tags would you like to add to your friend?").
				Options(huh.NewOptions("friend", "family", "colleague", "acquaintance")...).
				Value(&f.Tags),

			huh.NewConfirm().
				Key("done").
				Title("All done?").
				Validate(func(v bool) error {
					if !v {
						return fmt.Errorf("Welp, whenever you are ready")
					}
					return nil
				}).
				Affirmative("Yep").
				Negative("No"),
		),
	).
		WithShowHelp(false).
		WithShowErrors(false)

	return m
}

func (m NewFriendModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m NewFriendModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	var cmds []tea.Cmd

	// Process the form
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	if m.form.State == huh.StateCompleted {
		// Quit when the form is done.
		m.friend.Nicknames = strings.Split(m.nicknames, ",")

		cmds = append(cmds, tea.Quit)
	}

	return m, tea.Batch(cmds...)
}

func (m NewFriendModel) View() string {
	s := m.styles

	v := strings.TrimSuffix(m.form.View(), "\n\n")
	form := m.lg.NewStyle().Margin(1, 0).Render(v)

	errors := m.form.Errors()

	header := m.appBoundaryView("New Friend")

	if len(errors) > 0 {
		header = m.appErrorBoundaryView(m.errorView())
	}

	footer := m.appBoundaryView(m.form.Help().ShortHelpView(m.form.KeyBinds()))

	if len(errors) > 0 {
		footer = m.appErrorBoundaryView("")
	}

	return s.Base.Render(header + "\n" + form + "\n\n" + footer)
}

func (m NewFriendModel) errorView() string {
	var s string
	for _, err := range m.form.Errors() {
		s += err.Error()
	}
	return s
}

func (m NewFriendModel) appBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.HeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(indigo),
	)
}

func (m NewFriendModel) appErrorBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.ErrorHeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(red),
	)
}
