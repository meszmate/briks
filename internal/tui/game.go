package tui

import (
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

	// DAS/ARR state.
	dasAction   config.Action
	dasStarted  bool
	dasActive   bool
	dasTimer    time.Time
	dasInterval int
	arrInterval int
}

// NewGameModel creates a new gameplay model.
func NewGameModel(cfg *config.Config, keys *config.KeyBindings, rainbow *theme.RainbowState) GameModel {
	return GameModel{
		engine:      game.NewEngine(cfg.StartLevel, cfg.PreviewCount),
		keys:        keys,
		rainbow:     rainbow,
		dasInterval: cfg.DAS,
		arrInterval: cfg.ARR,
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

func (g GameModel) dasTick() tea.Cmd {
	d := time.Duration(g.arrInterval) * time.Millisecond
	if d < time.Millisecond {
		d = time.Millisecond
	}
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return DASTickMsg{Time: t}
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

	case DASTickMsg:
		if g.engine.State == game.StatePlaying && g.dasStarted {
			if !g.dasActive {
				// DAS delay just elapsed — activate ARR.
				g.dasActive = true
			}
			switch g.dasAction {
			case config.ActionMoveLeft:
				g.engine.MoveLeft()
			case config.ActionMoveRight:
				g.engine.MoveRight()
			case config.ActionSoftDrop:
				g.engine.SoftDrop()
			}
			return g, g.dasTick()
		}
		return g, nil

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
		return g.startDAS(config.ActionMoveLeft)
	case config.ActionMoveRight:
		g.engine.MoveRight()
		return g.startDAS(config.ActionMoveRight)
	case config.ActionSoftDrop:
		g.engine.SoftDrop()
		return g.startDAS(config.ActionSoftDrop)
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

func (g GameModel) startDAS(action config.Action) (GameModel, tea.Cmd) {
	g.dasAction = action
	g.dasStarted = true
	g.dasActive = false
	g.dasTimer = time.Now()

	// Start DAS delay then ARR.
	d := time.Duration(g.dasInterval) * time.Millisecond
	return g, tea.Tick(d, func(t time.Time) tea.Msg {
		return DASTickMsg{Time: t}
	})
}

// View renders the gameplay screen.
func (g GameModel) View(s Styles, cfg *config.Config, rainbow *theme.RainbowState) string {
	board := RenderBoard(g.engine, s, cfg.GhostPiece, cfg.ShowGrid, rainbow)
	hold := RenderHoldPanel(g.engine.HoldPiece, g.engine.HoldUsed, s, rainbow)
	next := RenderNextPanel(g.engine.NextPieces(), s, rainbow)
	stats := RenderStatsPanel(g.engine.Scorer, s)

	leftPanel := lipgloss.NewStyle().
		Width(12).
		MarginRight(1).
		Render(lipgloss.JoinVertical(lipgloss.Left, hold, "", stats))

	rightPanel := lipgloss.NewStyle().
		Width(12).
		MarginLeft(1).
		Render(next)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, board, rightPanel)
}
