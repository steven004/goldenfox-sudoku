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
// It aggregates general stats (times, counts) and calls RecordWin for progress logic.
func (us *UserStats) UpdateStats(record PuzzleRecord) {
	// 1. General Stats Updates
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

	// 2. Calculate Progress (Ranking)
	us.RecordWin(record)
}

// RecordWin handles the progress and leveling logic for a won game
func (us *UserStats) RecordWin(record PuzzleRecord) {
	// Challenger Level Up Logic
	// Calculate the gap between played difficulty and user level
	// Beginner(0) vs Level 1 (Beginner) -> Gap 0?
	// Note: difficulty enum might need int casting or mapping to be comparable to Level (int).
	// Engine: Beginner=0, Easy=1, but User Level starts at 1.
	playedLevel := int(record.Difficulty) + 1
	diffDelta := playedLevel - us.Level

	// BASE PROGRESS
	progressGain := 1

	if diffDelta > 0 {
		// Playing "Up" -> Bonus Points
		// e.g. Lv.1 plays Hard(4) -> Delta = 3 -> Gain = 1 + 3 = 4
		progressGain += diffDelta // Corrected formula: (Diff - Level) + 1 implies Delta + 1 if Delta is (Diff - Level)
		// Wait, my previous formula was `progressGain += diffDelta`.
		// If base was 1, then 1 + Delta. Correct.
	} else if diffDelta == 0 {
		// Playing "Fair" -> Standard Gain (1) unless Stale

		// Check "Behind Progress" Rule (Stale Game)
		gameProgress := 0 // Default
		if len(record.ID) == 12 {
			progChar := record.ID[1] // 2nd char
			if progChar >= '0' && progChar <= '9' {
				gameProgress = int(progChar-'0') - 4
			}
		}

		if gameProgress < us.Progress {
			// Rule: Same diff but behind progress -> No Gain
			progressGain = 0
		}
	} else {
		// Playing "Down" (Smurfing) -> No Gain
		progressGain = 0
	}

	// Apply Gain
	if progressGain > 0 {
		if us.Progress < 0 {
			us.Progress = progressGain // Reset to positive gain immediately
		} else {
			us.Progress += progressGain
		}
	}

	// Check Promotion
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
func (us *UserStats) RecordLoss(difficulty engine.DifficultyLevel) {
	// Calculate Gap
	playedLevel := int(difficulty) + 1
	diffDelta := playedLevel - us.Level

	// BASE PENALTY
	penalty := 1

	if diffDelta > 0 {
		// Playing "Up" -> No Penalty
		// If you try Hard and fail, it's okay.
		penalty = 0
	} else if diffDelta < 0 {
		// Playing "Down" and failing -> Massive Penalty
		// Lv.5 loses to Easy -> Embarrassing.
		penalty = 2
	}

	// Apply Penalty
	if penalty > 0 {
		if us.Progress > 0 {
			us.Progress = -penalty // Reset to negative immediately
		} else {
			us.Progress -= penalty
		}
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
