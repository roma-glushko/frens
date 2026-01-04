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

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type EditorOptions struct {
	Title       string
	Placeholder string
	SyntaxHint  string
}

type errMsg error

type EditorForm struct {
	Title      string
	SyntaxHint string
	Textarea   textarea.Model
	err        error
}

func NewEditorForm(o EditorOptions) EditorForm {
	ti := textarea.New()
	ti.ShowLineNumbers = false
	ti.Placeholder = o.Placeholder
	ti.SetWidth(100)
	ti.SetHeight(10)
	ti.Focus()

	return EditorForm{
		Title:      o.Title,
		Textarea:   ti,
		err:        nil,
		SyntaxHint: o.SyntaxHint,
	}
}

func (m EditorForm) Init() tea.Cmd {
	return textarea.Blink
}

func (m EditorForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type { //nolint:exhaustive // intentionally not exhaustive
		case tea.KeyEsc:
			if m.Textarea.Focused() {
				m.Textarea.Blur()
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		default:
			if !m.Textarea.Focused() {
				cmd = m.Textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.Textarea, cmd = m.Textarea.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m EditorForm) View() string {
	return fmt.Sprintf(
		"\n%s\n\n%s\n\n%s\n\n%s",
		m.Title,
		m.Textarea.View(),
		"Syntax: "+m.SyntaxHint,
		"(ctrl+c to quit)",
	) + "\n\n"
}
