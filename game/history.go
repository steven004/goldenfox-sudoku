package game

import (
	"errors"

	"github.com/steven004/goldenfox-sudoku/engine"
)

// HistoryManager handles the undo stack for the game
type HistoryManager struct {
	stack []*engine.SudokuBoard
}

// NewHistoryManager creates a new history manager
func NewHistoryManager() *HistoryManager {
	return &HistoryManager{
		stack: make([]*engine.SudokuBoard, 0),
	}
}

// Push saves a snapshot of the board to the history stack
func (hm *HistoryManager) Push(board *engine.SudokuBoard) {
	if board == nil {
		return
	}
	// Clone the board to ensure we store a snapshot, not a reference
	snapshot := board.Clone()
	hm.stack = append(hm.stack, snapshot)
}

// Pop restores the last saved board state
func (hm *HistoryManager) Pop() (*engine.SudokuBoard, error) {
	if len(hm.stack) == 0 {
		return nil, errors.New("history is empty")
	}

	lastIndex := len(hm.stack) - 1
	previousBoard := hm.stack[lastIndex]

	// Remove the last element
	hm.stack = hm.stack[:lastIndex]

	return previousBoard, nil
}

// Clear resets the history stack
func (hm *HistoryManager) Clear() {
	hm.stack = make([]*engine.SudokuBoard, 0)
}

// Count returns the number of moves in the history
func (hm *HistoryManager) Count() int {
	return len(hm.stack)
}
