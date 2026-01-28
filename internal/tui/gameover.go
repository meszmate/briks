package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/meszmate/briks/internal/config"
	"github.com/meszmate/briks/internal/game"
)

// GameOverModel represents the game over screen.
type GameOverModel struct {
	score    int
	level    int
	lines    int
	pieces   int
	rank     int // 0 = didn't make high scores
	isNewHS  bool
}

// NewGameOverModel creates a game over model and saves the score.
func NewGameOverModel(engine *game.Engine, hs *config.HighScores) GameOverModel {
	m := GameOverModel{
		score:  engine.Scorer.Score,
		level:  engine.Scorer.Level,
		lines:  engine.Scorer.Lines,
		pieces: engine.PiecesPlaced,
	}

	m.isNewHS = hs.IsHighScore(m.score)
	if m.isNewHS {
		entry := config.HighScore{
			Score:  m.score,
			Level:  m.level,
			Lines:  m.lines,
			Pieces: m.pieces,
			Date:   time.Now(),
		}
		m.rank = hs.Add(entry)
		hs.Save()
	}

	return m
}

// View renders the game over screen.
func (m GameOverModel) View(s Styles) string {
	var sb strings.Builder

	title := lipgloss.NewStyle().
		Foreground(s.Theme.Main).
		Bold(true).
		Render("GAME OVER")

	sb.WriteString(title)
	sb.WriteString("\n\n")

	if m.isNewHS {
		sb.WriteString(lipgloss.NewStyle().
			Foreground(s.Theme.Main).
			Bold(true).
			Render(fmt.Sprintf("NEW HIGH SCORE! #%d", m.rank)))
		sb.WriteString("\n\n")
	}

	labelStyle := lipgloss.NewStyle().Foreground(s.Theme.Sub).Width(10)
	valueStyle := lipgloss.NewStyle().Foreground(s.Theme.FG).Bold(true)

	sb.WriteString(labelStyle.Render("Score") + valueStyle.Render(fmt.Sprintf("%d", m.score)))
	sb.WriteString("\n")
	sb.WriteString(labelStyle.Render("Level") + valueStyle.Render(fmt.Sprintf("%d", m.level)))
	sb.WriteString("\n")
	sb.WriteString(labelStyle.Render("Lines") + valueStyle.Render(fmt.Sprintf("%d", m.lines)))
	sb.WriteString("\n")
	sb.WriteString(labelStyle.Render("Pieces") + valueStyle.Render(fmt.Sprintf("%d", m.pieces)))
	sb.WriteString("\n\n")

	sb.WriteString(lipgloss.NewStyle().Foreground(s.Theme.Sub).Faint(true).Render("r to restart, q/esc to menu"))

	return lipgloss.NewStyle().
		Padding(2, 4).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(s.Theme.Sub).
		Render(sb.String())
}
