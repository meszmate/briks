package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/meszmate/briks/internal/config"
	"github.com/meszmate/briks/internal/tui"
)

func main() {
	cfg := config.Load()
	keys := config.LoadKeyBindings()
	hs := config.LoadHighScores()

	app := tui.NewApp(cfg, keys, hs)

	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
