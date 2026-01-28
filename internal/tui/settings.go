package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/meszmate/briks/internal/config"
	"github.com/meszmate/briks/internal/theme"
)

type settingItem struct {
	label string
	key   string
}

var settingItems = []settingItem{
	{"Theme", "theme"},
	{"Start Level", "start_level"},
	{"Ghost Piece", "ghost_piece"},
	{"Show Grid", "show_grid"},
	{"Preview Count", "preview_count"},
	{"DAS (ms)", "das"},
	{"ARR (ms)", "arr"},
}

// SettingsModel handles the settings screen.
type SettingsModel struct {
	cursor int
	cfg    *config.Config
	styles Styles
}

// NewSettingsModel creates a new settings model.
func NewSettingsModel(cfg *config.Config, s Styles) SettingsModel {
	return SettingsModel{cfg: cfg, styles: s}
}

// Update handles settings input.
func (m SettingsModel) Update(msg tea.KeyMsg, cfg *config.Config) SettingsModel {
	switch msg.String() {
	case "j", "down":
		m.cursor = (m.cursor + 1) % len(settingItems)
	case "k", "up":
		m.cursor = (m.cursor - 1 + len(settingItems)) % len(settingItems)
	case "l", "right":
		m.cycleValue(cfg, 1)
	case "h", "left":
		m.cycleValue(cfg, -1)
	}
	m.cfg = cfg
	return m
}

func (m *SettingsModel) cycleValue(cfg *config.Config, dir int) {
	key := settingItems[m.cursor].key
	switch key {
	case "theme":
		names := theme.ThemeNames
		idx := 0
		for i, n := range names {
			if n == cfg.Theme {
				idx = i
				break
			}
		}
		idx = (idx + dir + len(names)) % len(names)
		cfg.Theme = names[idx]
	case "start_level":
		cfg.StartLevel += dir
		if cfg.StartLevel < 1 {
			cfg.StartLevel = 20
		}
		if cfg.StartLevel > 20 {
			cfg.StartLevel = 1
		}
	case "ghost_piece":
		cfg.GhostPiece = !cfg.GhostPiece
	case "show_grid":
		cfg.ShowGrid = !cfg.ShowGrid
	case "preview_count":
		cfg.PreviewCount += dir
		if cfg.PreviewCount < 1 {
			cfg.PreviewCount = 5
		}
		if cfg.PreviewCount > 5 {
			cfg.PreviewCount = 1
		}
	case "das":
		cfg.DAS += dir * 10
		if cfg.DAS < 50 {
			cfg.DAS = 50
		}
		if cfg.DAS > 500 {
			cfg.DAS = 500
		}
	case "arr":
		cfg.ARR += dir * 5
		if cfg.ARR < 0 {
			cfg.ARR = 0
		}
		if cfg.ARR > 200 {
			cfg.ARR = 200
		}
	}
}

func getValue(cfg *config.Config, key string) string {
	switch key {
	case "theme":
		return cfg.Theme
	case "start_level":
		return fmt.Sprintf("%d", cfg.StartLevel)
	case "ghost_piece":
		if cfg.GhostPiece {
			return "on"
		}
		return "off"
	case "show_grid":
		if cfg.ShowGrid {
			return "on"
		}
		return "off"
	case "preview_count":
		return fmt.Sprintf("%d", cfg.PreviewCount)
	case "das":
		return fmt.Sprintf("%d", cfg.DAS)
	case "arr":
		return fmt.Sprintf("%d", cfg.ARR)
	default:
		return ""
	}
}

// View renders the settings screen.
func (m SettingsModel) View(s Styles) string {
	t := s.Theme
	var sb strings.Builder

	title := lipgloss.NewStyle().
		Foreground(t.Main).
		Bold(true).
		Render("SETTINGS")

	sb.WriteString(title)
	sb.WriteString("\n\n")

	for i, item := range settingItems {
		labelStyle := lipgloss.NewStyle().Width(14)
		value := getValue(m.cfg, item.key)

		if i == m.cursor {
			sb.WriteString(lipgloss.NewStyle().Foreground(t.Main).Render(" > "))
			sb.WriteString(labelStyle.Foreground(t.FG).Render(item.label))
			sb.WriteString(lipgloss.NewStyle().Foreground(t.Main).Render("< " + value + " >"))
		} else {
			sb.WriteString(lipgloss.NewStyle().Foreground(t.Sub).Render("   "))
			sb.WriteString(labelStyle.Foreground(t.Sub).Render(item.label))
			sb.WriteString(lipgloss.NewStyle().Foreground(t.Sub).Render("  " + value))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(lipgloss.NewStyle().
		Foreground(t.SubAlt).
		Render("   h/l change  j/k navigate  q save"))

	return sb.String()
}
