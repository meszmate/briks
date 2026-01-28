package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/meszmate/briks/internal/config"
	"github.com/meszmate/briks/internal/theme"
)

// Screen represents the active screen.
type Screen int

const (
	ScreenMenu Screen = iota
	ScreenGame
	ScreenPause
	ScreenGameOver
	ScreenSettings
	ScreenHighScores
	ScreenKeyBinds
)

const (
	MinWidth  = 60
	MinHeight = 28
)

// App is the root Bubble Tea model.
type App struct {
	screen Screen
	width  int
	height int

	cfg        *config.Config
	keys       *config.KeyBindings
	highScores *config.HighScores
	styles     Styles
	rainbow    *theme.RainbowState

	menu     MenuModel
	game     GameModel
	pause    PauseModel
	gameOver GameOverModel
	settings SettingsModel
	scores   HighScoresModel
	keyBinds KeyBindsModel
}

// NewApp creates the root application model.
func NewApp(cfg *config.Config, keys *config.KeyBindings, hs *config.HighScores) App {
	t := theme.GetTheme(cfg.Theme)
	s := NewStyles(t)
	rb := theme.NewRainbowState()

	app := App{
		screen:     ScreenMenu,
		cfg:        cfg,
		keys:       keys,
		highScores: hs,
		styles:     s,
		rainbow:    rb,
	}

	app.menu = NewMenuModel(s)
	app.settings = NewSettingsModel(cfg, s)
	app.scores = NewHighScoresModel(hs, s)
	app.keyBinds = NewKeyBindsModel(keys, s)

	return app
}

func (a App) Init() tea.Cmd {
	return nil
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		return a, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return a, tea.Quit
		}
	}

	switch a.screen {
	case ScreenMenu:
		return a.updateMenu(msg)
	case ScreenGame:
		return a.updateGame(msg)
	case ScreenPause:
		return a.updatePause(msg)
	case ScreenGameOver:
		return a.updateGameOver(msg)
	case ScreenSettings:
		return a.updateSettings(msg)
	case ScreenHighScores:
		return a.updateHighScores(msg)
	case ScreenKeyBinds:
		return a.updateKeyBinds(msg)
	}

	return a, nil
}

func (a App) View() string {
	if a.width < MinWidth || a.height < MinHeight {
		msg := lipgloss.NewStyle().
			Foreground(a.styles.Theme.Main).
			Bold(true).
			Render("Terminal too small\nMinimum: 60x28")
		return lipgloss.Place(a.width, a.height, lipgloss.Center, lipgloss.Center, msg)
	}

	var content string

	switch a.screen {
	case ScreenMenu:
		content = a.menu.View(a.styles)
	case ScreenGame:
		content = a.game.View(a.styles, a.cfg, a.rainbow)
	case ScreenPause:
		content = a.pause.View(a.styles)
	case ScreenGameOver:
		content = a.gameOver.View(a.styles)
	case ScreenSettings:
		content = a.settings.View(a.styles)
	case ScreenHighScores:
		content = a.scores.View(a.styles)
	case ScreenKeyBinds:
		content = a.keyBinds.View(a.styles)
	}

	return lipgloss.Place(a.width, a.height, lipgloss.Center, lipgloss.Center, content)
}

// refreshedStyles returns new styles from the current config theme.
func refreshedStyles(cfg *config.Config) Styles {
	t := theme.GetTheme(cfg.Theme)
	return NewStyles(t)
}

// Screen transition helpers.

func (a App) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return a, tea.Quit
		case "j", "down":
			a.menu.Next()
		case "k", "up":
			a.menu.Prev()
		case "enter", "l":
			switch a.menu.Selected() {
			case 0: // Play
				a.game = NewGameModel(a.cfg, a.keys, a.rainbow)
				a.screen = ScreenGame
				return a, a.game.Init()
			case 1: // Settings
				a.settings = NewSettingsModel(a.cfg, a.styles)
				a.screen = ScreenSettings
			case 2: // High Scores
				a.scores = NewHighScoresModel(a.highScores, a.styles)
				a.screen = ScreenHighScores
			case 3: // Key Bindings
				a.keyBinds = NewKeyBindsModel(a.keys, a.styles)
				a.screen = ScreenKeyBinds
			case 4: // Quit
				return a, tea.Quit
			}
		}
	}
	return a, nil
}

func (a App) updateGame(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.game, cmd = a.game.Update(msg, a.keys)

	if a.game.paused {
		a.pause = NewPauseModel()
		a.screen = ScreenPause
		return a, nil
	}

	if a.game.gameOver {
		a.gameOver = NewGameOverModel(a.game.engine, a.highScores)
		a.screen = ScreenGameOver
		return a, nil
	}

	return a, cmd
}

func (a App) updatePause(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "p", "esc":
			a.screen = ScreenGame
			a.game.paused = false
			return a, a.game.resumeTick()
		case "q":
			a.screen = ScreenMenu
			a.menu = NewMenuModel(a.styles)
		case "r":
			a.game = NewGameModel(a.cfg, a.keys, a.rainbow)
			a.screen = ScreenGame
			return a, a.game.Init()
		}
	}
	return a, nil
}

func (a App) updateGameOver(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			a.game = NewGameModel(a.cfg, a.keys, a.rainbow)
			a.screen = ScreenGame
			return a, a.game.Init()
		case "q", "esc", "enter":
			a.screen = ScreenMenu
			a.menu = NewMenuModel(a.styles)
		}
	}
	return a, nil
}

func (a App) updateSettings(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			a.cfg.Save()
			a.styles = refreshedStyles(a.cfg)
			a.screen = ScreenMenu
			a.menu = NewMenuModel(a.styles)
			return a, nil
		default:
			a.settings = a.settings.Update(msg, a.cfg)
			a.styles = refreshedStyles(a.cfg)
			a.settings.styles = a.styles
		}
	}
	return a, nil
}

func (a App) updateHighScores(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "enter":
			a.screen = ScreenMenu
			a.menu = NewMenuModel(a.styles)
		}
	}
	return a, nil
}

func (a App) updateKeyBinds(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if a.keyBinds.listening {
			a.keyBinds = a.keyBinds.HandleKey(msg, a.keys)
			return a, nil
		}
		switch msg.String() {
		case "q":
			a.keys.Save()
			a.screen = ScreenMenu
			a.menu = NewMenuModel(a.styles)
		case "esc":
			if a.keyBinds.listening {
				a.keyBinds.listening = false
			} else {
				a.keys.Save()
				a.screen = ScreenMenu
				a.menu = NewMenuModel(a.styles)
			}
		default:
			a.keyBinds = a.keyBinds.Update(msg, a.keys)
		}
	}
	return a, nil
}
