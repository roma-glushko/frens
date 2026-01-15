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

package log

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

// BulletChar is the bullet character for list items
const BulletChar = "•"

type OutputHandler = func(w io.Writer, data any)

func TextOutputHandler(w io.Writer, data any) {
	_, _ = fmt.Fprintln(w, data)
}

func JSONOutputHandler(w io.Writer, data any) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		_, _ = fmt.Fprintf(w, "Error formatting JSON output: %v\n", err)
		return
	}

	_, _ = fmt.Fprintln(w, string(jsonData))
}

// Shared styles for CLI output
var (
	SuccessStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("2")) // green
	HeaderStyle   = lipgloss.NewStyle().Bold(true)
	WarnStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("3")) // yellow
	ErrorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("1")) // red
	MutedStyle    = lipgloss.NewStyle().Faint(true)
	CountStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6")) // cyan
	LabelStyle    = lipgloss.NewStyle().Bold(true)
	TagStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("6")) // cyan
	LocationStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5")) // magenta
	IDStyle       = lipgloss.NewStyle().Faint(true)
)

// Success prints a success message with a checkmark
func Success(msg string) {
	Info(SuccessStyle.Render(" ✔") + " " + msg + "\n")
}

// Successf prints a formatted success message with a checkmark
func Successf(format string, args ...any) {
	Success(fmt.Sprintf(format, args...))
}

// Header prints a section header
func Header(msg string) {
	Info("==> " + HeaderStyle.Render(msg) + "\n")
}

// Headerf prints a formatted section header
func Headerf(format string, args ...any) {
	Header(fmt.Sprintf(format, args...))
}

// Empty prints an empty state message
func Empty(entity string) {
	Info(MutedStyle.Render("No "+entity+" found.") + "\n")
}

// Found prints a count of found items
func Found(count int, singular, plural string) {
	word := plural
	if count == 1 {
		word = singular
	}

	Info(fmt.Sprintf("Found %s %s:\n", CountStyle.Render(strconv.Itoa(count)), word))
}

// Progress prints a progress/status message
func Progress(msg string) {
	Info(MutedStyle.Render("› "+msg) + "\n")
}

// Progressf prints a formatted progress/status message
func Progressf(format string, args ...any) {
	Progress(fmt.Sprintf(format, args...))
}

// Canceled prints a cancellation message
func Canceled(msg string) {
	Info(" ↩ " + msg + "\n")
}

// Deleted prints a deletion confirmation
func Deleted(entity string) {
	Info(SuccessStyle.Render(" ✔") + " " + entity + " deleted.\n")
}

// WarnPrompt formats a warning message for prompts
func WarnPrompt(msg string) string {
	return WarnStyle.Render(msg)
}

// ConfirmPrompt formats a confirmation prompt message
func ConfirmPrompt(msg string) string {
	return WarnStyle.Render(msg) + " [y/N]: "
}

// Bullet prints a bulleted list item
func Bullet(msg string) {
	Info(" • " + msg + "\n")
}

// Bulletf prints a formatted bulleted list item
func Bulletf(format string, args ...any) {
	Bullet(fmt.Sprintf(format, args...))
}
