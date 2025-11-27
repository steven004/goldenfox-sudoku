package game

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/steven004/goldenfox-sudoku/engine"
)

// PuzzleRecord represents a single played puzzle
type PuzzleRecord struct {
	ID          string                 `json:"id"`
	Predefined  string                 `json:"predefined"`  // String representation of initial board
	FinalState  string                 `json:"final_state"` // String representation of final board
	IsSolved    bool                   `json:"is_solved"`
	TimeElapsed time.Duration          `json:"time_elapsed"`
	PlayedAt    time.Time              `json:"played_at"`
	Difficulty  engine.DifficultyLevel `json:"difficulty"`
	Mistakes    int                    `json:"mistakes"`
}

// UserStats holds aggregated statistics for the user
type UserStats struct {
	Level        int                                `json:"level"` // User's experience level
	JoinTime     time.Time                          `json:"join_time"`
	TotalSolved  int                                `json:"total_solved"`
	BestTimes    map[engine.DifficultyLevel]float64 `json:"best_times"`    // In seconds
	AverageTimes map[engine.DifficultyLevel]float64 `json:"average_times"` // In seconds
	TotalTimes   map[engine.DifficultyLevel]float64 `json:"total_times"`   // In seconds (helper for average)
	SolvedCounts map[engine.DifficultyLevel]int     `json:"solved_counts"` // Helper for average
}

// UserData is the root structure for user persistence
type UserData struct {
	Stats   UserStats      `json:"stats"`
	History []PuzzleRecord `json:"history"`
	mu      sync.RWMutex   `json:"-"`
}

// NewUserData creates a new UserData instance with defaults
func NewUserData() *UserData {
	return &UserData{
		Stats: UserStats{
			Level:        1,
			JoinTime:     time.Now(),
			BestTimes:    make(map[engine.DifficultyLevel]float64),
			AverageTimes: make(map[engine.DifficultyLevel]float64),
			TotalTimes:   make(map[engine.DifficultyLevel]float64),
			SolvedCounts: make(map[engine.DifficultyLevel]int),
		},
		History: make([]PuzzleRecord, 0),
	}
}

// Save saves the user data to a JSON file
func (ud *UserData) Save(filename string) error {
	ud.mu.RLock()
	defer ud.mu.RUnlock()

	data, err := json.MarshalIndent(ud, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal user data: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write user data file: %w", err)
	}

	return nil
}

// LoadUserData loads user data from a JSON file
func LoadUserData(filename string) (*UserData, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return NewUserData(), nil // Return new if not exists
		}
		return nil, fmt.Errorf("failed to read user data file: %w", err)
	}

	var ud UserData
	if err := json.Unmarshal(data, &ud); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	// Initialize maps if nil (in case of empty JSON or old version)
	if ud.Stats.BestTimes == nil {
		ud.Stats.BestTimes = make(map[engine.DifficultyLevel]float64)
	}
	if ud.Stats.AverageTimes == nil {
		ud.Stats.AverageTimes = make(map[engine.DifficultyLevel]float64)
	}
	if ud.Stats.TotalTimes == nil {
		ud.Stats.TotalTimes = make(map[engine.DifficultyLevel]float64)
	}
	if ud.Stats.SolvedCounts == nil {
		ud.Stats.SolvedCounts = make(map[engine.DifficultyLevel]int)
	}

	return &ud, nil
}

// AddPuzzleRecord adds a record and updates stats
func (ud *UserData) AddPuzzleRecord(record PuzzleRecord) {
	ud.mu.Lock()
	defer ud.mu.Unlock()

	ud.History = append(ud.History, record)

	if record.IsSolved {
		ud.Stats.TotalSolved++
		ud.Stats.SolvedCounts[record.Difficulty]++

		seconds := record.TimeElapsed.Seconds()

		// Update Total Time for average calculation
		ud.Stats.TotalTimes[record.Difficulty] += seconds

		// Recalculate Average
		count := float64(ud.Stats.SolvedCounts[record.Difficulty])
		ud.Stats.AverageTimes[record.Difficulty] = ud.Stats.TotalTimes[record.Difficulty] / count

		// Update Best Time
		currentBest, exists := ud.Stats.BestTimes[record.Difficulty]
		if !exists || seconds < currentBest {
			ud.Stats.BestTimes[record.Difficulty] = seconds
		}

		// Simple Level Up Logic (e.g., every 5 puzzles)
		// This is a placeholder for more complex logic
		if ud.Stats.TotalSolved%5 == 0 {
			ud.Stats.Level++
		}
	}
}

// UpsertPuzzleRecord adds or updates a record and updates stats if newly solved
func (ud *UserData) UpsertPuzzleRecord(record PuzzleRecord) {
	ud.mu.Lock()
	defer ud.mu.Unlock()

	var existingIdx = -1
	for i := range ud.History {
		if ud.History[i].ID == record.ID {
			existingIdx = i
			break
		}
	}

	if existingIdx != -1 {
		// Update existing
		oldRecord := ud.History[existingIdx]
		ud.History[existingIdx] = record

		// If it wasn't solved before but is now, update stats
		if !oldRecord.IsSolved && record.IsSolved {
			ud.updateStats(record)
		}
	} else {
		// Append new
		ud.History = append(ud.History, record)
		if record.IsSolved {
			ud.updateStats(record)
		}
	}
}

// updateStats updates user statistics based on a solved record
func (ud *UserData) updateStats(record PuzzleRecord) {
	ud.Stats.TotalSolved++
	ud.Stats.SolvedCounts[record.Difficulty]++

	seconds := record.TimeElapsed.Seconds()

	// Update Total Time for average calculation
	ud.Stats.TotalTimes[record.Difficulty] += seconds

	// Recalculate Average
	count := float64(ud.Stats.SolvedCounts[record.Difficulty])
	ud.Stats.AverageTimes[record.Difficulty] = ud.Stats.TotalTimes[record.Difficulty] / count

	// Update Best Time
	currentBest, exists := ud.Stats.BestTimes[record.Difficulty]
	if !exists || seconds < currentBest {
		ud.Stats.BestTimes[record.Difficulty] = seconds
	}

	// Simple Level Up Logic
	if ud.Stats.TotalSolved%5 == 0 {
		ud.Stats.Level++
	}
}
