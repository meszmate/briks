package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/meszmate/briks/internal/config"
	"github.com/meszmate/briks/internal/game"
	"github.com/meszmate/briks/internal/theme"
)

// GameModel handles the gameplay screen.
type GameModel struct {
	engine   *game.Engine
	keys     *config.KeyBindings
	rainbow  *theme.RainbowState
	paused   bool
	gameOver bool
}

// NewGameModel creates a new gameplay model.
func NewGameModel(cfg *config.Config, keys *config.KeyBindings, rainbow *theme.RainbowState) GameModel {
	return GameModel{
		engine:  game.NewEngine(cfg.StartLevel, cfg.PreviewCount),
		keys:    keys,
		rainbow: rainbow,
	}
}

// Init returns the initial commands for the game.
func (g GameModel) Init() tea.Cmd {
	return tea.Batch(
		g.gravityTick(),
		g.lockTick(),
		g.rainbowTick(),
	)
}

func (g GameModel) gravityTick() tea.Cmd {
	interval := g.engine.Scorer.GravityInterval()
	d := time.Duration(float64(time.Second) * interval)
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return TickMsg{Time: t}
	})
}

func (g GameModel) lockTick() tea.Cmd {
	return tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
		return LockTickMsg{Time: t}
	})
}

func (g GameModel) rainbowTick() tea.Cmd {
	return tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
		return RainbowTickMsg{Time: t}
	})
}

func (g GameModel) resumeTick() tea.Cmd {
	return tea.Batch(
		g.gravityTick(),
		g.lockTick(),
		g.rainbowTick(),
	)
}

// Update processes messages for the game screen.
func (g GameModel) Update(msg tea.Msg, keys *config.KeyBindings) (GameModel, tea.Cmd) {
	if g.paused || g.gameOver {
		return g, nil
	}

	switch msg := msg.(type) {
	case TickMsg:
		if g.engine.State == game.StatePlaying {
			g.engine.Tick()
			if g.engine.State == game.StateGameOver {
				g.gameOver = true
				return g, nil
			}
		}
		return g, g.gravityTick()

	case LockTickMsg:
		if g.engine.State == game.StatePlaying {
			g.engine.CheckLock()
			if g.engine.State == game.StateGameOver {
				g.gameOver = true
				return g, nil
			}
		}
		return g, g.lockTick()

	case RainbowTickMsg:
		if g.rainbow != nil {
			g.rainbow.Tick(0.05)
		}
		return g, g.rainbowTick()

	case tea.KeyMsg:
		return g.handleKey(msg, keys)
	}

	return g, nil
}

func (g GameModel) handleKey(msg tea.KeyMsg, keys *config.KeyBindings) (GameModel, tea.Cmd) {
	key := msg.String()

	action, ok := keys.MatchAction(key)
	if !ok {
		return g, nil
	}

	switch action {
	case config.ActionMoveLeft:
		g.engine.MoveLeft()
	case config.ActionMoveRight:
		g.engine.MoveRight()
	case config.ActionSoftDrop:
		g.engine.SoftDrop()
	case config.ActionHardDrop:
		g.engine.HardDrop()
		if g.engine.State == game.StateGameOver {
			g.gameOver = true
			return g, nil
		}
	case config.ActionRotateCW:
		g.engine.RotateCW()
	case config.ActionRotateCCW:
		g.engine.RotateCCW()
	case config.ActionHold:
		g.engine.Hold()
		if g.engine.State == game.StateGameOver {
			g.gameOver = true
			return g, nil
		}
	case config.ActionPause:
		g.paused = true
		return g, nil
	}

	return g, nil
}

// View renders the gameplay screen.
func (g GameModel) View(s Styles, cfg *config.Config, rainbow *theme.RainbowState) string {
	t := s.Theme

	board := RenderBoard(g.engine, s, cfg.GhostPiece, cfg.ShowGrid, rainbow)
	hold := RenderHoldPanel(g.engine.HoldPiece, g.engine.HoldUsed, s, rainbow)
	next := RenderNextPanel(g.engine.NextPieces(), s, rainbow)
	stats := RenderStatsPanel(g.engine.Scorer, s)

	// Build left panel
	var leftSb strings.Builder
	leftSb.WriteString(hold)
	leftSb.WriteString("\n\n")
	leftSb.WriteString(stats)

	leftPanel := lipgloss.NewStyle().
		Width(12).
		Render(leftSb.String())

	// Build right panel
	rightPanel := lipgloss.NewStyle().
		Width(12).
		Render(next)

	// Build help text
	helpStyle := lipgloss.NewStyle().Foreground(t.SubAlt)
	help := helpStyle.Render(fmt.Sprintf("h/l move  j drop  k rotate  c hold  space hard drop  p pause"))

	// Combine panels
	gameRow := lipgloss.JoinHorizontal(lipgloss.Top,
		leftPanel,
		"  ",
		board,
		"  ",
		rightPanel,
	)

	return lipgloss.JoinVertical(lipgloss.Center, gameRow, "", help)
}
