package game

import (
	"fmt"
	"time"

	"github.com/steven004/goldenfox-sudoku/engine"
)

// GameManager manages the game state and coordinates gameplay
type GameManager struct {
	// Game Components
	timer     *GameTimer
	history   *HistoryManager
	generator engine.PuzzleGenerator
	userData  *UserData

	// Game State
	initialBoard    *engine.SudokuBoard
	currentBoard    *engine.SudokuBoard
	difficulty      engine.DifficultyLevel
	difficultyIndex float64 // Specific difficulty index (e.g. 1.2)
	currentGameID   string

	// Conflict State
	conflictRow   int
	conflictCol   int
	conflictValue int

	// Limits
	eraseCount int
	undoCount  int
}

// NewGameManager creates a new game manager with the specified generator
func NewGameManager(generator engine.PuzzleGenerator) *GameManager {
	// Load user data
	ud, err := LoadUserData("user_data.json")
	if err != nil {
		fmt.Printf("Warning: Failed to load user data: %v\n", err)
		ud = NewUserData()
	}

	return &GameManager{
		generator: generator,

		conflictRow:   -1,
		conflictCol:   -1,
		conflictValue: 0,
		timer:         NewGameTimer(),
		history:       NewHistoryManager(),
		userData:      ud,
	}
}

// NewGame generates and starts a new puzzle
// difficultyOverride: optional difficulty string (e.g., "Hard"). If empty, uses User Level.
func (gm *GameManager) NewGame(difficultyOverride string) error {
	// Check for abandonment of previous game
	if gm.currentBoard != nil && !gm.currentBoard.IsSolved() && gm.timer.IsRunning() {
		// Previous game was in progress and not solved -> Record Loss
		gm.userData.RecordLoss(gm.difficulty)
		// Save the loss state
		gm.SaveCurrentGame()
	}

	var targetDifficulty engine.DifficultyLevel

	if difficultyOverride != "" {
		// Manual Difficulty Selection
		switch difficultyOverride {
		case "Beginner":
			targetDifficulty = engine.Beginner
		case "Easy":
			targetDifficulty = engine.Easy
		case "Medium":
			targetDifficulty = engine.Medium
		case "Hard":
			targetDifficulty = engine.Hard
		case "Expert":
			targetDifficulty = engine.Expert
		case "FoxGod":
			targetDifficulty = engine.FoxGod
		default:
			// Fallback to User Level if string is invalid
			targetDifficulty = gm.getDifficultyFromLevel()
		}
	} else {
		// No override -> Use User Level
		targetDifficulty = gm.getDifficultyFromLevel()
	}

	// Determine extra clues based on difficulty and progress
	// Note: If playing "Up" or "Down", we should technically use the delta.
	// But for simplicity, we use the User's current progress within their current level
	// to adjust clues slightly.
	// Or, if manually selecting, maybe we should DISABLE dynamic clues?
	// For now, let's keep it simple: Dynamic Clues always apply relative to the DIFFICULTY base.
	progress := gm.userData.Stats.Progress
	extraClues := CalculateDynamicClues(targetDifficulty, progress)

	// Generate a new puzzle
	puzzle, diffIndex, err := gm.generator.Generate(targetDifficulty, extraClues)
	if err != nil {
		return fmt.Errorf("failed to generate puzzle: %w", err)
	}

	// Store the initial board (for restart)
	gm.initialBoard = puzzle.Clone()

	// Set current board
	gm.currentBoard = puzzle

	// Store difficulty and index
	gm.difficulty = targetDifficulty
	gm.difficultyIndex = diffIndex

	gm.ResetConflict()

	// Reset Timer
	gm.timer.Reset()

	// Reset Limits & History
	gm.eraseCount = 0
	gm.undoCount = 0
	gm.history.Clear()

	// Generate new Game ID
	gm.currentGameID = fmt.Sprintf("%d", time.Now().UnixNano())

	return nil
}

// Helper to determine difficulty from User Level
func (gm *GameManager) getDifficultyFromLevel() engine.DifficultyLevel {
	userLevel := gm.userData.Stats.Level
	if userLevel < 1 {
		userLevel = 1
	}
	if userLevel > 6 {
		userLevel = 6
	}
	return engine.DifficultyLevel(userLevel - 1)
}

