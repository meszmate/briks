package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const highscoreFile = "highscores.json"

const MaxHighScores = 10

// HighScore represents a single high score entry.
type HighScore struct {
	Score  int       `json:"score"`
	Level  int       `json:"level"`
	Lines  int       `json:"lines"`
	Pieces int       `json:"pieces"`
	Date   time.Time `json:"date"`
}

// HighScores manages the top scores list.
type HighScores struct {
	Scores []HighScore `json:"scores"`
}

func highscorePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configDir, highscoreFile), nil
}

// LoadHighScores reads high scores from disk.
func LoadHighScores() *HighScores {
	path, err := highscorePath()
	if err != nil {
		return &HighScores{}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return &HighScores{}
	}

	hs := &HighScores{}
	if err := json.Unmarshal(data, hs); err != nil {
		return &HighScores{}
	}

	return hs
}

// Save writes high scores to disk.
func (hs *HighScores) Save() error {
	path, err := highscorePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(hs, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// Add inserts a score and returns its rank (1-based), or 0 if it didn't make the list.
func (hs *HighScores) Add(score HighScore) int {
	hs.Scores = append(hs.Scores, score)
	sort.Slice(hs.Scores, func(i, j int) bool {
		return hs.Scores[i].Score > hs.Scores[j].Score
	})

	if len(hs.Scores) > MaxHighScores {
		hs.Scores = hs.Scores[:MaxHighScores]
	}

	for i, s := range hs.Scores {
		if s.Score == score.Score && s.Date.Equal(score.Date) {
			return i + 1
		}
	}

	return 0
}

// IsHighScore checks if a score would make the top list.
func (hs *HighScores) IsHighScore(score int) bool {
	if len(hs.Scores) < MaxHighScores {
		return true
	}
	return score > hs.Scores[len(hs.Scores)-1].Score
}
