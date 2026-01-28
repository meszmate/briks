package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/meszmate/briks/internal/config"
)

// KeyBindsModel handles the key binding configuration screen.
type KeyBindsModel struct {
	cursor    int
	listening bool
	keys      *config.KeyBindings
}

// NewKeyBindsModel creates a new key bindings model.
func NewKeyBindsModel(keys *config.KeyBindings, s Styles) KeyBindsModel {
	return KeyBindsModel{keys: keys}
}

// Update handles navigation input.
func (m KeyBindsModel) Update(msg tea.KeyMsg, keys *config.KeyBindings) KeyBindsModel {
	m.keys = keys
	switch msg.String() {
	case "j", "down":
		m.cursor = (m.cursor + 1) % len(config.AllActions)
	case "k", "up":
		m.cursor = (m.cursor - 1 + len(config.AllActions)) % len(config.AllActions)
	case "enter":
		m.listening = true
	}
	return m
}

// HandleKey processes a key press while listening for a new binding.
func (m KeyBindsModel) HandleKey(msg tea.KeyMsg, keys *config.KeyBindings) KeyBindsModel {
	if !m.listening {
		return m
	}

	key := msg.String()
	if key == "esc" {
		m.listening = false
		return m
	}

	action := config.AllActions[m.cursor]
	keys.SetBinding(action, []string{key})
	m.keys = keys
	m.listening = false
	return m
}

// View renders the key bindings screen.
func (m KeyBindsModel) View(s Styles) string {
	var sb strings.Builder

	title := lipgloss.NewStyle().
		Foreground(s.Theme.Main).
		Bold(true).
		Render("KEY BINDINGS")

	sb.WriteString(title)
	sb.WriteString("\n\n")

	for i, action := range config.AllActions {
		label := config.ActionLabel(action)
		labelStyle := lipgloss.NewStyle().Width(16)

		var keyStr string
		if m.listening && i == m.cursor {
			keyStr = lipgloss.NewStyle().
				Foreground(s.Theme.Main).
				Bold(true).
				Blink(true).
				Render("press a key...")
		} else {
			binds := m.keys.Bindings[action]
			keyStr = lipgloss.NewStyle().
				Foreground(s.Theme.FG).
				Render(strings.Join(binds, ", "))
		}

		if i == m.cursor {
			sb.WriteString(lipgloss.NewStyle().
				Foreground(s.Theme.Main).
				Bold(true).
				Render("▸ "))
			sb.WriteString(labelStyle.Foreground(s.Theme.Main).Render(label))
			sb.WriteString(keyStr)
		} else {
			sb.WriteString(lipgloss.NewStyle().
				Foreground(s.Theme.Sub).
				Render("  "))
			sb.WriteString(labelStyle.Foreground(s.Theme.Sub).Render(label))
			sb.WriteString(keyStr)
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(lipgloss.NewStyle().
		Foreground(s.Theme.Sub).
		Faint(true).
		Render("enter to rebind, esc to cancel, q to save & back"))

	return lipgloss.NewStyle().
		Padding(2, 4).
		Render(sb.String())
}
