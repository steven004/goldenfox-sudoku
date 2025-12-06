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
	ID              string                 `json:"id"`
	Predefined      string                 `json:"predefined"`  // String representation of initial board
	FinalState      string                 `json:"final_state"` // String representation of final board
	IsSolved        bool                   `json:"is_solved"`
	TimeElapsed     time.Duration          `json:"time_elapsed"`
	PlayedAt        time.Time              `json:"played_at"`
	Difficulty      engine.DifficultyLevel `json:"difficulty"`
	DifficultyIndex float64                `json:"difficulty_index"`
	Mistakes        int                    `json:"mistakes"`
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
		Stats:   NewUserStats(),
		History: make([]PuzzleRecord, 0),
	}
}

// getUserDataPath returns the platform-specific path for user data
func getUserDataPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "GoldenFoxSudoku", "user_data.json"), nil
}

// Save saves the user data to a JSON file
func (ud *UserData) Save(filename string) error {
	path, err := getUserDataPath()
	if err != nil {
		return fmt.Errorf("failed to get user data path: %w", err)
	}

	ud.mu.RLock()
	defer ud.mu.RUnlock()

	data, err := json.MarshalIndent(ud, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal user data: %w", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write user data file: %w", err)
	}

	return nil
}

// LoadUserData loads user data from the standard location
func LoadUserData(filename string) (*UserData, error) {
	path, err := getUserDataPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get user data path: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewUserData(), nil
		}
		return nil, fmt.Errorf("failed to read user data file: %w", err)
	}

	var ud UserData
	if err := json.Unmarshal(data, &ud); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	// Ensure maps are initialized (in case of loading old data or empty maps)
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
		oldRecord := ud.History[existingIdx]
		ud.History[existingIdx] = record
		if !oldRecord.IsSolved && record.IsSolved {
			ud.Stats.UpdateStats(record)
		}
	} else {
		ud.History = append(ud.History, record)
		if record.IsSolved {
			ud.Stats.UpdateStats(record)
		}
	}
}

// RecordLoss records a failed or abandoned game
func (ud *UserData) RecordLoss(difficulty engine.DifficultyLevel) {
	ud.mu.Lock()
	defer ud.mu.Unlock()
	ud.Stats.RecordLoss(difficulty)
}

// GetPendingGamesCount returns the number of unfinished games in history
func (ud *UserData) GetPendingGamesCount() int {
	ud.mu.RLock()
	defer ud.mu.RUnlock()

	count := 0
	for _, record := range ud.History {
		if !record.IsSolved {
			count++
		}
	}
	return count
}

// GetWinRate returns the percentage of games won
func (ud *UserData) GetWinRate() float64 {
	ud.mu.RLock()
	defer ud.mu.RUnlock()

	total := len(ud.History)
	if total == 0 {
		return 0.0
	}
	return (float64(ud.Stats.TotalSolved) / float64(total)) * 100
}

// GetGamesAtDifficulty returns the number of games played (started) at a specific difficulty
func (ud *UserData) GetGamesAtDifficulty(diff engine.DifficultyLevel) int {
	ud.mu.RLock()
	defer ud.mu.RUnlock()

	count := 0
	for _, record := range ud.History {
		if record.Difficulty == diff {
			count++
		}
	}
	return count
}
