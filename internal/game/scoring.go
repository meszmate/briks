package game

import "math"

// Scorer tracks score, level, lines, combos, and back-to-back state.
type Scorer struct {
	Score      int
	Level      int
	Lines      int
	Combo      int
	BackToBack bool
}

// NewScorer creates a scorer starting at the given level.
func NewScorer(startLevel int) *Scorer {
	return &Scorer{Level: startLevel}
}

// AddSoftDrop adds points for soft dropping.
func (s *Scorer) AddSoftDrop(cells int) {
	s.Score += cells
}

// AddHardDrop adds points for hard dropping.
func (s *Scorer) AddHardDrop(cells int) {
	s.Score += cells * 2
}

// AddLineClear processes a line clear event and returns the clear type.
func (s *Scorer) AddLineClear(linesCleared int, isTSpin bool) LineClearType {
	if linesCleared == 0 {
		s.Combo = 0
		return ClearNone
	}

	var clearType LineClearType
	var basePoints int

	if isTSpin {
		switch linesCleared {
		case 1:
			clearType = ClearTSpinSingle
			basePoints = 800
		case 2:
			clearType = ClearTSpinDouble
			basePoints = 1200
		case 3:
			clearType = ClearTSpinTriple
			basePoints = 1600
		}
	} else {
		switch linesCleared {
		case 1:
			clearType = ClearSingle
			basePoints = 100
		case 2:
			clearType = ClearDouble
			basePoints = 300
		case 3:
			clearType = ClearTriple
			basePoints = 500
		case 4:
			clearType = ClearTetris
			basePoints = 800
		}
	}

	points := basePoints * s.Level

	// Back-to-back bonus for Tetris or T-Spin clears.
	isDifficult := clearType == ClearTetris || isTSpin
	if isDifficult && s.BackToBack {
		points = points * 3 / 2
	}
	if isDifficult {
		s.BackToBack = true
	} else {
		s.BackToBack = false
	}

	// Combo bonus.
	if s.Combo > 0 {
		points += 50 * s.Combo * s.Level
	}
	s.Combo++

	s.Score += points
	s.Lines += linesCleared

	// Level up every 10 lines.
	newLevel := s.Lines/10 + 1
	if newLevel > s.Level {
		s.Level = newLevel
	}

	return clearType
}

// GravityInterval returns the gravity interval in seconds for the current level.
// Formula: (0.8 - ((level-1) * 0.007)) ^ (level-1)
func (s *Scorer) GravityInterval() float64 {
	level := float64(s.Level)
	interval := math.Pow(0.8-((level-1)*0.007), level-1)
	if interval < 0.01 {
		interval = 0.01
	}
	return interval
}
