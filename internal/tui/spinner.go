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

package tui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty"
)

var spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6")) // cyan

// TaskResult holds the result of an async task
type TaskResult struct {
	Err error
}

// taskDoneMsg signals that the task is complete
type taskDoneMsg struct {
	err error
}

// SpinnerModel is a bubble tea model for showing a spinner during an operation
type SpinnerModel struct {
	spinner spinner.Model
	message string
	task    func() error
	done    bool
	err     error
}

// NewSpinner creates a new spinner model
func NewSpinner(message string, task func() error) SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	return SpinnerModel{
		spinner: s,
		message: message,
		task:    task,
	}
}

func (m SpinnerModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.runTask(),
	)
}

func (m SpinnerModel) runTask() tea.Cmd {
	return func() tea.Msg {
		err := m.task()
		return taskDoneMsg{err: err}
	}
}

func (m SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case taskDoneMsg:
		m.done = true
		m.err = msg.err

		return m, tea.Quit

	case spinner.TickMsg:
		var cmd tea.Cmd

		m.spinner, cmd = m.spinner.Update(msg)

		return m, cmd
	}

	return m, nil
}

func (m SpinnerModel) View() string {
	if m.done {
		return ""
	}

	return fmt.Sprintf("%s %s", m.spinner.View(), m.message)
}

// Error returns the error from the task, if any
func (m SpinnerModel) Error() error {
	return m.err
}

// RunWithSpinner runs a task with a spinner if the terminal is interactive
// Falls back to simple progress message if not a TTY
func RunWithSpinner(message string, task func() error) error {
	// Check if we're in a TTY
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		fmt.Printf("› %s\n", message)
		return task()
	}

	model := NewSpinner(message, task)
	p := tea.NewProgram(model)

	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	if m, ok := finalModel.(SpinnerModel); ok {
		return m.Error()
	}

	return nil
}
