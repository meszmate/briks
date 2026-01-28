package theme

import "github.com/charmbracelet/lipgloss"

// Theme defines the color palette for the TUI.
type Theme struct {
	Name   string
	BG     lipgloss.Color // background
	FG     lipgloss.Color // primary foreground
	Main   lipgloss.Color // accent / highlight
	Sub    lipgloss.Color // secondary text
	SubAlt lipgloss.Color // dim / border

	// Piece colors.
	PieceI lipgloss.Color
	PieceO lipgloss.Color
	PieceT lipgloss.Color
	PieceS lipgloss.Color
	PieceZ lipgloss.Color
	PieceJ lipgloss.Color
	PieceL lipgloss.Color

	Ghost lipgloss.Color
	Grid  lipgloss.Color
}