// RestartGame resets the current puzzle to its initial state
func (gm *GameManager) RestartGame() error {
	if gm.initialBoard == nil {
		return fmt.Errorf("no game in progress")
	}

	// Clone the initial board
	gm.currentBoard = gm.initialBoard.Clone()

	gm.ResetConflict()

	// Reset Timer
	gm.timer.Reset()

	// Reset Limits & History
	gm.eraseCount = 0
	gm.undoCount = 0
	gm.history.Clear()

	// Generate new Game ID
	gm.currentGameID = fmt.Sprintf("%d", time.Now().UnixNano())

	return nil
}

// ... (existing methods) ...

// ResetConflict clears the conflict state
func (gm *GameManager) ResetConflict() {
	gm.conflictRow = -1
	gm.conflictCol = -1
	gm.conflictValue = 0
}

// GetConflictInfo returns the current conflict state
func (gm *GameManager) GetConflictInfo() (row, col, val int, hasConflict bool) {
	if gm.conflictRow == -1 {
		return 0, 0, 0, false
	}
	return gm.conflictRow, gm.conflictCol, gm.conflictValue, true
}

// InputNumber attempts to input a number into the selected cell (Pen Mode)
// Now accepts coordinates directly.
func (gm *GameManager) InputNumber(row, col, val int) error {
	if gm.currentBoard == nil {
		return fmt.Errorf("no game in progress")
	}

	if row < 0 || row > 8 || col < 0 || col > 8 {
		return fmt.Errorf("invalid position: [%d][%d]", row, col)
	}

	cell := &gm.currentBoard.Cells[row][col]

	// Cannot edit given cells
	if cell.Given {
		return fmt.Errorf("cannot edit given cell")
	}

	// Save state for Undo
	gm.history.Push(gm.currentBoard)

	// PEN MODE: Set value
	if cell.Value == val {
		// Tapping same number clears it
		cell.Value = 0
	} else {
		cell.Value = val
		// Force clear note for this number in this cell if setting value
		delete(cell.Candidates, val)
	}

	// Valdation/Stats Logic
	isValid := gm.currentBoard.IsValidMove(row, col, val)

	if isValid {
		// Valid move
		gm.currentBoard.RemoveCandidateFromPeers(row, col, val)
		gm.ResetConflict()

		// Check for win
		if gm.currentBoard.IsSolved() {
			gm.timer.Stop()
			if err := gm.SaveCurrentGame(); err != nil {
				fmt.Printf("Error auto-saving game: %v\n", err)
			}
		} else {
			go gm.SaveCurrentGame()
		}
	} else {
		// Invalid
		gm.conflictRow = row
		gm.conflictCol = col
		gm.conflictValue = val
	}

	return nil
}

// ToggleCandidate toggles a candidate note in the selected cell
func (gm *GameManager) ToggleCandidate(row, col, val int) error {
	if gm.currentBoard == nil {
		return fmt.Errorf("no game in progress")
	}

	if row < 0 || row > 8 || col < 0 || col > 8 {
		return fmt.Errorf("invalid position: [%d][%d]", row, col)
	}

	cell := &gm.currentBoard.Cells[row][col]

	// Cannot edit given cells
	if cell.Given {
		return fmt.Errorf("cannot edit given cell")
	}

	// Only allow notes if cell is empty (Value == 0)
	if cell.Value != 0 {
		return nil
	}

	// Save state for Undo
	gm.history.Push(gm.currentBoard)

	if cell.Candidates == nil {
		cell.Candidates = make(map[int]bool)
	}
	if cell.Candidates[val] {
		delete(cell.Candidates, val)
	} else {
		cell.Candidates[val] = true
	}

	go gm.SaveCurrentGame()

	return nil
}

// ClearCell removes the value or candidates from a cell
func (gm *GameManager) ClearCell(row, col int) error {
	if gm.currentBoard == nil {
		return fmt.Errorf("no game in progress")
	}

	if row < 0 || row > 8 || col < 0 || col > 8 {
		return fmt.Errorf("invalid position: [%d][%d]", row, col)
	}

	// Check if cell is given
	if gm.currentBoard.Cells[row][col].Given {
		return fmt.Errorf("cannot clear given cell at [%d][%d]", row, col)
	}

	// LOCKOUT LOGIC
	if gm.conflictRow != -1 {
		if row == gm.conflictRow && col == gm.conflictCol {
			// User is clearing the conflicting cell -> Resolve conflict
			gm.ResetConflict()
			return nil // Done, nothing on board to clear (since it wasn't saved)
		} else {
			return fmt.Errorf("must resolve conflict at [%d][%d] first", gm.conflictRow+1, gm.conflictCol+1)
		}
	}

	// Normal Clear
	// Check limit
	if gm.eraseCount >= 3 {
		return fmt.Errorf("no erase chances left")
	}

	// Save state to history
	gm.history.Push(gm.currentBoard)

	if err := gm.currentBoard.ClearCell(row, col); err != nil {
		return err
	}

	gm.eraseCount++

	// Clear value
	gm.currentBoard.Cells[row][col].Value = 0
	gm.currentBoard.Cells[row][col].Candidates = make(map[int]bool)

	// Auto-save on clear
	go gm.SaveCurrentGame()

	return nil
}

