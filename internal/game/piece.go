package game

// Piece represents an active tetromino on the board.
type Piece struct {
	Type     PieceType
	Rotation Rotation
	Pos      Position // top-left corner of bounding box
}

// Cells returns the filled cell positions for the piece's current rotation state.
func (p *Piece) Cells() []Position {
	offsets := PieceRotations[p.Type][p.Rotation]
	cells := make([]Position, len(offsets))
	for i, off := range offsets {
		cells[i] = Position{Row: p.Pos.Row + off.Row, Col: p.Pos.Col + off.Col}
	}
	return cells
}

// Clone returns a copy of the piece.
func (p *Piece) Clone() Piece {
	return Piece{Type: p.Type, Rotation: p.Rotation, Pos: p.Pos}
}

// PieceRotations defines the cell offsets for each piece type and rotation.
// Offsets are relative to the piece's position (top-left of bounding box).
var PieceRotations = map[PieceType][4][]Position{
	PieceI: {
		// Rot0 - horizontal I
		{{0, 0}, {0, 1}, {0, 2}, {0, 3}},
		// Rot1
		{{0, 2}, {1, 2}, {2, 2}, {3, 2}},
		// Rot2
		{{2, 0}, {2, 1}, {2, 2}, {2, 3}},
		// Rot3
		{{0, 1}, {1, 1}, {2, 1}, {3, 1}},
	},
	PieceO: {
		{{0, 0}, {0, 1}, {1, 0}, {1, 1}},
		{{0, 0}, {0, 1}, {1, 0}, {1, 1}},
		{{0, 0}, {0, 1}, {1, 0}, {1, 1}},
		{{0, 0}, {0, 1}, {1, 0}, {1, 1}},
	},
	PieceT: {
		{{0, 1}, {1, 0}, {1, 1}, {1, 2}},
		{{0, 1}, {1, 1}, {1, 2}, {2, 1}},
		{{1, 0}, {1, 1}, {1, 2}, {2, 1}},
		{{0, 1}, {1, 0}, {1, 1}, {2, 1}},
	},
	PieceS: {
		{{0, 1}, {0, 2}, {1, 0}, {1, 1}},
		{{0, 1}, {1, 1}, {1, 2}, {2, 2}},
		{{1, 1}, {1, 2}, {2, 0}, {2, 1}},
		{{0, 0}, {1, 0}, {1, 1}, {2, 1}},
	},
	PieceZ: {
		{{0, 0}, {0, 1}, {1, 1}, {1, 2}},
		{{0, 2}, {1, 1}, {1, 2}, {2, 1}},
		{{1, 0}, {1, 1}, {2, 1}, {2, 2}},
		{{0, 1}, {1, 0}, {1, 1}, {2, 0}},
	},
	PieceJ: {
		{{0, 0}, {1, 0}, {1, 1}, {1, 2}},
		{{0, 1}, {0, 2}, {1, 1}, {2, 1}},
		{{1, 0}, {1, 1}, {1, 2}, {2, 2}},
		{{0, 1}, {1, 1}, {2, 0}, {2, 1}},
	},
	PieceL: {
		{{0, 2}, {1, 0}, {1, 1}, {1, 2}},
		{{0, 1}, {1, 1}, {2, 1}, {2, 2}},
		{{1, 0}, {1, 1}, {1, 2}, {2, 0}},
		{{0, 0}, {0, 1}, {1, 1}, {2, 1}},
	},
}

// SRS wall kick data.
// WallKicksJLSTZ is for J, L, S, T, Z pieces.
var WallKicksJLSTZ = map[[2]Rotation][]Position{
	{Rot0, Rot1}: {{0, 0}, {-1, 0}, {-1, -1}, {0, 2}, {-1, 2}},
	{Rot1, Rot0}: {{0, 0}, {1, 0}, {1, 1}, {0, -2}, {1, -2}},
	{Rot1, Rot2}: {{0, 0}, {1, 0}, {1, 1}, {0, -2}, {1, -2}},
	{Rot2, Rot1}: {{0, 0}, {-1, 0}, {-1, -1}, {0, 2}, {-1, 2}},
	{Rot2, Rot3}: {{0, 0}, {1, 0}, {1, -1}, {0, 2}, {1, 2}},
	{Rot3, Rot2}: {{0, 0}, {-1, 0}, {-1, 1}, {0, -2}, {-1, -2}},
	{Rot3, Rot0}: {{0, 0}, {-1, 0}, {-1, 1}, {0, -2}, {-1, -2}},
	{Rot0, Rot3}: {{0, 0}, {1, 0}, {1, -1}, {0, 2}, {1, 2}},
}

// WallKicksI is for the I piece.
var WallKicksI = map[[2]Rotation][]Position{
	{Rot0, Rot1}: {{0, 0}, {-2, 0}, {1, 0}, {-2, 1}, {1, -2}},
	{Rot1, Rot0}: {{0, 0}, {2, 0}, {-1, 0}, {2, -1}, {-1, 2}},
	{Rot1, Rot2}: {{0, 0}, {-1, 0}, {2, 0}, {-1, -2}, {2, 1}},
	{Rot2, Rot1}: {{0, 0}, {1, 0}, {-2, 0}, {1, 2}, {-2, -1}},
	{Rot2, Rot3}: {{0, 0}, {2, 0}, {-1, 0}, {2, -1}, {-1, 2}},
	{Rot3, Rot2}: {{0, 0}, {-2, 0}, {1, 0}, {-2, 1}, {1, -2}},
	{Rot3, Rot0}: {{0, 0}, {1, 0}, {-2, 0}, {1, 2}, {-2, -1}},
	{Rot0, Rot3}: {{0, 0}, {-1, 0}, {2, 0}, {-1, -2}, {2, 1}},
}

// GetWallKicks returns the wall kick offsets for a rotation attempt.
func GetWallKicks(pieceType PieceType, from, to Rotation) []Position {
	key := [2]Rotation{from, to}
	if pieceType == PieceI {
		return WallKicksI[key]
	}
	if pieceType == PieceO {
		return []Position{{0, 0}}
	}
	return WallKicksJLSTZ[key]
}

// SpawnPosition returns the starting position for a piece type.
// Pieces spawn at the top of the visible area, centered horizontally.
func SpawnPosition(pt PieceType) Position {
	// BufferRows = 4, so visible area starts at row 4
	// Spawn pieces so their top cells are at row 4 (first visible row)
	switch pt {
	case PieceI:
		// I piece in Rot0 has all cells at row 0, so Pos.Row = 4 puts cells at row 4
		return Position{Row: BufferRows, Col: 3}
	case PieceO:
		// O piece has cells at rows 0-1, so Pos.Row = 4 puts cells at rows 4-5
		return Position{Row: BufferRows, Col: 4}
	default:
		// Most pieces have cells at rows 0-1, so Pos.Row = 4 puts cells at rows 4-5
		return Position{Row: BufferRows, Col: 3}
	}
}
