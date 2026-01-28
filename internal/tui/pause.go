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
	var sb strings.Builder

	title := lipgloss.NewStyle().
		Foreground(s.Theme.Main).
		Bold(true).
		Render("PAUSED")

	sb.WriteString(title)
	sb.WriteString("\n\n")
	sb.WriteString(lipgloss.NewStyle().Foreground(s.Theme.FG).Render("p/esc  resume"))
	sb.WriteString("\n")
	sb.WriteString(lipgloss.NewStyle().Foreground(s.Theme.FG).Render("r      restart"))
	sb.WriteString("\n")
	sb.WriteString(lipgloss.NewStyle().Foreground(s.Theme.FG).Render("q      quit to menu"))

	return lipgloss.NewStyle().
		Padding(2, 4).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(s.Theme.Sub).
		Render(sb.String())
}