// GetBoard returns the current board state (read-only access)
func (gm *GameManager) GetBoard() *engine.SudokuBoard {
	return gm.currentBoard
}

// GetDifficulty returns the current game's difficulty level
func (gm *GameManager) GetDifficulty() engine.DifficultyLevel {
	return gm.difficulty
}

// IsSolved checks if the current puzzle is completely and correctly solved
func (gm *GameManager) IsSolved() bool {
	if gm.currentBoard == nil {
		return false
	}
	return gm.currentBoard.IsSolved()
}

// FindConflicts returns all cells that violate Sudoku rules
func (gm *GameManager) FindConflicts() []engine.Coordinate {
	if gm.currentBoard == nil {
		return nil
	}
	return gm.currentBoard.FindConflicts()
}

// GetHistory returns the list of puzzle records
func (gm *GameManager) GetHistory() []PuzzleRecord {
	if gm.userData == nil {
		return nil
	}
	gm.userData.mu.RLock()
	defer gm.userData.mu.RUnlock()

	// Return a copy to avoid race conditions
	history := make([]PuzzleRecord, len(gm.userData.History))
	copy(history, gm.userData.History)
	return history
}

// SaveCurrentGame saves the current game state (even if not finished)
func (gm *GameManager) SaveCurrentGame() error {
	if gm.currentBoard == nil {
		return fmt.Errorf("no game in progress")
	}
	if gm.currentGameID == "" {
		return fmt.Errorf("no current game ID to save")
	}

	// Calculate Time Elapsed
	elapsed := gm.timer.GetElapsedDuration()

	record := PuzzleRecord{
		ID:              gm.currentGameID,
		Predefined:      gm.initialBoard.String(),
		FinalState:      gm.currentBoard.String(),
		IsSolved:        gm.currentBoard.IsSolved(),
		TimeElapsed:     elapsed,
		PlayedAt:        time.Now(),
		Difficulty:      gm.difficulty,
		DifficultyIndex: gm.difficultyIndex, // Added: Save difficulty index
		Mistakes:        0,                  // TODO: Track mistakes
	}

	// Use Upsert to handle both new and existing records + stats
	gm.userData.UpsertPuzzleRecord(record)

	return gm.userData.Save("user_data.json")
}

// LoadGame loads a game from a history record
func (gm *GameManager) LoadGame(id string) error {
	var record *PuzzleRecord

	gm.userData.mu.RLock()
	for i := range gm.userData.History {
		if gm.userData.History[i].ID == id {
			record = &gm.userData.History[i]
			break
		}
	}
	gm.userData.mu.RUnlock()

	if record == nil {
		return fmt.Errorf("record not found: %s", id)
	}

	// Reconstruct Boards
	initial, err := engine.ParseBoard(record.Predefined)
	if err != nil {
		return fmt.Errorf("failed to parse initial board: %w", err)
	}

	current, err := engine.ParseBoard(record.FinalState)
	if err != nil {
		return fmt.Errorf("failed to parse current board: %w", err)
	}

	// Restore 'Given' status based on initial board
	// Any non-zero value in initial board is a Given
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			val, _ := initial.GetValue(r, c)
			if val != 0 {
				initial.Cells[r][c].Given = true
				current.Cells[r][c].Given = true
			}
		}
	}

	gm.initialBoard = initial
	gm.currentBoard = current
	gm.difficulty = record.Difficulty
	gm.difficultyIndex = record.DifficultyIndex // Added: Load difficulty index
	gm.currentGameID = record.ID                // Set current ID to loaded ID

	// Restore Timer
	// We want to resume from where we left off.
	// startTime = Now - TimeElapsed
	gm.timer.SetStartTime(time.Now().Add(-record.TimeElapsed))

	if record.IsSolved {
		gm.timer.SetEndTime(time.Now()) // Mark as ended
	} else {
		gm.timer.SetEndTime(time.Time{}) // Ensure running
	}

	gm.ResetConflict()
	gm.history.Clear() // Clear undo stack on load

	return nil
}

