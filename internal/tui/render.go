package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/meszmate/briks/internal/game"
	"github.com/meszmate/briks/internal/theme"
)

const (
	cellWidth = 2 // 2 characters per cell for visual squareness
)

// RenderBoard renders the visible portion of the board with the active and ghost pieces.
func RenderBoard(engine *game.Engine, styles Styles, showGhost, showGrid bool, rainbow *theme.RainbowState) string {
	t := styles.Theme

	// Build a visible grid with colors.
	grid := make([][]lipgloss.Color, game.VisibleRows)
	for r := 0; r < game.VisibleRows; r++ {
		grid[r] = make([]lipgloss.Color, game.BoardWidth)
		for c := 0; c < game.BoardWidth; c++ {
			cell := engine.Board.GetVisibleCell(r, c)
			if cell != game.Empty {
				grid[r][c] = pieceColorToLipgloss(cell, t, rainbow)
			}
		}
	}

	// Draw ghost piece.
	if showGhost && engine.Current != nil {
		ghostCells := engine.GhostCells()
		for _, gc := range ghostCells {
			vr := gc.Row - game.BufferRows
			if vr >= 0 && vr < game.VisibleRows && gc.Col >= 0 && gc.Col < game.BoardWidth {
				if grid[vr][gc.Col] == "" {
					grid[vr][gc.Col] = t.Ghost
				}
			}
		}
	}

	// Draw current piece.
	if engine.Current != nil {
		cells := engine.Current.Cells()
		color := pieceColorToLipgloss(game.PieceColor(engine.Current.Type), t, rainbow)
		for _, cell := range cells {
			vr := cell.Row - game.BufferRows
			if vr >= 0 && vr < game.VisibleRows && cell.Col >= 0 && cell.Col < game.BoardWidth {
				grid[vr][cell.Col] = color
			}
		}
	}

	// Render the grid to string.
	var sb strings.Builder

	// Top border.
	sb.WriteString(lipgloss.NewStyle().Foreground(t.Sub).Render("╭" + strings.Repeat("──", game.BoardWidth) + "╮"))
	sb.WriteString("\n")

	for r := 0; r < game.VisibleRows; r++ {
		sb.WriteString(lipgloss.NewStyle().Foreground(t.Sub).Render("│"))
		for c := 0; c < game.BoardWidth; c++ {
			if grid[r][c] != "" {
				sb.WriteString(lipgloss.NewStyle().Background(grid[r][c]).Render("  "))
			} else if showGrid {
				sb.WriteString(lipgloss.NewStyle().Background(t.BG).Foreground(t.Grid).Render(" ·"))
			} else {
				sb.WriteString(lipgloss.NewStyle().Background(t.BG).Render("  "))
			}
		}
		sb.WriteString(lipgloss.NewStyle().Foreground(t.Sub).Render("│"))
		sb.WriteString("\n")
	}

	// Bottom border.
	sb.WriteString(lipgloss.NewStyle().Foreground(t.Sub).Render("╰" + strings.Repeat("──", game.BoardWidth) + "╯"))

	return sb.String()
}

// RenderPiecePreview renders a small preview of a piece type.
func RenderPiecePreview(pt game.PieceType, t theme.Theme, rainbow *theme.RainbowState) string {
	offsets := game.PieceRotations[pt][game.Rot0]
	color := pieceColorToLipgloss(game.PieceColor(pt), t, rainbow)

	// Find bounding box.
	minR, maxR, minC, maxC := 10, -1, 10, -1
	for _, off := range offsets {
		if off.Row < minR {
			minR = off.Row
		}
		if off.Row > maxR {
			maxR = off.Row
		}
		if off.Col < minC {
			minC = off.Col
		}
		if off.Col > maxC {
			maxC = off.Col
		}
	}

	height := maxR - minR + 1
	width := maxC - minC + 1
	grid := make([][]bool, height)
	for i := range grid {
		grid[i] = make([]bool, width)
	}
	for _, off := range offsets {
		grid[off.Row-minR][off.Col-minC] = true
	}

	var sb strings.Builder
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			if grid[r][c] {
				sb.WriteString(lipgloss.NewStyle().Background(color).Render("  "))
			} else {
				sb.WriteString(lipgloss.NewStyle().Background(t.BG).Render("  "))
			}
		}
		if r < height-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// RenderHoldPanel renders the hold piece panel.
