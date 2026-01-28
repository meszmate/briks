package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// PauseModel represents the pause overlay.
type PauseModel struct{}

// NewPauseModel creates a new pause model.
func NewPauseModel() PauseModel {
	return PauseModel{}
}

// View renders the pause screen.
func (p PauseModel) View(s Styles) string {
	t := s.Theme
	var sb strings.Builder

	title := lipgloss.NewStyle().
		Foreground(t.Main).
		Bold(true).
		Render("PAUSED")

	sb.WriteString(title)
	sb.WriteString("\n\n")

	dimStyle := lipgloss.NewStyle().Foreground(t.Sub)
	keyStyle := lipgloss.NewStyle().Foreground(t.FG)

	sb.WriteString(keyStyle.Render("p") + dimStyle.Render(" resume"))
	sb.WriteString("\n")
	sb.WriteString(keyStyle.Render("r") + dimStyle.Render(" restart"))
	sb.WriteString("\n")
	sb.WriteString(keyStyle.Render("q") + dimStyle.Render(" quit"))

	return lipgloss.NewStyle().
		Padding(1, 3).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(t.SubAlt).
		Render(sb.String())
}
