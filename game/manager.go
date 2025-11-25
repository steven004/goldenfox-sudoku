package game

import (
	"fmt"

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
}

// NewGameManager creates a new game manager with the specified generator
func NewGameManager(generator engine.PuzzleGenerator) *GameManager {
	return &GameManager{
		generator:   generator,
		selectedRow: -1, // No selection initially
		selectedCol: -1,
		pencilMode:  false,
	}
}

// NewGame generates and starts a new puzzle of the specified difficulty
func (gm *GameManager) NewGame(difficulty engine.DifficultyLevel) error {
	// Generate a new puzzle
	puzzle, err := gm.generator.Generate(difficulty)
	if err != nil {
		return fmt.Errorf("failed to generate puzzle: %w", err)
	}

	// Store the initial board (for restart)
	gm.initialBoard = puzzle.Clone()

	// Set current board
	gm.currentBoard = puzzle

	// Store difficulty
	gm.difficulty = difficulty

	// Reset selection
	gm.selectedRow = -1
	gm.selectedCol = -1
	gm.pencilMode = false

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

	return nil
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

	if gm.pencilMode {
		// Add pencil note
		return gm.currentBoard.AddCandidate(row, col, val)
	} else {
		// Place number
		// First check if it's a valid move
		if !gm.currentBoard.IsValidMove(row, col, val) {
			// Still allow the move but it will create conflicts
			// The GUI can highlight conflicts using FindConflicts()
		}

		// Set the value
		if err := gm.currentBoard.SetValue(row, col, val); err != nil {
			return err
		}

		// Auto-remove this candidate from peer cells
		gm.currentBoard.RemoveCandidateFromPeers(row, col, val)

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

	// Clear the value
	if err := gm.currentBoard.ClearCell(row, col); err != nil {
		return err
	}

	// Also clear all candidates
	gm.currentBoard.Cells[row][col].Candidates = make(map[int]bool)

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