// GetUserLevel returns the user's current level
func (gm *GameManager) GetUserLevel() int {
	if gm.userData == nil {
		return 1
	}
	gm.userData.mu.RLock()
	defer gm.userData.mu.RUnlock()
	return gm.userData.Stats.Level
}

// GetGameState returns the current game state for the UI
func (gm *GameManager) GetGameState() GameState {
	board := engine.SudokuBoard{}
	if gm.currentBoard != nil {
		// Clone the board to avoid modifying the actual game state
		board = *gm.currentBoard.Clone()

		// Inject transient conflict if any
		if gm.conflictRow != -1 {
			board.Cells[gm.conflictRow][gm.conflictCol].Value = gm.conflictValue
		}

		// Check for conflicts and mark invalid cells
		// We use the cloned board's FindConflicts method to include the transient value
		conflicts := board.FindConflicts()
		for _, coord := range conflicts {
			board.Cells[coord.Row][coord.Col].IsInvalid = true
		}
	}

	// Get stats safely
	level := 1
	gamesPlayed := 0
	avgTimeStr := "--:--"
	winRate := 0.0
	pendingGames := 0
	currentDiffCount := 0
	progress := 0
	remainingCells := 81 // Default if no board

	if gm.currentBoard != nil {
		// Count remaining cells
		filled := 0
		for r := 0; r < 9; r++ {
			for c := 0; c < 9; c++ {
				if gm.currentBoard.Cells[r][c].Value != 0 {
					filled++
				}
			}
		}
		remainingCells = 81 - filled
	}

	if gm.userData != nil {
		gm.userData.mu.RLock()
		level = gm.userData.Stats.Level
		gamesPlayed = len(gm.userData.History)
		winRate = gm.userData.GetWinRate()
		pendingGames = gm.userData.GetPendingGamesCount()
		currentDiffCount = gm.userData.GetGamesAtDifficulty(gm.difficulty)
		progress = gm.userData.Stats.Progress

		// Average time for CURRENT difficulty
		if avg, ok := gm.userData.Stats.AverageTimes[gm.difficulty]; ok && avg > 0 {
			// Format seconds to MM:SS
			minutes := int(avg) / 60
			seconds := int(avg) % 60
			avgTimeStr = fmt.Sprintf("%02d:%02d", minutes, seconds)
		}
		gm.userData.mu.RUnlock()
	}

	return GameState{
		Board:                  board,
		Mistakes:               0,
		EraseCount:             gm.eraseCount,
		UndoCount:              gm.undoCount,
		ElapsedSeconds:         int(gm.timer.GetElapsedDuration().Seconds()),
		Difficulty:             gm.difficulty.String(),
		DifficultyIndex:        gm.difficultyIndex,
		IsSolved:               gm.currentBoard != nil && gm.currentBoard.IsSolved(),
		UserLevel:              level,
		GamesPlayed:            gamesPlayed,
		AverageTime:            avgTimeStr,
		WinRate:                winRate,
		PendingGames:           pendingGames,
		CurrentDifficultyCount: currentDiffCount,
		Progress:               progress,
		RemainingCells:         remainingCells,
	}
}

// Undo reverts the last move
func (gm *GameManager) Undo() error {
	if gm.currentBoard == nil {
		return fmt.Errorf("no game in progress")
	}

	if gm.undoCount >= 3 {
		return fmt.Errorf("no undo chances left")
	}

	// Pop last state
	previousBoard, err := gm.history.Pop()
	if err != nil {
		return fmt.Errorf("nothing to undo")
	}

	// Restore board
	gm.currentBoard = previousBoard
	gm.undoCount++

	// Reset conflict state as we reverted to a (presumably) valid state
	gm.ResetConflict()

	// Auto-save
	go gm.SaveCurrentGame()

	return nil
}

// GetLastGameID returns the ID of the last played game from history, or empty string if none
func (gm *GameManager) GetLastGameID() string {
	if gm.userData == nil {
		return ""
	}
	gm.userData.mu.RLock()
	defer gm.userData.mu.RUnlock()

	if len(gm.userData.History) == 0 {
		return ""
	}

	// Return the ID of the last record
	// Assuming History is appended to, so last element is latest
	lastRecord := gm.userData.History[len(gm.userData.History)-1]

	// Only return if it's NOT solved (i.e., resumeable)
	if !lastRecord.IsSolved {
		return lastRecord.ID
	}

	return ""
}
