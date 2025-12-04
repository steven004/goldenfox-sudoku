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
	Level        int                                `json:"level"`    // User's experience level (1-6)
	Progress     int                                `json:"progress"` // -2 to +4 (at +5 level up, at -3 level down)
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
			Progress:     0,
			JoinTime:     time.Now(),
			BestTimes:    make(map[engine.DifficultyLevel]float64),
			AverageTimes: make(map[engine.DifficultyLevel]float64),
			TotalTimes:   make(map[engine.DifficultyLevel]float64),
			SolvedCounts: make(map[engine.DifficultyLevel]int),
		},
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
			ud.updateStats(record)
		}
	} else {
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

	// Level Up Logic
	// Win: Progress +1
	if ud.Stats.Progress < 0 {
		ud.Stats.Progress = 1 // Reset to +1 if coming from negative
	} else {
		ud.Stats.Progress++
	}

	if ud.Stats.Progress >= 5 {
		if ud.Stats.Level < 6 { // Max Level 6 (FoxGod)
			ud.Stats.Level++
			ud.Stats.Progress = 0 // Reset after promotion
		} else {
			ud.Stats.Progress = 5 // Cap at max progress for max level
		}
	}
}

// RecordLoss records a failed or abandoned game
func (ud *UserData) RecordLoss() {
	ud.mu.Lock()
	defer ud.mu.Unlock()

	// Loss: Progress -1
	if ud.Stats.Progress > 0 {
		ud.Stats.Progress = -1 // Reset to -1 if coming from positive
	} else {
		ud.Stats.Progress--
	}

	// Level Down Logic
	// Fail 3 consecutive games (Progress reaches -3) -> Level Down
	if ud.Stats.Progress <= -3 {
		if ud.Stats.Level > 1 { // Min Level 1 (Beginner)
			ud.Stats.Level--
			ud.Stats.Progress = 0 // Reset after demotion
		} else {
			ud.Stats.Progress = -3 // Cap at min progress for min level
		}
	}
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
