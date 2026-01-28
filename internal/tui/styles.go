package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/meszmate/briks/internal/theme"
)

// Styles holds all Lipgloss styles derived from the active theme.
type Styles struct {
	Theme theme.Theme

	// Base styles.
	Background lipgloss.Style
	Title      lipgloss.Style
	Subtitle   lipgloss.Style
	Text       lipgloss.Style
	Dim        lipgloss.Style
	Highlight  lipgloss.Style
	Selected   lipgloss.Style

	// Board styles.
	BoardBorder lipgloss.Style
	Cell        lipgloss.Style
	EmptyCell   lipgloss.Style
	GridCell    lipgloss.Style
	GhostCell   lipgloss.Style

	// Panel styles.
	Panel      lipgloss.Style
	PanelTitle lipgloss.Style

	// Menu styles.
	MenuItem         lipgloss.Style
	MenuItemSelected lipgloss.Style
}

// NewStyles creates styles from a theme.
func NewStyles(t theme.Theme) Styles {
	s := Styles{Theme: t}

	s.Background = lipgloss.NewStyle().
		Background(t.BG).
		Foreground(t.FG)

	s.Title = lipgloss.NewStyle().
		Foreground(t.Main).
		Bold(true)

	s.Subtitle = lipgloss.NewStyle().
		Foreground(t.Sub)

	s.Text = lipgloss.NewStyle().
		Foreground(t.FG)

	s.Dim = lipgloss.NewStyle().
		Foreground(t.Sub)

	s.Highlight = lipgloss.NewStyle().
		Foreground(t.Main).
		Bold(true)

	s.Selected = lipgloss.NewStyle().
		Foreground(t.Main)

	s.BoardBorder = lipgloss.NewStyle().
		Foreground(t.SubAlt)

	s.Cell = lipgloss.NewStyle().
		Background(t.BG)

	s.EmptyCell = lipgloss.NewStyle().
		Background(t.BG).
		Foreground(t.BG)

	s.GridCell = lipgloss.NewStyle().
		Background(t.BG).
		Foreground(t.Grid)

	s.GhostCell = lipgloss.NewStyle().
		Background(t.Ghost)

	s.Panel = lipgloss.NewStyle().
		Foreground(t.FG).
		PaddingLeft(1).
		PaddingRight(1)

	s.PanelTitle = lipgloss.NewStyle().
		Foreground(t.Main).
		Bold(true)

	s.MenuItem = lipgloss.NewStyle().
		Foreground(t.Sub).
		PaddingLeft(2)

	s.MenuItemSelected = lipgloss.NewStyle().
		Foreground(t.Main).
		Bold(true).
		PaddingLeft(2)

	return s
}

// CellStyle returns the style for a cell with the given color.
func (s *Styles) CellStyle(c lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Background(c).
		Foreground(c)
}
