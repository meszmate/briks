package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var menuItems = []string{
	"Play",
	"Settings",
	"High Scores",
	"Key Bindings",
	"Quit",
}

// MenuModel represents the main menu.
type MenuModel struct {
	cursor int
}

// NewMenuModel creates a new menu.
func NewMenuModel(s Styles) MenuModel {
	return MenuModel{}
}

// Selected returns the currently selected menu index.
func (m *MenuModel) Selected() int {
	return m.cursor
}

// Next moves the cursor down.
func (m *MenuModel) Next() {
	m.cursor = (m.cursor + 1) % len(menuItems)
}

// Prev moves the cursor up.
func (m *MenuModel) Prev() {
	m.cursor = (m.cursor - 1 + len(menuItems)) % len(menuItems)
}

// View renders the menu.
func (m MenuModel) View(s Styles) string {
	var sb strings.Builder

	title := lipgloss.NewStyle().
		Foreground(s.Theme.Main).
		Bold(true).
		MarginBottom(1).
		Render("BRIKS")

	subtitle := lipgloss.NewStyle().
		Foreground(s.Theme.Sub).
		MarginBottom(2).
		Render("Terminal Tetris")

	sb.WriteString(title)
	sb.WriteString("\n")
	sb.WriteString(subtitle)
	sb.WriteString("\n\n")

	for i, item := range menuItems {
		if i == m.cursor {
			sb.WriteString(lipgloss.NewStyle().
				Foreground(s.Theme.Main).
				Bold(true).
				Render("▸ " + item))
		} else {
			sb.WriteString(lipgloss.NewStyle().
				Foreground(s.Theme.Sub).
				Render("  " + item))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(lipgloss.NewStyle().
		Foreground(s.Theme.Sub).
		Faint(true).
		Render("j/k or ↑/↓ to navigate, enter to select"))

	return lipgloss.NewStyle().
		Padding(2, 4).
		Render(sb.String())
}
