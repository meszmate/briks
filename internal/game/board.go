package game

const (
	BoardWidth  = 10
	BoardHeight = 24 // 20 visible + 4 buffer rows at top
	VisibleRows = 20
	BufferRows  = 4
)

// Board represents the Tetris playing field.
type Board struct {
	Cells [BoardHeight][BoardWidth]CellColor
}

// NewBoard creates an empty board.
func NewBoard() *Board {
	return &Board{}
}

// InBounds checks if a position is within the board.
func (b *Board) InBounds(pos Position) bool {
	return pos.Row >= 0 && pos.Row < BoardHeight && pos.Col >= 0 && pos.Col < BoardWidth
}

// IsEmpty checks if a cell is empty (and in bounds).
func (b *Board) IsEmpty(pos Position) bool {
	if !b.InBounds(pos) {
		return false
	}
	return b.Cells[pos.Row][pos.Col] == Empty
}

// IsOccupied checks if a cell is occupied (or out of bounds).
func (b *Board) IsOccupied(pos Position) bool {
	return !b.IsEmpty(pos)
}

// ValidPosition checks if a piece can exist at its current position.
func (b *Board) ValidPosition(p *Piece) bool {
	for _, cell := range p.Cells() {
		if !b.InBounds(cell) || b.Cells[cell.Row][cell.Col] != Empty {
			return false
		}
	}
	return true
}

// PlacePiece locks a piece onto the board.
func (b *Board) PlacePiece(p *Piece) {
	color := PieceColor(p.Type)
	for _, cell := range p.Cells() {
		if b.InBounds(cell) {
			b.Cells[cell.Row][cell.Col] = color
		}
	}
}

// ClearLines removes completed lines and returns how many were cleared
// and the row indices that were cleared.
func (b *Board) ClearLines() (int, []int) {
	var cleared []int
	for row := BoardHeight - 1; row >= 0; row-- {
		full := true
		for col := 0; col < BoardWidth; col++ {
			if b.Cells[row][col] == Empty {
				full = false
				break
			}
		}
		if full {
			cleared = append(cleared, row)
		}
	}

	if len(cleared) == 0 {
		return 0, nil
	}

	// Remove cleared rows and shift everything down.
	newCells := [BoardHeight][BoardWidth]CellColor{}
	writeRow := BoardHeight - 1
	for readRow := BoardHeight - 1; readRow >= 0; readRow-- {
		isClear := false
		for _, cr := range cleared {
			if readRow == cr {
				isClear = true
				break
			}
		}
		if !isClear {
			newCells[writeRow] = b.Cells[readRow]
			writeRow--
		}
	}
	b.Cells = newCells

	return len(cleared), cleared
}

// GhostPosition returns the position a piece would be at if hard dropped.
func (b *Board) GhostPosition(p *Piece) Position {
	ghost := p.Clone()
	for {
		ghost.Pos.Row++
		if !b.ValidPosition(&ghost) {
			ghost.Pos.Row--
			break
		}
	}
	return ghost.Pos
}

// IsAboveVisible checks if any occupied cell is in the buffer zone (above visible area).
func (b *Board) IsAboveVisible() bool {
	for row := 0; row < BufferRows; row++ {
		for col := 0; col < BoardWidth; col++ {
			if b.Cells[row][col] != Empty {
				return true
			}
		}
	}
	return false
}

// GetVisibleCell returns the color of a cell in visible coordinates (0 = top visible row).
func (b *Board) GetVisibleCell(visRow, col int) CellColor {
	return b.Cells[visRow+BufferRows][col]
}
