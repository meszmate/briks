package game

// CellColor represents the color type of a cell on the board.
type CellColor int

const (
	Empty CellColor = iota
	ColorI
	ColorO
	ColorT
	ColorS
	ColorZ
	ColorJ
	ColorL
	ColorGhost
)

// GameState represents the current state of the game.
type GameState int

const (
	StatePlaying GameState = iota
	StatePaused
	StateGameOver
)

// LineClearType categorizes how lines were cleared.
type LineClearType int

const (
	ClearNone LineClearType = iota
	ClearSingle
	ClearDouble
	ClearTriple
	ClearTetris
	ClearTSpinSingle
	ClearTSpinDouble
	ClearTSpinTriple
)

// Position represents a row/column coordinate.
type Position struct {
	Row int
	Col int
}

// Rotation represents a piece's rotation state (0–3).
type Rotation int

const (
	Rot0 Rotation = iota
	Rot1
	Rot2
	Rot3
)

// PieceType identifies a tetromino piece.
type PieceType int

const (
	PieceI PieceType = iota
	PieceO
	PieceT
	PieceS
	PieceZ
	PieceJ
	PieceL
)

var AllPieceTypes = []PieceType{PieceI, PieceO, PieceT, PieceS, PieceZ, PieceJ, PieceL}

// PieceColor returns the CellColor for a given piece type.
func PieceColor(p PieceType) CellColor {
	switch p {
	case PieceI:
		return ColorI
	case PieceO:
		return ColorO
	case PieceT:
		return ColorT
	case PieceS:
		return ColorS
	case PieceZ:
		return ColorZ
	case PieceJ:
		return ColorJ
	case PieceL:
		return ColorL
	default:
		return Empty
	}
}

// PieceName returns the name string for a piece type.
func PieceName(p PieceType) string {
	switch p {
	case PieceI:
		return "I"
	case PieceO:
		return "O"
	case PieceT:
		return "T"
	case PieceS:
		return "S"
	case PieceZ:
		return "Z"
	case PieceJ:
		return "J"
	case PieceL:
		return "L"
	default:
		return "?"
	}
}