func RenderHoldPanel(holdPiece *game.PieceType, holdUsed bool, styles Styles, rainbow *theme.RainbowState) string {
	t := styles.Theme
	var sb strings.Builder

	sb.WriteString(styles.PanelTitle.Render("HOLD"))
	sb.WriteString("\n")
	sb.WriteString(lipgloss.NewStyle().Foreground(t.Sub).Render("╭────────╮"))
	sb.WriteString("\n")

	if holdPiece != nil {
		preview := RenderPiecePreview(*holdPiece, t, rainbow)
		lines := strings.Split(preview, "\n")
		for i := 0; i < 2; i++ {
			sb.WriteString(lipgloss.NewStyle().Foreground(t.Sub).Render("│"))
			if i < len(lines) {
				content := lines[i]
				if holdUsed {
					content = lipgloss.NewStyle().Faint(true).Render(content)
				}
				// Pad to 8 chars (4 cells * 2 chars).
				padded := content + lipgloss.NewStyle().Background(t.BG).Render(strings.Repeat(" ", 8-lipgloss.Width(content)))
				sb.WriteString(padded)
			} else {
				sb.WriteString(lipgloss.NewStyle().Background(t.BG).Render("        "))
			}
			sb.WriteString(lipgloss.NewStyle().Foreground(t.Sub).Render("│"))
			sb.WriteString("\n")
		}
	} else {
		for i := 0; i < 2; i++ {
			sb.WriteString(lipgloss.NewStyle().Foreground(t.Sub).Render("│"))
			sb.WriteString(lipgloss.NewStyle().Background(t.BG).Render("        "))
			sb.WriteString(lipgloss.NewStyle().Foreground(t.Sub).Render("│"))
			sb.WriteString("\n")
		}
	}

	sb.WriteString(lipgloss.NewStyle().Foreground(t.Sub).Render("╰────────╯"))

	return sb.String()
}

// RenderNextPanel renders the next pieces preview panel.
func RenderNextPanel(pieces []game.PieceType, styles Styles, rainbow *theme.RainbowState) string {
	t := styles.Theme
	var sb strings.Builder

	sb.WriteString(styles.PanelTitle.Render("NEXT"))
	sb.WriteString("\n")
	sb.WriteString(lipgloss.NewStyle().Foreground(t.Sub).Render("╭────────╮"))
	sb.WriteString("\n")

	for i, pt := range pieces {
		preview := RenderPiecePreview(pt, t, rainbow)
		lines := strings.Split(preview, "\n")
		for j := 0; j < 2; j++ {
			sb.WriteString(lipgloss.NewStyle().Foreground(t.Sub).Render("│"))
			if j < len(lines) {
				content := lines[j]
				padded := content + lipgloss.NewStyle().Background(t.BG).Render(strings.Repeat(" ", 8-lipgloss.Width(content)))
				sb.WriteString(padded)
			} else {
				sb.WriteString(lipgloss.NewStyle().Background(t.BG).Render("        "))
			}
			sb.WriteString(lipgloss.NewStyle().Foreground(t.Sub).Render("│"))
			sb.WriteString("\n")
		}
		if i < len(pieces)-1 {
			sb.WriteString(lipgloss.NewStyle().Foreground(t.Sub).Render("│"))
			sb.WriteString(lipgloss.NewStyle().Background(t.BG).Render("        "))
			sb.WriteString(lipgloss.NewStyle().Foreground(t.Sub).Render("│"))
			sb.WriteString("\n")
		}
	}

	sb.WriteString(lipgloss.NewStyle().Foreground(t.Sub).Render("╰────────╯"))

	return sb.String()
}

// RenderStatsPanel renders the score/level/lines panel.
func RenderStatsPanel(scorer *game.Scorer, styles Styles) string {
	var sb strings.Builder

	sb.WriteString(styles.PanelTitle.Render("SCORE"))
	sb.WriteString("\n")
	sb.WriteString(styles.Highlight.Render(fmt.Sprintf("%d", scorer.Score)))
	sb.WriteString("\n\n")

	sb.WriteString(styles.PanelTitle.Render("LEVEL"))
	sb.WriteString("\n")
	sb.WriteString(styles.Text.Render(fmt.Sprintf("%d", scorer.Level)))
	sb.WriteString("\n\n")

	sb.WriteString(styles.PanelTitle.Render("LINES"))
	sb.WriteString("\n")
	sb.WriteString(styles.Text.Render(fmt.Sprintf("%d", scorer.Lines)))
	sb.WriteString("\n\n")

	sb.WriteString(styles.PanelTitle.Render("COMBO"))
	sb.WriteString("\n")
	combo := scorer.Combo
	if combo > 0 {
		combo-- // display active combo count
	}
	sb.WriteString(styles.Text.Render(fmt.Sprintf("%d", combo)))

	return sb.String()
}

func pieceColorToLipgloss(c game.CellColor, t theme.Theme, rainbow *theme.RainbowState) lipgloss.Color {
	if rainbow != nil && t.Name == "rainbow" {
		switch c {
		case game.ColorI:
			return rainbow.PieceColor(0)
		case game.ColorO:
			return rainbow.PieceColor(1)
		case game.ColorT:
			return rainbow.PieceColor(2)
		case game.ColorS:
			return rainbow.PieceColor(3)
		case game.ColorZ:
			return rainbow.PieceColor(4)
		case game.ColorJ:
			return rainbow.PieceColor(5)
		case game.ColorL:
			return rainbow.PieceColor(6)
		}
	}

	switch c {
	case game.ColorI:
		return t.PieceI
	case game.ColorO:
		return t.PieceO
	case game.ColorT:
		return t.PieceT
	case game.ColorS:
		return t.PieceS
	case game.ColorZ:
		return t.PieceZ
	case game.ColorJ:
		return t.PieceJ
	case game.ColorL:
		return t.PieceL
	case game.ColorGhost:
		return t.Ghost
	default:
		return t.BG
	}
}
