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
	t := s.Theme
	var sb strings.Builder

	title := lipgloss.NewStyle().
		Foreground(t.Main).
		Bold(true).
		Render("HIGH SCORES")

	sb.WriteString(title)
	sb.WriteString("\n\n")

	if len(m.scores.Scores) == 0 {
		sb.WriteString(lipgloss.NewStyle().
			Foreground(t.Sub).
			Render("No scores yet. Play a game!"))
	} else {
		// Header
		headerStyle := lipgloss.NewStyle().Foreground(t.Sub)
		sb.WriteString(headerStyle.Render(fmt.Sprintf("   %-4s %10s %6s %6s   %s", "#", "Score", "Level", "Lines", "Date")))
		sb.WriteString("\n")
		sb.WriteString(lipgloss.NewStyle().Foreground(t.SubAlt).Render("   " + strings.Repeat("─", 42)))
		sb.WriteString("\n")

		for i, hs := range m.scores.Scores {
			dateStr := hs.Date.Format("2006-01-02")

			var rankStr string
			rankStyle := lipgloss.NewStyle().Foreground(t.Sub)
			scoreStyle := lipgloss.NewStyle().Foreground(t.FG)

			if i == 0 {
				rankStyle = rankStyle.Foreground(t.Main).Bold(true)
				scoreStyle = scoreStyle.Foreground(t.Main).Bold(true)
				rankStr = " 1."
			} else {
				rankStr = fmt.Sprintf("%2d.", i+1)
			}

			sb.WriteString(rankStyle.Render(fmt.Sprintf("   %-4s", rankStr)))
			sb.WriteString(scoreStyle.Render(fmt.Sprintf("%10d", hs.Score)))
			sb.WriteString(lipgloss.NewStyle().Foreground(t.FG).Render(fmt.Sprintf(" %6d", hs.Level)))
			sb.WriteString(lipgloss.NewStyle().Foreground(t.FG).Render(fmt.Sprintf(" %6d", hs.Lines)))
			sb.WriteString(lipgloss.NewStyle().Foreground(t.Sub).Render(fmt.Sprintf("   %s", dateStr)))
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\n")
	sb.WriteString(lipgloss.NewStyle().
		Foreground(t.SubAlt).
		Render("   q back"))

	return sb.String()
}
