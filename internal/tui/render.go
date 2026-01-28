package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/meszmate/briks/internal/game"
	"github.com/meszmate/briks/internal/theme"
)

const (
	// Block characters for pieces (foreground colored, no background)
	blockFull  = "██"
	blockGhost = "░░"
	blockEmpty = "  "
	blockGrid  = "· "
)

// RenderBoard renders the visible portion of the board with the active and ghost pieces.
func RenderBoard(engine *game.Engine, styles Styles, showGhost, showGrid bool, rainbow *theme.RainbowState) string {
	t := styles.Theme

	// Build a visible grid with colors.
	type cell struct {
		color   lipgloss.Color
		isGhost bool
	}
	grid := make([][]cell, game.VisibleRows)
	for r := 0; r < game.VisibleRows; r++ {
		grid[r] = make([]cell, game.BoardWidth)
		for c := 0; c < game.BoardWidth; c++ {
			cellColor := engine.Board.GetVisibleCell(r, c)
			if cellColor != game.Empty {
				grid[r][c] = cell{color: pieceColorToLipgloss(cellColor, t, rainbow)}
			}
		}
	}

	// Draw ghost piece.
	if showGhost && engine.Current != nil {
		ghostCells := engine.GhostCells()
		ghostColor := pieceColorToLipgloss(game.PieceColor(engine.Current.Type), t, rainbow)
		for _, gc := range ghostCells {
			vr := gc.Row - game.BufferRows
			if vr >= 0 && vr < game.VisibleRows && gc.Col >= 0 && gc.Col < game.BoardWidth {
				if grid[vr][gc.Col].color == "" {
					grid[vr][gc.Col] = cell{color: ghostColor, isGhost: true}
				}
			}
		}
	}

	// Draw current piece.
	if engine.Current != nil {
		cells := engine.Current.Cells()
		color := pieceColorToLipgloss(game.PieceColor(engine.Current.Type), t, rainbow)
		for _, c := range cells {
			vr := c.Row - game.BufferRows
			if vr >= 0 && vr < game.VisibleRows && c.Col >= 0 && c.Col < game.BoardWidth {
				grid[vr][c.Col] = cell{color: color, isGhost: false}
			}
		}
	}

	// Render the grid to string.
	var sb strings.Builder
	borderStyle := lipgloss.NewStyle().Foreground(t.Sub)

	// Top border.
	sb.WriteString(borderStyle.Render("┌" + strings.Repeat("──", game.BoardWidth) + "┐"))
	sb.WriteString("\n")

	for r := 0; r < game.VisibleRows; r++ {
		sb.WriteString(borderStyle.Render("│"))
		for c := 0; c < game.BoardWidth; c++ {
			if grid[r][c].color != "" {
				style := lipgloss.NewStyle().Foreground(grid[r][c].color)
				if grid[r][c].isGhost {
					sb.WriteString(style.Render(blockGhost))
				} else {
					sb.WriteString(style.Render(blockFull))
				}
			} else if showGrid {
				sb.WriteString(lipgloss.NewStyle().Foreground(t.SubAlt).Render(blockGrid))
			} else {
				sb.WriteString(blockEmpty)
			}
		}
		sb.WriteString(borderStyle.Render("│"))
		sb.WriteString("\n")
	}

	// Bottom border.
	sb.WriteString(borderStyle.Render("└" + strings.Repeat("──", game.BoardWidth) + "┘"))

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

	style := lipgloss.NewStyle().Foreground(color)
	var sb strings.Builder
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			if grid[r][c] {
				sb.WriteString(style.Render(blockFull))
			} else {
				sb.WriteString(blockEmpty)
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

	titleStyle := lipgloss.NewStyle().Foreground(t.Sub).Bold(true)
	borderStyle := lipgloss.NewStyle().Foreground(t.SubAlt)

	sb.WriteString(titleStyle.Render("HOLD"))
	sb.WriteString("\n")
	sb.WriteString(borderStyle.Render("┌────────┐"))
	sb.WriteString("\n")

	if holdPiece != nil {
		preview := RenderPiecePreview(*holdPiece, t, rainbow)
		lines := strings.Split(preview, "\n")

		// Center the piece in the box (4 cells wide = 8 chars)
		for i := 0; i < 2; i++ {
			sb.WriteString(borderStyle.Render("│"))
			if i < len(lines) {
				content := lines[i]
				if holdUsed {
					content = lipgloss.NewStyle().Faint(true).Render(content)
				}
				width := lipgloss.Width(content)
				pad := (8 - width) / 2
				sb.WriteString(strings.Repeat(" ", pad))
				sb.WriteString(content)
				sb.WriteString(strings.Repeat(" ", 8-width-pad))
			} else {
				sb.WriteString("        ")
			}
			sb.WriteString(borderStyle.Render("│"))
			sb.WriteString("\n")
		}
	} else {
		for i := 0; i < 2; i++ {
			sb.WriteString(borderStyle.Render("│"))
			sb.WriteString("        ")
			sb.WriteString(borderStyle.Render("│"))
			sb.WriteString("\n")
		}
	}

	sb.WriteString(borderStyle.Render("└────────┘"))

	return sb.String()
}

// RenderNextPanel renders the next pieces preview panel.
func RenderNextPanel(pieces []game.PieceType, styles Styles, rainbow *theme.RainbowState) string {
	t := styles.Theme
	var sb strings.Builder

	titleStyle := lipgloss.NewStyle().Foreground(t.Sub).Bold(true)
	borderStyle := lipgloss.NewStyle().Foreground(t.SubAlt)

	sb.WriteString(titleStyle.Render("NEXT"))
	sb.WriteString("\n")
	sb.WriteString(borderStyle.Render("┌────────┐"))
	sb.WriteString("\n")

	for i, pt := range pieces {
		preview := RenderPiecePreview(pt, t, rainbow)
		lines := strings.Split(preview, "\n")
		for j := 0; j < 2; j++ {
			sb.WriteString(borderStyle.Render("│"))
			if j < len(lines) {
				content := lines[j]
				width := lipgloss.Width(content)
				pad := (8 - width) / 2
				sb.WriteString(strings.Repeat(" ", pad))
				sb.WriteString(content)
				sb.WriteString(strings.Repeat(" ", 8-width-pad))
			} else {
				sb.WriteString("        ")
			}
			sb.WriteString(borderStyle.Render("│"))
			sb.WriteString("\n")
		}
		if i < len(pieces)-1 {
			sb.WriteString(borderStyle.Render("│        │"))
			sb.WriteString("\n")
		}
	}

	sb.WriteString(borderStyle.Render("└────────┘"))

	return sb.String()
}

// RenderStatsPanel renders the score/level/lines panel.
func RenderStatsPanel(scorer *game.Scorer, styles Styles) string {
	t := styles.Theme
	var sb strings.Builder

	labelStyle := lipgloss.NewStyle().Foreground(t.Sub)
	valueStyle := lipgloss.NewStyle().Foreground(t.FG).Bold(true)
	highlightStyle := lipgloss.NewStyle().Foreground(t.Main).Bold(true)

	sb.WriteString(labelStyle.Render("SCORE"))
	sb.WriteString("\n")
	sb.WriteString(highlightStyle.Render(fmt.Sprintf("%d", scorer.Score)))
	sb.WriteString("\n\n")

	sb.WriteString(labelStyle.Render("LEVEL"))
	sb.WriteString("\n")
	sb.WriteString(valueStyle.Render(fmt.Sprintf("%d", scorer.Level)))
	sb.WriteString("\n\n")

	sb.WriteString(labelStyle.Render("LINES"))
	sb.WriteString("\n")
	sb.WriteString(valueStyle.Render(fmt.Sprintf("%d", scorer.Lines)))

	if scorer.Combo > 1 {
		sb.WriteString("\n\n")
		sb.WriteString(labelStyle.Render("COMBO"))
		sb.WriteString("\n")
		sb.WriteString(highlightStyle.Render(fmt.Sprintf("%d", scorer.Combo-1)))
	}

	if scorer.BackToBack {
		sb.WriteString("\n\n")
		sb.WriteString(highlightStyle.Render("B2B"))
	}

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
	default:
		return t.FG
	}
}
