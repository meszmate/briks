package game

import "time"

const (
	MaxLockResets = 15
	LockDelay     = 500 * time.Millisecond
)

// Engine orchestrates the game: board, active piece, bag, scorer, state.
type Engine struct {
	Board        *Board
	Bag          *Bag
	Scorer       *Scorer
	State        GameState
	Current      *Piece
	HoldPiece    *PieceType
	HoldUsed     bool
	PreviewCount int

	// Lock delay tracking.
	LockTimer   time.Time
	LockResets  int
	LockStarted bool

	// T-Spin detection.
	LastMoveWasRotation bool

	// Stats.
	PiecesPlaced int
	StartTime    time.Time
}

// NewEngine creates a new game engine.
func NewEngine(startLevel, previewCount int) *Engine {
	e := &Engine{
		Board:        NewBoard(),
		Bag:          NewBag(),
		Scorer:       NewScorer(startLevel),
		State:        StatePlaying,
		PreviewCount: previewCount,
		StartTime:    time.Now(),
	}
	e.spawnPiece()
	return e
}

// spawnPiece pulls the next piece from the bag and places it at the spawn position.
// Returns false if the piece can't be placed (game over).
func (e *Engine) spawnPiece() bool {
	pt := e.Bag.Next()
	p := &Piece{
		Type:     pt,
		Rotation: Rot0,
		Pos:      SpawnPosition(pt),
	}

	if !e.Board.ValidPosition(p) {
		e.State = StateGameOver
		return false
	}

	e.Current = p
	e.HoldUsed = false
	e.LockStarted = false
	e.LockResets = 0
	e.LastMoveWasRotation = false
	return true
}

// NextPieces returns the upcoming pieces for preview.
func (e *Engine) NextPieces() []PieceType {
	return e.Bag.Preview(e.PreviewCount)
}

// canMoveDown checks if the current piece can move down.
func (e *Engine) canMoveDown() bool {
	if e.Current == nil {
		return false
	}
	test := e.Current.Clone()
	test.Pos.Row++
	return e.Board.ValidPosition(&test)
}

// MoveLeft moves the current piece left.
func (e *Engine) MoveLeft() bool {
	if e.State != StatePlaying || e.Current == nil {
		return false
	}
	test := e.Current.Clone()
	test.Pos.Col--
	if e.Board.ValidPosition(&test) {
		e.Current.Pos.Col--
		e.LastMoveWasRotation = false
		e.resetLockIfNeeded()
		return true
	}
	return false
}

// MoveRight moves the current piece right.
func (e *Engine) MoveRight() bool {
	if e.State != StatePlaying || e.Current == nil {
		return false
	}
	test := e.Current.Clone()
	test.Pos.Col++
	if e.Board.ValidPosition(&test) {
		e.Current.Pos.Col++
		e.LastMoveWasRotation = false
		e.resetLockIfNeeded()
		return true
	}
	return false
}

// MoveDown moves the current piece down by one.
// Returns true if successful.
func (e *Engine) MoveDown() bool {
	if e.State != StatePlaying || e.Current == nil {
		return false
	}
	test := e.Current.Clone()
	test.Pos.Row++
	if e.Board.ValidPosition(&test) {
		e.Current.Pos.Row++
		e.LastMoveWasRotation = false
		// Reset lock when piece moves down
		e.LockStarted = false
		e.LockResets = 0
		return true
	}
	return false
}

// SoftDrop moves down and awards soft drop points.
func (e *Engine) SoftDrop() bool {
	if e.MoveDown() {
		e.Scorer.AddSoftDrop(1)
		return true
	}
	return false
}

// HardDrop instantly drops and locks the piece.
func (e *Engine) HardDrop() LineClearType {
	if e.State != StatePlaying || e.Current == nil {
		return ClearNone
	}
	dropDist := 0
	for e.canMoveDown() {
		e.Current.Pos.Row++
		dropDist++
	}
	e.Scorer.AddHardDrop(dropDist)
	return e.lockPiece()
}

// RotateCW rotates the piece clockwise using SRS wall kicks.
func (e *Engine) RotateCW() bool {
	if e.State != StatePlaying || e.Current == nil {
		return false
	}
	return e.rotate(nextRotCW(e.Current.Rotation))
}

// RotateCCW rotates the piece counter-clockwise using SRS wall kicks.
func (e *Engine) RotateCCW() bool {
	if e.State != StatePlaying || e.Current == nil {
		return false
	}
	return e.rotate(nextRotCCW(e.Current.Rotation))
}

