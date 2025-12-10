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
// It prioritizes a local file if it exists (legacy/dev), otherwise defaults to OS Config Dir
func ResolveUserDataPath(filename string) (string, error) {
	// 1. If filename is provided, check if it exists locally
	if filename != "" {
		if _, err := os.Stat(filename); err == nil {
			return filename, nil
		}
	}

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
				// The "Replay" scenario implies we reused the ID.
				ud.CompletedHistory = append(ud.CompletedHistory, record)

				// Scoring Check:
				// "game behind current difficult or same difficult but behind progress, will not gain win score"
				// We need current difficulty/progress from Stats to compare against record's difficulty/context.
				// However, record ID encodes the context! But let's check current stats.

				// Current User State
				currentDiff := engine.DifficultyLevel(ud.Stats.Level - 1)
				if ud.Stats.Level > 6 {
					currentDiff = engine.FoxGod
				} // 6 is FoxGod in engine, Level is 1-6?
				// Actually engine defines Beginner=0..FoxGod=5. Level is 1..6.
				// userLevel 1 = Beginner(0), 6 = FoxGod(5).

				currentProgress := ud.Stats.Progress

				recordDiff := record.Difficulty
				// We need to decode Progress from ID or assume we compare against current?
				// "behind current difficult"
				isLowerDiff := recordDiff < currentDiff

				// "same difficult but behind progress"
				// We need the progress this game was created with.
				// We can parse it from ID or store it? PuzzleRecord doesn't store 'ProgressAtCreation'.
				// But we can infer it from the ID!
				// ID: D P ... (1st digit Diff+1, 2nd digit Prog+4)
				// Let's parse ID.

				shouldScore := true
				if len(record.ID) == 12 {
					// Parse Progress Char
					progChar := record.ID[1]         // '0'..'9'
					progVal := int(progChar-'0') - 4 // Convert back to -2..4

					if isLowerDiff {
						shouldScore = false
					} else if recordDiff == currentDiff {
						if progVal < currentProgress {
							shouldScore = false
						}
					}
				}

				if shouldScore {
					ud.updateStatsLocked(record)
				}
			}
		} else if completedIdx != -1 {
			// Already in Completed, just update (Replay case where it wasn't in pending?)
			// This happens if we loaded a completed game to replay, but didn't save it to pending yet?
			// User said: "replayed... also put into uncompleted list".
			// So it should have been in Pending if we did the load logic right.
			// But if we are here, just update.
			oldRecord := ud.CompletedHistory[completedIdx]
			if record.TimeElapsed < oldRecord.TimeElapsed {
				ud.CompletedHistory[completedIdx] = record
			}
		} else {
			// New record directly solved (unlikely)
			ud.CompletedHistory = append(ud.CompletedHistory, record)
			ud.updateStatsLocked(record)
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

// updateStatsLocked is a helper that assumes lock is held
func (ud *UserData) updateStatsLocked(record PuzzleRecord) {
	ud.Stats.TotalSolved++
	ud.Stats.SolvedCounts[record.Difficulty]++

	seconds := record.TimeElapsed.Seconds()
	ud.Stats.TotalTimes[record.Difficulty] += seconds

	count := float64(ud.Stats.SolvedCounts[record.Difficulty])
	ud.Stats.AverageTimes[record.Difficulty] = ud.Stats.TotalTimes[record.Difficulty] / count

	currentBest, exists := ud.Stats.BestTimes[record.Difficulty]
	if !exists || seconds < currentBest {
		ud.Stats.BestTimes[record.Difficulty] = seconds
	}

	// ---- Progress & Leveling Logic ----

	currentLevel := ud.Stats.Level
	currentDiffIndex := currentLevel - 1 // Level 1 -> 0 (Beginner)
	recordDiffIndex := int(record.Difficulty)

	// Calculate Progress Gain
	progressGain := 0

	if recordDiffIndex > currentDiffIndex {
		// Rule 2: Higher Difficulty -> Gain (Diff - CurrentLevel) + 1
		// Note: "Diff" here likely means the level value (1-6) or index (0-5)?
		// User said: "(difficult - currentLevel) + 1"
		// If Level 1 User plays Level 2 (Easy): (2 - 1) + 1 = 2 gain.
		// Using Indices: (1 - 0) + 1 = 2. Matches.
		progressGain = (recordDiffIndex - currentDiffIndex) + 1
	} else if recordDiffIndex == currentDiffIndex {
		// Same Difficulty
		// Check "Behind Progress" Rule
		// Need to extract original progress from ID
		gameProgress := 0 // Default
		if len(record.ID) == 12 {
			progChar := record.ID[1] // 2nd char
			if progChar >= '0' && progChar <= '9' {
				gameProgress = int(progChar-'0') - 4
			}
		}

		if gameProgress < ud.Stats.Progress {
			// Rule 1: Same diff but behind progress -> No Gain
			progressGain = 0
		} else {
			// Standard Gain
			progressGain = 1
		}
	} else {
		// Lower Difficulty
		// Rule 1: Lower -> No Gain
		progressGain = 0
	}

	// Apply Gain
	if progressGain > 0 {
		// Momentum Reset: If negative, jump to gain?
		// Previous logic: "If negative, reset to 1".
		// Let's refine: If negative, add gain to it? Or reset?
		// User's "Progress + 1" implies strictly adding.
		// But earlier "Reset to 1" was standard.
		// Let's stick to: If negative, reset to baseline + (gain-1)?
		// Simplest Interpretation: If < 0, Progress = progressGain. If >= 0, Progress += progressGain.

		if ud.Stats.Progress < 0 {
			ud.Stats.Progress = progressGain
		} else {
			ud.Stats.Progress += progressGain
		}
	}

	// Level Up Check
	// Max Progress is 4 (triggers at 5).
	// If we gain multiple points, we might skip.
	if ud.Stats.Progress >= 5 {
		if ud.Stats.Level < 6 {
			ud.Stats.Level++
			ud.Stats.Progress = 0 // Reset to 0 after level up
		} else {
			ud.Stats.Progress = 5 // Cap at max
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
	return len(ud.PendingHistory)
}

// GetWinRate returns the percentage of games won
func (ud *UserData) GetWinRate() float64 {
	ud.mu.RLock()
	defer ud.mu.RUnlock()

	total := len(ud.CompletedHistory) + len(ud.PendingHistory)
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
