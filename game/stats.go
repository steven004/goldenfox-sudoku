package game

import (
	"time"

	"github.com/steven004/goldenfox-sudoku/engine"
)

// UserStats holds aggregated statistics for the user
type UserStats struct {
	Level        int                                `json:"level"`    // User's experience level (1-6)
	Progress     int                                `json:"progress"` // -2 to +4 (at +5 level up, at -3 level down)
	JoinTime     time.Time                          `json:"join_time"`
	TotalSolved  int                                `json:"total_solved"`
	BestTimes    map[engine.DifficultyLevel]float64 `json:"best_times"`    // In seconds
	AverageTimes map[engine.DifficultyLevel]float64 `json:"average_times"` // In seconds
	TotalTimes   map[engine.DifficultyLevel]float64 `json:"total_times"`   // In seconds (helper for average)
	SolvedCounts map[engine.DifficultyLevel]int     `json:"solved_counts"` // Helper for average
}

// NewUserStats creates a new UserStats instance with defaults
func NewUserStats() UserStats {
	return UserStats{
		Level:        1,
		Progress:     0,
		JoinTime:     time.Now(),
		BestTimes:    make(map[engine.DifficultyLevel]float64),
		AverageTimes: make(map[engine.DifficultyLevel]float64),
		TotalTimes:   make(map[engine.DifficultyLevel]float64),
		SolvedCounts: make(map[engine.DifficultyLevel]int),
	}
}

// UpdateStats updates user statistics based on a solved record
func (us *UserStats) UpdateStats(record PuzzleRecord) {
	us.TotalSolved++

	if us.SolvedCounts == nil {
		us.SolvedCounts = make(map[engine.DifficultyLevel]int)
	}
	us.SolvedCounts[record.Difficulty]++

	seconds := record.TimeElapsed.Seconds()

	// Update Total Time for average calculation
	if us.TotalTimes == nil {
		us.TotalTimes = make(map[engine.DifficultyLevel]float64)
	}
	us.TotalTimes[record.Difficulty] += seconds

	// Recalculate Average
	count := float64(us.SolvedCounts[record.Difficulty])

	if us.AverageTimes == nil {
		us.AverageTimes = make(map[engine.DifficultyLevel]float64)
	}
	us.AverageTimes[record.Difficulty] = us.TotalTimes[record.Difficulty] / count

	// Update Best Time
	if us.BestTimes == nil {
		us.BestTimes = make(map[engine.DifficultyLevel]float64)
	}
	currentBest, exists := us.BestTimes[record.Difficulty]
	if !exists || seconds < currentBest {
		us.BestTimes[record.Difficulty] = seconds
	}

	// Level Up Logic
	// Win: Progress +1
	if us.Progress < 0 {
		us.Progress = 1 // Reset to +1 if coming from negative
	} else {
		us.Progress++
	}

	if us.Progress >= 5 {
		if us.Level < 6 { // Max Level 6 (FoxGod)
			us.Level++
			us.Progress = 0 // Reset after promotion
		} else {
			us.Progress = 5 // Cap at max progress for max level
		}
	}
}

// RecordLoss records a failed or abandoned game
func (us *UserStats) RecordLoss() {
	// Loss: Progress -1
	if us.Progress > 0 {
		us.Progress = -1 // Reset to -1 if coming from positive
	} else {
		us.Progress--
	}

	// Level Down Logic
	// Fail 3 consecutive games (Progress reaches -3) -> Level Down
	if us.Progress <= -3 {
		if us.Level > 1 { // Min Level 1 (Beginner)
			us.Level--
			us.Progress = 0 // Reset after demotion
		} else {
			us.Progress = -3 // Cap at min progress for min level
		}
	}
}
