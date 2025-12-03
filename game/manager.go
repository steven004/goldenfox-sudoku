package game

import (
	"fmt"
	"time"

	"github.com/steven004/goldenfox-sudoku/engine"
)

// GameManager manages the game state and coordinates gameplay
type GameManager struct {
	currentBoard *engine.SudokuBoard
	initialBoard *engine.SudokuBoard
	generator    engine.PuzzleGenerator
	difficulty   engine.DifficultyLevel
	selectedRow  int
	selectedCol  int
	pencilMode   bool

	// Conflict State
	conflictRow   int
	conflictCol   int
	conflictValue int

	// Timer State
	startTime  time.Time
	endTime    time.Time
	pausedTime time.Duration

	// User Data
	userData      *UserData
	currentGameID string

	// Limits & History
	eraseCount  int
	undoCount   int
	moveHistory []*engine.SudokuBoard
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
		generator:     generator,
		selectedRow:   -1,
		selectedCol:   -1,
		pencilMode:    false,
		conflictRow:   -1,
		conflictCol:   -1,
		conflictValue: 0,
		startTime:     time.Now(),
		userData:      ud,
		moveHistory:   make([]*engine.SudokuBoard, 0),
	}
}

// NewGame generates and starts a new puzzle based on User Level
func (gm *GameManager) NewGame(difficulty engine.DifficultyLevel) error {
	// Check for abandonment of previous game
	if gm.currentBoard != nil && !gm.currentBoard.IsSolved() && gm.endTime.IsZero() {
		// Previous game was in progress and not solved -> Record Loss
		gm.userData.RecordLoss()
		// Save the loss state
		gm.SaveCurrentGame()
	}

	// Determine difficulty from User Level
	// Level 1 -> Beginner (0), Level 6 -> FoxGod (5)
	userLevel := gm.userData.Stats.Level
	if userLevel < 1 {
		userLevel = 1
	}
	if userLevel > 6 {
		userLevel = 6
	}
	targetDifficulty := engine.DifficultyLevel(userLevel - 1)

	// Generate a new puzzle
	puzzle, err := gm.generator.Generate(targetDifficulty)
	if err != nil {
		return fmt.Errorf("failed to generate puzzle: %w", err)
	}

	// Store the initial board (for restart)
	gm.initialBoard = puzzle.Clone()

	// Set current board
	gm.currentBoard = puzzle

	// Store difficulty
	gm.difficulty = targetDifficulty

	// Reset selection
	gm.selectedRow = -1
	gm.selectedCol = -1
	gm.pencilMode = false
	gm.ResetConflict()

	// Reset Timer
	gm.startTime = time.Now()
	gm.endTime = time.Time{} // Reset end time
	gm.pausedTime = 0

	// Reset Limits & History
	gm.eraseCount = 0
	gm.undoCount = 0
	gm.moveHistory = make([]*engine.SudokuBoard, 0)

	// Generate new Game ID
	gm.currentGameID = fmt.Sprintf("%d", time.Now().UnixNano())

	return nil
}

// RestartGame resets the current puzzle to its initial state
func (gm *GameManager) RestartGame() error {
	if gm.initialBoard == nil {
		return fmt.Errorf("no game in progress")
	}

	// Clone the initial board
	gm.currentBoard = gm.initialBoard.Clone()

	// Reset selection and mode
	gm.selectedRow = -1
	gm.selectedCol = -1
	gm.pencilMode = false
	gm.ResetConflict()

	// Reset Timer
	gm.startTime = time.Now()
	gm.endTime = time.Time{} // Reset end time
	gm.pausedTime = 0

	// Generate new Game ID
	gm.currentGameID = fmt.Sprintf("%d", time.Now().UnixNano())

	return nil
}

// ... (existing methods) ...

// ... (existing methods) ...

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

// SelectCell sets the currently selected cell
func (gm *GameManager) SelectCell(row, col int) error {
	if row < 0 || row > 8 || col < 0 || col > 8 {
		return fmt.Errorf("invalid cell position: [%d][%d]", row, col)
	}

	gm.selectedRow = row
	gm.selectedCol = col
	return nil
}

// GetSelectedCell returns the currently selected cell coordinates
func (gm *GameManager) GetSelectedCell() (row, col int, hasSelection bool) {
	if gm.selectedRow == -1 || gm.selectedCol == -1 {
		return 0, 0, false
	}
	return gm.selectedRow, gm.selectedCol, true
}

// ClearSelection clears the current cell selection
func (gm *GameManager) ClearSelection() {
	gm.selectedRow = -1
	gm.selectedCol = -1
}

