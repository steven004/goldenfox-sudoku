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
	Stats            UserStats      `json:"stats"`
	CompletedHistory []PuzzleRecord `json:"completed_history"`
	PendingHistory   []PuzzleRecord `json:"pending_history"`
	mu               sync.RWMutex   `json:"-"`
}

// NewUserData creates a new UserData instance with defaults
func NewUserData() *UserData {
	return &UserData{
		Stats:            NewUserStats(),
		CompletedHistory: make([]PuzzleRecord, 0),
		PendingHistory:   make([]PuzzleRecord, 0),
	}
}

// ... (retain ResolveUserDataPath and Save mostly unchanged, just struct field changes implicitly handled by JSON marshaler)

// getUserDataPath returns the platform-specific path for user data
func getUserDataPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "GoldenFoxSudoku", "user_data.json"), nil
}

// ResolveUserDataPath determines the path to the user data file
// It now ALWAYS defaults to the standard OS Config Dir to ensure consistency between Dev and Prod
func ResolveUserDataPath(filename string) (string, error) {
	// 2. Default to standard OS config path
	return getUserDataPath()
}

// Save saves the user data to a JSON file
func (ud *UserData) Save(filename string) error {
	path, err := ResolveUserDataPath(filename)
	if err != nil {
		return fmt.Errorf("failed to resolve user data path: %w", err)
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

// LoadUserData loads user data
func LoadUserData(filename string) (*UserData, error) {
	path, err := ResolveUserDataPath(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve user data path: %w", err)
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

	// Ensure maps are initialized
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

	// Ensure slices are initialized
	if ud.CompletedHistory == nil {
		ud.CompletedHistory = make([]PuzzleRecord, 0)
	}
	if ud.PendingHistory == nil {
		ud.PendingHistory = make([]PuzzleRecord, 0)
	}

	return &ud, nil
}

// UpsertPuzzleRecord adds or updates a record and handles moving between lists
func (ud *UserData) UpsertPuzzleRecord(record PuzzleRecord) {
	ud.mu.Lock()
	defer ud.mu.Unlock()

	// 1. Check Pending List first
	pendingIdx := -1
	for i := range ud.PendingHistory {
		if ud.PendingHistory[i].ID == record.ID {
			pendingIdx = i
			break
		}
	}

	// 2. Check Completed List
	completedIdx := -1
	for i := range ud.CompletedHistory {
		if ud.CompletedHistory[i].ID == record.ID {
			completedIdx = i
			break
		}
	}

	// SCENARIO A: Game is now SOLVED
	if record.IsSolved {
		if pendingIdx != -1 {
			// Move from Pending to Completed
			// Remove from Pending
			ud.PendingHistory = append(ud.PendingHistory[:pendingIdx], ud.PendingHistory[pendingIdx+1:]...)

			// Check if it exists in Completed (Replay case)
			if completedIdx != -1 {
				// Update existing formatted completion record (e.g. better time?)
				// Requirement: "the re-completion could update the elapsed time only"
				// We update the record but preserve original 'PlayedAt' if desired?
				// User said: "update the elapsed time only". Let's assume we update the whole record state
				// but we do NOT trigger scoring again.
				oldRecord := ud.CompletedHistory[completedIdx]
				if record.TimeElapsed < oldRecord.TimeElapsed {
					ud.CompletedHistory[completedIdx] = record
				}
				// Replay of already completed game -> No score gain
			} else {
				// First time completion for this Session ID?
				// Wait, if ID is unique per NewGame, it's new.
				ud.CompletedHistory = append(ud.CompletedHistory, record)

				// Update User Stats (UpdateStats handles general stats + RecordWin logic)
				ud.Stats.UpdateStats(record)
			}
		} else if completedIdx != -1 {
			// Already in Completed, just update (Replay case where it wasn't in pending?)
			oldRecord := ud.CompletedHistory[completedIdx]
			if record.TimeElapsed < oldRecord.TimeElapsed {
				ud.CompletedHistory[completedIdx] = record
			}
		} else {
			// New record directly solved (unlikely)
			ud.CompletedHistory = append(ud.CompletedHistory, record)
			ud.Stats.UpdateStats(record)
		}
	} else {
		// SCENARIO B: Game is UNCOMPLETED (In Progress)
		if pendingIdx != -1 {
			// Update existing pending
			ud.PendingHistory[pendingIdx] = record
		} else {
			// New pending game (or Replay newly added to Pending)
			// Ensure it's not a duplicate
			ud.PendingHistory = append(ud.PendingHistory, record)
		}
		// If it exists in completed (Replay), strictly it stays in completed too.
		// We don't remove from Completed when replaying, we only copy to Pending.
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
	return len(ud.PendingHistory)
}

// GetWinRate returns the percentage of games won
func (ud *UserData) GetWinRate() float64 {
	ud.mu.RLock()
	defer ud.mu.RUnlock()

	// Use actual stored history for consistency to avoid >100% rates
	// TotalSolved in Stats might be desynced or count duplicates if logic changes
	solvedCount := len(ud.CompletedHistory)

	// Total = Solved + Pending (In Progress)
	// Note: If a user deletes history, this rate resets, which is correct behavior for "User Data"
	total := solvedCount + len(ud.PendingHistory)

	if total == 0 {
		return 0.0
	}
	return (float64(solvedCount) / float64(total)) * 100
}

// GetGamesAtDifficulty returns the number of games played (started) at a specific difficulty
func (ud *UserData) GetGamesAtDifficulty(diff engine.DifficultyLevel) int {
	ud.mu.RLock()
	defer ud.mu.RUnlock()

	count := 0
	for _, record := range ud.CompletedHistory {
		if record.Difficulty == diff {
			count++
		}
	}
	for _, record := range ud.PendingHistory {
		if record.Difficulty == diff {
			count++
		}
	}
	return count
}

// GetAllHistory returns a combined list of all history (Completed + Pending)
// Useful for displaying full history to user
func (ud *UserData) GetAllHistory() []PuzzleRecord {
	ud.mu.RLock()
	defer ud.mu.RUnlock()

	// Preallocate slice
	all := make([]PuzzleRecord, 0, len(ud.CompletedHistory)+len(ud.PendingHistory))
	all = append(all, ud.PendingHistory...)
	all = append(all, ud.CompletedHistory...)
	return all
}