func (e *Engine) rotate(newRot Rotation) bool {
	kicks := GetWallKicks(e.Current.Type, e.Current.Rotation, newRot)
	for _, kick := range kicks {
		test := e.Current.Clone()
		test.Rotation = newRot
		test.Pos.Col += kick.Col
		test.Pos.Row -= kick.Row // SRS: positive row = up
		if e.Board.ValidPosition(&test) {
			e.Current.Rotation = test.Rotation
			e.Current.Pos = test.Pos
			e.LastMoveWasRotation = true
			e.resetLockIfNeeded()
			return true
		}
	}
	return false
}

// Hold swaps the current piece with the hold piece.
func (e *Engine) Hold() bool {
	if e.State != StatePlaying || e.Current == nil || e.HoldUsed {
		return false
	}
	currentType := e.Current.Type
	if e.HoldPiece != nil {
		// Swap with held piece.
		heldType := *e.HoldPiece
		e.HoldPiece = &currentType
		p := &Piece{
			Type:     heldType,
			Rotation: Rot0,
			Pos:      SpawnPosition(heldType),
		}
		if !e.Board.ValidPosition(p) {
			e.State = StateGameOver
			return false
		}
		e.Current = p
	} else {
		e.HoldPiece = &currentType
		e.spawnPiece()
	}
	e.HoldUsed = true
	e.LockStarted = false
	e.LockResets = 0
	e.LastMoveWasRotation = false
	return true
}

// Tick advances the game by one gravity step.
func (e *Engine) Tick() LineClearType {
	if e.State != StatePlaying || e.Current == nil {
		return ClearNone
	}

	// Try to move down
	if e.MoveDown() {
		return ClearNone
	}

	// Piece can't move down - start or continue lock delay
	if !e.LockStarted {
		e.LockStarted = true
		e.LockTimer = time.Now()
		return ClearNone
	}

	// Check if lock delay expired
	if time.Since(e.LockTimer) >= LockDelay {
		return e.lockPiece()
	}

	return ClearNone
}

// CheckLock checks if the lock delay has expired.
func (e *Engine) CheckLock() LineClearType {
	if e.State != StatePlaying || e.Current == nil || !e.LockStarted {
		return ClearNone
	}

	// If piece can now move down, cancel lock
	if e.canMoveDown() {
		e.LockStarted = false
		e.LockResets = 0
		return ClearNone
	}

	// Check if lock delay expired
	if time.Since(e.LockTimer) >= LockDelay {
		return e.lockPiece()
	}

	return ClearNone
}

// GhostPosition returns where the current piece would land.
func (e *Engine) GhostPosition() Position {
	if e.Current == nil {
		return Position{}
	}
	return e.Board.GhostPosition(e.Current)
}

// GhostCells returns the ghost piece cells.
func (e *Engine) GhostCells() []Position {
	if e.Current == nil {
		return nil
	}
	ghost := e.Current.Clone()
	ghost.Pos = e.GhostPosition()
	return ghost.Cells()
}

func (e *Engine) lockPiece() LineClearType {
	if e.Current == nil {
		return ClearNone
	}

	// Detect T-spin before placing.
	isTSpin := e.detectTSpin()

	// Place the piece on the board
	e.Board.PlacePiece(e.Current)
	e.PiecesPlaced++

	// Clear lines
	linesCleared, _ := e.Board.ClearLines()
	clearType := e.Scorer.AddLineClear(linesCleared, isTSpin)

	// Spawn next piece
	if !e.spawnPiece() {
		return clearType
	}

	return clearType
}

func (e *Engine) detectTSpin() bool {
	if e.Current.Type != PieceT || !e.LastMoveWasRotation {
		return false
	}

	// Check 3 of 4 corners of the T's bounding box are occupied.
	row := e.Current.Pos.Row
	col := e.Current.Pos.Col
	corners := []Position{
		{row, col},
		{row, col + 2},
		{row + 2, col},
		{row + 2, col + 2},
	}

	occupied := 0
	for _, c := range corners {
		if !e.Board.InBounds(c) || e.Board.Cells[c.Row][c.Col] != Empty {
			occupied++
		}
	}

	return occupied >= 3
}

func (e *Engine) resetLockIfNeeded() {
	if e.LockStarted && e.LockResets < MaxLockResets {
		e.LockTimer = time.Now()
		e.LockResets++
	}
}

func nextRotCW(r Rotation) Rotation {
	return (r + 1) % 4
}

func nextRotCCW(r Rotation) Rotation {
	return (r + 3) % 4
}

// ElapsedTime returns the game duration.
func (e *Engine) ElapsedTime() time.Duration {
	return time.Since(e.StartTime)
}

// PiecesPerSecond returns the placement rate.
func (e *Engine) PiecesPerSecond() float64 {
	elapsed := e.ElapsedTime().Seconds()
	if elapsed == 0 {
		return 0
	}
	return float64(e.PiecesPlaced) / elapsed
}
