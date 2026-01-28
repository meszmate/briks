package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/meszmate/briks/internal/config"
)

// HighScoresModel displays the high scores table.
type HighScoresModel struct {
	scores *config.HighScores
}

// NewHighScoresModel creates a new high scores model.
func NewHighScoresModel(hs *config.HighScores, s Styles) HighScoresModel {
	return HighScoresModel{scores: hs}
}

// View renders the high scores table.
func (m HighScoresModel) View(s Styles) string {
	var sb strings.Builder

	title := lipgloss.NewStyle().
		Foreground(s.Theme.Main).
		Bold(true).
		Render("HIGH SCORES")

	sb.WriteString(title)
	sb.WriteString("\n\n")

	if len(m.scores.Scores) == 0 {
		sb.WriteString(lipgloss.NewStyle().
			Foreground(s.Theme.Sub).
			Render("No scores yet. Play a game!"))
	} else {
		// Header.
		headerStyle := lipgloss.NewStyle().Foreground(s.Theme.Sub).Bold(true)
		sb.WriteString(headerStyle.Render(fmt.Sprintf("  %-4s  %-10s  %-6s  %-6s  %-12s", "#", "Score", "Level", "Lines", "Date")))
		sb.WriteString("\n")
		sb.WriteString(lipgloss.NewStyle().Foreground(s.Theme.SubAlt).Render(strings.Repeat("─", 48)))
		sb.WriteString("\n")

		for i, hs := range m.scores.Scores {
			rankStyle := lipgloss.NewStyle().Foreground(s.Theme.Sub)
			scoreStyle := lipgloss.NewStyle().Foreground(s.Theme.FG)
			dateStr := hs.Date.Format("2006-01-02")

			if i == 0 {
				rankStyle = rankStyle.Foreground(s.Theme.Main).Bold(true)
				scoreStyle = scoreStyle.Foreground(s.Theme.Main).Bold(true)
			}

			sb.WriteString(rankStyle.Render(fmt.Sprintf("  %-4d", i+1)))
			sb.WriteString(scoreStyle.Render(fmt.Sprintf("  %-10d", hs.Score)))
			sb.WriteString(lipgloss.NewStyle().Foreground(s.Theme.FG).Render(fmt.Sprintf("  %-6d", hs.Level)))
			sb.WriteString(lipgloss.NewStyle().Foreground(s.Theme.FG).Render(fmt.Sprintf("  %-6d", hs.Lines)))
			sb.WriteString(lipgloss.NewStyle().Foreground(s.Theme.Sub).Render(fmt.Sprintf("  %-12s", dateStr)))
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\n")
	sb.WriteString(lipgloss.NewStyle().
		Foreground(s.Theme.Sub).
		Faint(true).
		Render("q/esc to go back"))

	return lipgloss.NewStyle().
		Padding(2, 4).
		Render(sb.String())
}