// TogglePencilMode switches between number input and pencil note mode
func (gm *GameManager) TogglePencilMode() {
	gm.pencilMode = !gm.pencilMode
}

// IsPencilMode returns true if pencil mode is active
func (gm *GameManager) IsPencilMode() bool {
	return gm.pencilMode
}

// SetPencilMode explicitly sets the pencil mode state
func (gm *GameManager) SetPencilMode(enabled bool) {
	gm.pencilMode = enabled
}

// InputNumber places a number or pencil note at the specified position
func (gm *GameManager) InputNumber(row, col, val int) error {
	if gm.currentBoard == nil {
		return fmt.Errorf("no game in progress")
	}

	if row < 0 || row > 8 || col < 0 || col > 8 {
		return fmt.Errorf("invalid position: [%d][%d]", row, col)
	}

	if val < 1 || val > 9 {
		return fmt.Errorf("invalid value: %d (must be 1-9)", val)
	}

	// Check if cell is given
	if gm.currentBoard.Cells[row][col].Given {
		return fmt.Errorf("cannot modify given cell at [%d][%d]", row, col)
	}

	// Check if already solved (don't allow moves if game is over)
	if !gm.endTime.IsZero() {
		return nil
	}

	// LOCKOUT LOGIC: Check if there is an active conflict
	if gm.conflictRow != -1 {
		// If there is a conflict, user MUST interact with the conflicting cell
		if row != gm.conflictRow || col != gm.conflictCol {
			return fmt.Errorf("must resolve conflict at [%d][%d] first", gm.conflictRow+1, gm.conflictCol+1)
		}
		// If interacting with the conflicting cell, proceed to validation below
	}

	if gm.pencilMode {
		// Pencil mode is allowed even during conflict?
		// User said "erase it before fill other numbers".
		// Let's assume pencil marks are fine or blocked?
		// Safest is to BLOCK pencil marks on OTHER cells too.
		if gm.conflictRow != -1 && (row != gm.conflictRow || col != gm.conflictCol) {
			return fmt.Errorf("must resolve conflict at [%d][%d] first", gm.conflictRow+1, gm.conflictCol+1)
		}

		// Toggle candidate: if exists, remove it; otherwise, add it.
		candidates, err := gm.currentBoard.GetCandidates(row, col)
		if err != nil {
			return err
		}

		exists := false
		for _, c := range candidates {
			if c == val {
				exists = true
				break
			}
		}

		if exists {
			return gm.currentBoard.RemoveCandidate(row, col, val)
		} else {
			return gm.currentBoard.AddCandidate(row, col, val)
		}
	} else {
		// Place number

		// Check validity
		isValid := gm.currentBoard.IsValidMove(row, col, val)

		if isValid {
			// Valid move: Set value and clear any conflict state
			// Save state to history before modifying
			gm.pushHistory()

			if err := gm.currentBoard.SetValue(row, col, val); err != nil {
				return err
			}
			gm.currentBoard.RemoveCandidateFromPeers(row, col, val)
			gm.ResetConflict() // Conflict resolved

			// Check for win condition
			if gm.currentBoard.IsSolved() && gm.endTime.IsZero() {
				gm.endTime = time.Now()

				// Auto-save on win
				// Auto-save on win
				if err := gm.SaveCurrentGame(); err != nil {
					fmt.Printf("Error auto-saving game: %v\n", err)
				}
			} else {
				// Auto-save progress on valid move
				// We ignore errors here to avoid interrupting gameplay
				go gm.SaveCurrentGame()
			}
		} else {
			// Invalid move: Do NOT set value on board. Set Transient Conflict.
			gm.conflictRow = row
			gm.conflictCol = col
			gm.conflictValue = val
			// Note: We do NOT call SetValue, so board data remains clean.
		}

		return nil
	}
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
	gm.pushHistory()

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

// IsNumberComplete returns true if all 9 instances of a number are placed
func (gm *GameManager) IsNumberComplete(val int) bool {
	if gm.currentBoard == nil {
		return false
	}
	return gm.currentBoard.CountNumber(val) == 9
}

// GetCellValue returns the value at the specified position
func (gm *GameManager) GetCellValue(row, col int) (int, error) {
	if gm.currentBoard == nil {
		return 0, fmt.Errorf("no game in progress")
	}
	return gm.currentBoard.GetValue(row, col)
}

// GetCellCandidates returns the pencil notes for a cell
func (gm *GameManager) GetCellCandidates(row, col int) ([]int, error) {
	if gm.currentBoard == nil {
		return nil, fmt.Errorf("no game in progress")
	}
	return gm.currentBoard.GetCandidates(row, col)
}

// IsCellGiven returns true if the cell is a given clue
func (gm *GameManager) IsCellGiven(row, col int) bool {
	if gm.currentBoard == nil {
		return false
	}
	if row < 0 || row > 8 || col < 0 || col > 8 {
		return false
	}
	return gm.currentBoard.Cells[row][col].Given
}

// GetMistakes returns the current mistake count
func (gm *GameManager) GetMistakes() int {
	// TODO: Implement actual mistake tracking
	return 0
}

// GetElapsedTime returns the formatted elapsed time string
func (gm *GameManager) GetElapsedTime() string {
	if gm.currentBoard == nil {
		return "00:00"
	}

	var elapsed time.Duration
	if !gm.endTime.IsZero() {
		elapsed = gm.endTime.Sub(gm.startTime)
	} else {
		elapsed = time.Since(gm.startTime)
	}

	// Format as MM:SS
	minutes := int(elapsed.Minutes())
	seconds := int(elapsed.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
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
	var elapsed time.Duration
	if !gm.endTime.IsZero() {
		elapsed = gm.endTime.Sub(gm.startTime)
	} else {
		elapsed = time.Since(gm.startTime)
	}

	record := PuzzleRecord{
		ID:          gm.currentGameID,
		Predefined:  gm.initialBoard.String(),
		FinalState:  gm.currentBoard.String(),
		IsSolved:    gm.currentBoard.IsSolved(),
		TimeElapsed: elapsed,
		PlayedAt:    time.Now(),
		Difficulty:  gm.difficulty,
		Mistakes:    0, // TODO: Track mistakes
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
	gm.currentGameID = record.ID // Set current ID to loaded ID

	// Restore Timer
	// We want to resume from where we left off.
	// startTime = Now - TimeElapsed
	gm.startTime = time.Now().Add(-record.TimeElapsed)

	if record.IsSolved {
		gm.endTime = time.Now() // Mark as ended
	} else {
		gm.endTime = time.Time{} // Ensure running
	}

	gm.pausedTime = 0
	gm.ResetConflict()

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

	selected := gm.selectedRow != -1 && gm.selectedCol != -1

	// Get stats safely
	level := 1
	gamesPlayed := 0
	avgTimeStr := "--:--"
	winRate := 0.0
	pendingGames := 0
	currentDiffCount := 0
	consecutiveWins := 0
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
		consecutiveWins = gm.userData.Stats.ConsecutiveWins

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
		SelectedRow:            gm.selectedRow,
		SelectedCol:            gm.selectedCol,
		IsSelected:             selected,
		PencilMode:             gm.pencilMode,
		Mistakes:               0,
		EraseCount:             gm.eraseCount,
		UndoCount:              gm.undoCount,
		TimeElapsed:            gm.GetElapsedTime(),
		Difficulty:             gm.difficulty.String(),
		IsSolved:               gm.currentBoard != nil && gm.currentBoard.IsSolved(),
		UserLevel:              level,
		GamesPlayed:            gamesPlayed,
		AverageTime:            avgTimeStr,
		WinRate:                winRate,
		PendingGames:           pendingGames,
		CurrentDifficultyCount: currentDiffCount,
		WinsForNextLevel:       5 - consecutiveWins,
		RemainingCells:         remainingCells,
	}
}

// pushHistory saves the current board state to history
func (gm *GameManager) pushHistory() {
	if gm.currentBoard == nil {
		return
	}
	// Clone current board
	snapshot := gm.currentBoard.Clone()
	gm.moveHistory = append(gm.moveHistory, snapshot)
}

// Undo reverts the last move
func (gm *GameManager) Undo() error {
	if gm.currentBoard == nil {
		return fmt.Errorf("no game in progress")
	}

	if gm.undoCount >= 3 {
		return fmt.Errorf("no undo chances left")
	}

	if len(gm.moveHistory) == 0 {
		return fmt.Errorf("nothing to undo")
	}

	// Pop last state
	lastIndex := len(gm.moveHistory) - 1
	previousBoard := gm.moveHistory[lastIndex]
	gm.moveHistory = gm.moveHistory[:lastIndex]

	// Restore board
	gm.currentBoard = previousBoard
	gm.undoCount++

	// Reset conflict state as we reverted to a (presumably) valid state
	gm.ResetConflict()

	// Auto-save
	go gm.SaveCurrentGame()

	return nil
}
