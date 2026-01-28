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
	t := s.Theme
	var sb strings.Builder

	// Logo
	logo := `
 ██████╗ ██████╗ ██╗██╗  ██╗███████╗
 ██╔══██╗██╔══██╗██║██║ ██╔╝██╔════╝
 ██████╔╝██████╔╝██║█████╔╝ ███████╗
 ██╔══██╗██╔══██╗██║██╔═██╗ ╚════██║
 ██████╔╝██║  ██║██║██║  ██╗███████║
 ╚═════╝ ╚═╝  ╚═╝╚═╝╚═╝  ╚═╝╚══════╝`

	sb.WriteString(lipgloss.NewStyle().Foreground(t.Main).Render(logo))
	sb.WriteString("\n\n")

	// Menu items
	for i, item := range menuItems {
		if i == m.cursor {
			sb.WriteString(lipgloss.NewStyle().
				Foreground(t.Main).
				Bold(true).
				Render(" > " + item))
		} else {
			sb.WriteString(lipgloss.NewStyle().
				Foreground(t.Sub).
				Render("   " + item))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(lipgloss.NewStyle().
		Foreground(t.SubAlt).
		Render("   j/k navigate  enter select  q quit"))

	return sb.String()
}
