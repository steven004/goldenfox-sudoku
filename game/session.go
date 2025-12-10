package game

import (
	"fmt"
	"time"

	"github.com/steven004/goldenfox-sudoku/engine"
)

// GameSession represents a single active game session
type GameSession struct {
	// Board State
	initialBoard *engine.SudokuBoard
	currentBoard *engine.SudokuBoard

	// Metadata
	difficulty      engine.DifficultyLevel
	difficultyIndex float64
	gameID          string

	// Logic Components
	history *HistoryManager

	// Timekeeping (Inlined)
	startTime time.Time
	endTime   time.Time

	// Counters
	eraseCount int
	undoCount  int
}

// NewGameSession creates a new session
func NewGameSession(puzzle *engine.SudokuBoard, difficulty engine.DifficultyLevel, diffIndex float64, id string) *GameSession {
	return &GameSession{
		initialBoard:    puzzle.Clone(),
		currentBoard:    puzzle,
		difficulty:      difficulty,
		difficultyIndex: diffIndex,
		gameID:          id,
		history:         NewHistoryManager(),
		startTime:       time.Now(),
		eraseCount:      0,
		undoCount:       0,
	}
}

// GetElapsedDuration calculates time played
func (s *GameSession) GetElapsedDuration() time.Duration {
	if !s.endTime.IsZero() {
		return s.endTime.Sub(s.startTime)
	}
	return time.Since(s.startTime)
}

// InputNumber inputs a number (validates conflict) (Pure Logic)
func (s *GameSession) InputNumber(row, col, val int) error {
	if s.currentBoard == nil {
		return fmt.Errorf("no board")
	}

	if !s.currentBoard.IsValidMove(row, col, val) {
		return fmt.Errorf("invalid move (conflict)")
	}

	cell := &s.currentBoard.Cells[row][col]
	if cell.Given {
		return fmt.Errorf("cannot edit given cell")
	}

	// Push History
	s.history.Push(s.currentBoard)

	if cell.Value == val {
		cell.Value = 0 // Toggle off
	} else {
		cell.Value = val
		delete(cell.Candidates, val)
	}

	s.currentBoard.RemoveCandidateFromPeers(row, col, val)

	if s.currentBoard.IsSolved() {
		s.endTime = time.Now()
	}

	return nil
}

// ToggleCandidate toggles a pencil mark (Pure Logic)
func (s *GameSession) ToggleCandidate(row, col, val int) error {
	if s.currentBoard == nil {
		return fmt.Errorf("no board")
	}

	if !s.currentBoard.IsValidMove(row, col, val) {
		return fmt.Errorf("invalid candidate (conflict)")
	}

	cell := &s.currentBoard.Cells[row][col]
	if cell.Given {
		return fmt.Errorf("cannot edit given")
	}
	if cell.Value != 0 {
		return nil
	}

	s.history.Push(s.currentBoard)

	if cell.Candidates == nil {
		cell.Candidates = make(map[int]bool)
	}
	if cell.Candidates[val] {
		delete(cell.Candidates, val)
	} else {
		cell.Candidates[val] = true
	}

	return nil
}

// ClearCell clears a cell (Pure Logic)
func (s *GameSession) ClearCell(row, col int) error {
	// Check limit
	if s.eraseCount >= 3 {
		return fmt.Errorf("no erase chances left")
	}

	cell := &s.currentBoard.Cells[row][col]
	if cell.Given {
		return fmt.Errorf("cannot clear given")
	}

	s.history.Push(s.currentBoard)

	cell.Value = 0
	cell.Candidates = make(map[int]bool)

	s.eraseCount++
	return nil
}

// Undo reverts last move (Pure Logic)
func (s *GameSession) Undo() error {
	if s.undoCount >= 3 {
		return fmt.Errorf("no undo chances left")
	}

	prev, err := s.history.Pop()
	if err != nil {
		return fmt.Errorf("nothing to undo")
	}

	s.currentBoard = prev
	s.undoCount++
	return nil
}
